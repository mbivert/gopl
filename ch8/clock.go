package main

import (
	"io"
	"log"
	"net"
	"time"
	"flag"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	port := flag.String("port", ":8000", "listening port")
	flag.Parse()

	listener, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on "+*port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}