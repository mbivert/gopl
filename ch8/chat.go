package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	ch   chan<- string
	name string
}

var port = flag.String("p", ":8000", "listening port")

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

// Exercise 8.12
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

func clientWriter(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	name := conn.RemoteAddr().String()

	cli := client{ch, name}

	ch <- "Connected as " + name
	messages <- name + " has entered the chat"
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + ": " + input.Text()
	}

	leaving <- cli
	messages <- name + " has left the chat"
	conn.Close()
}

func main() {
	ln, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
