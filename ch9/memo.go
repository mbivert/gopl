package main

/*
 * Untested / WIP.
 *
 * This feels too complex, and we've altered request despite
 * not being explicitly asked to.
 *
 * Idea: a Func's execution can be altered by a done channel,
 * closed upon cancellation. Because we want that Func to be
 * potentially cancelled from within (e.g. imagine it's making
 * a http request), Get() should accept that cancellation channel,
 * and it should need to be forwarded in the request.
 *
 * We furthermore rely on a having a queue for each (func, key)
 * pair. A cancelled call() will trigger a goroutine waiting to
 * be able to write on that queue. In deliver, goroutines will
 * either wait for a result, or to be contacted on that queue:
 * the firts one to receive bits from the queue will trigger
 * a new computation.
 */

import "fmt"

// Func is the type of the function to memoize.
type Func func(key string, done chan<- struct{}) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
	done     chan<- struct{}
}

type Memo struct{ requests chan request }

// New returns a memoization of f. Clients must subsequently call Close.
func New(f Func) *Memo {

	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	var res result
	response := make(chan result)

	select {
	case <-done:
		res.err = fmt.Errorf("cancelled")
	case memo.requests <- request{key, response, done}:
		select {
		case <-done:
			res.err = fmt.Errorf("cancelled")
		case res = <-response:
		}
	}

	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)

	queues := make(map[string]chan<- struct{})

	for req := range memo.requests {
		queue, ok := queues[req.key]
		if !ok {
			queue = make(chan<- struct{})
		}
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req, queue) // call f(key)
		}
		go e.deliver(f, req, queue)
	}
}

func (e *entry) call(f Func, req request, queue <-chan struct{}) {
	// Evaluate the function.
	e.res.value, e.res.err := f(req.key, done)

	// Broadcast the ready condition only if f wasn't cancelled.
	select {
	case <-req.done:
		// This function call was cancelled, and the corresponding
		// deliver will also gets cancelled. But some other goys
		// may be waiting for delivery. So, try to contact one of them
		// via the queue, and ask him to issue another request.
		//
		// If there's no-one waiting for us now, there may be someone
		// in the future; because there's still a registered cache entry,
		// it'll reach deliver() and will need to be informed that the
		// computation needs to be restarted.
		//
		// But we don't want to block here, hence why we're launching
		// another goroutine.
		go func() {
			select {
			case queue <- struct{}:
			}
		}()
	default:
		close(e.ready)
	}
}

func (e *entry) deliver(f Func, req request, queue <-chan struct{}) {
	select {
	// Le Choosen One to restart the computation.
	case <-queue:
		goto Restart
	case <-req.done:
	// Wait for the ready condition.
	case <-e.ready:
		// Send the result to the client.
		select {
		case <-queue:
			goto Restart
		case <-req.done:
		case req.response <- e.res:
		}
	}
	return

Restart:
	go e.call(f, req, queue)
}
