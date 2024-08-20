package main

// « The message sent to the monitor goroutine must contain
// both the amount to withdraw and a new channel over which
// the monitor goroutine can send the boolean result
// back to Withdraw. »
type withdraw struct {
	amount int
	result chan bool
}

var (
	deposits  = make(chan int)
	balances  = make(chan int)
	withdraws = make(chan withdraw)
)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	w := withdraw{amount, make(chan bool)}
	withdraws <- w
	return <-w.result
}

func teller() {
	var balance int

	for {
		select {
		// well, guess we assume amount to be >= 0
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case w := <-withdraws:
			if balance-w.amount < 0 {
				w.result <- false
			} else {
				balance -= w.amount
				w.result <- true
			}
		}
	}
}

func main() {
	go teller()

	// side-effect: wait for teller() to get running.
	Deposit(500)

	println(Balance())
	println(Withdraw(400))
	println(Balance())
	println(Withdraw(200))
	println(Balance())
}
