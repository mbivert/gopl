package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type client struct {
	ch     chan<- string
	name   string
	ticker *time.Ticker
}

var port = flag.String("p", ":8000", "listening port")

// automatically disconnect clients after that much time of idling.
var timeout = 5 * time.Second

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
	timeouts = make(chan net.Conn) // 8.13
)

// 8.12
func lsClients(to client, clients map[client]bool) {
	names := make([]string, len(clients))
	i := 0
	for cli := range clients {
		names[i] = cli.name
		i++
	}
	to.ch <- "Users: " + strings.Join(names, ", ")
}

func broadcaster() {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.ch <- msg
			}
		case cli := <-entering:
			clients[cli] = true
			lsClients(cli, clients)
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

// separated from broadcaster so that we still
// disconnect clients who don't provide their name fast enough.
func timeouter() {
	for conn := range timeouts {
		conn.Close()
	}
}

func clientWriter(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)

	// Set a temporary name & a corresponding client:
	// we want the timeout goroutine below to be running
	// even when user hasn't provided a name yet.
	cli := client{
		ch,
		conn.RemoteAddr().String(),
		time.NewTicker(timeout),
	}

	// 8.13
	go func() {
		select {
		case <- cli.ticker.C:
			timeouts <- conn
		}
		// we'll get there, eventually. that is,
		// even if the client is already long gone.
		//
		// for real code, it might be better to avoid
		// accumulating such goroutines too much.
		cli.ticker.Stop()
	}()

	input := bufio.NewScanner(conn)
	fmt.Fprintf(conn, "Enter your name: ")
	ok := input.Scan()
	if !ok {
		conn.Close()
		return
	}
	cli.ticker.Reset(timeout)
	cli.name = input.Text()

	go clientWriter(conn, ch)

	ch <- "Connected as " + cli.name
	messages <- cli.name + " has entered the chat"
	entering <- cli

	for input.Scan() {
		messages <- cli.name + ": " + input.Text()
		cli.ticker.Reset(timeout) // 8.13
	}

	leaving <- cli
	messages <- cli.name + " has left the chat"

	// Already closed on timeout
	conn.Close()
}

func main() {
	ln, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	go timeouter()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
