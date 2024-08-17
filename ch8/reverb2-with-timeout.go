package main

// NOTE: strictly speaking, this is not the echo server from section 8.3,
// but the one from section 8.4.1 (slightly more complex)

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	var wg sync.WaitGroup

	ticker := time.NewTicker(1*time.Second)

	// Wait in the background for clients to send
	// us data; once he does, forward it to the data
	// channel
	//
	// NOTE: an unbuffered channel or a channel of size
	// one would work equally well here
	data   := make(chan string, 1)
	go func() {
		for {
			x := input.Scan()
			// ignoring errors
			if !x {
				// We could stop the ticker too, but it's
				// better to stop it later, regardless of
				// whether we stopped because of the timeout
				// or of the closed input stream.
				close(data)
				return
			}
			data <- input.Text()
		}
	}()

	for countdown := 10; countdown > 0; countdown-- {
		// Wait for one tick unless we receive something
		// from the client on data
		select {
		case <- ticker.C:
		case x, ok := <- data:
			// make sure to detect when the channel is closed
			// (e.g. client sends EOF), for otherwise, we'd reset
			// the countdown.
			if !ok {
				countdown = 0
				break
			}
			// client sent us something: reset the countdown
			countdown = 10
			wg.Add(1)
			go func() {
				defer wg.Done()
				echo(c, x, 1*time.Second)
			}()
		}
	}

	wg.Wait()
	ticker.Stop()
	tc, ok := c.(*net.TCPConn)
	if !ok {
		log.Println("(╯°□°）╯︵ ┻━┻")
		// NOTE: ignoring potential errors from input.Err()
		c.Close()
	} else {
		tc.CloseWrite()
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {

			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
}
