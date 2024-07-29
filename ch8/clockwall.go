package main

import (
	"strings"
	"os"
	"log"
	"fmt"
	"bufio"
	"net"
)

// trying to avoid channels
func main() {
	var cities []string
	var cs []*bufio.Reader

	for _, arg := range os.Args[1:] {
		xs := strings.SplitN(arg, "=", 2)
		if len(xs) != 2 {
			log.Fatal("Invalid argument format: "+arg)
		}

		c, err := net.Dial("tcp", xs[1])
		if err != nil {
			log.Fatal(err)
		}
		// remember, defering is executed at the end of
		// the surrounding function, not code block.
		defer c.Close()

		cities = append(cities, xs[0])
		cs = append(cs, bufio.NewReader(c))
	}

	for _, city := range cities {
		fmt.Printf("%15s ", city)
	}
	fmt.Printf("\n")
	for {
		for _, c := range cs {
			l, isPrefix, err := c.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			if isPrefix {
				log.Fatal("Clock servers shouldn't send so much: "+string(l))
			}
			fmt.Printf("%15s ", l)
		}
		fmt.Printf("\n")
	}

}
