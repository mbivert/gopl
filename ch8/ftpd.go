package main

// NOTE: this is a rough draft.
//
// Some tips:
//
//	- The server will typically read commands one line at a time from
//	the client, and is expected to answer messages starting with
//	a status code (integer) followed by some (optional IIRC) text.
//
//	- A core idea behind FTP is that there are two sockets; one
//	to send commands (the "initial" one), and one to send data. By
//	default, the second socket is opened *on the client*: the server
//	must connect to whatever the client is telling him to.
//
//	The passive mode allows to have both sockets opened by the server.
//
//	- For tests, you can use "ftp -d -e -4 localhost 8000" as a client
//
//	- The list of FTP server return codes[0], coupled with section 5.4
//	of RFC 959[1] should tell you pretty much what you need to do.
//
//	The codes are also listed in section 4.2 of the previous RFC.
//
// [0]: https://en.wikipedia.org/wiki/List_of_FTP_server_return_codes
// [1]: https://datatracker.ietf.org/doc/html/rfc959

import (
	"bufio"
	"flag"
	"net"
	"strings"
	"log"
	"io"
	"fmt"
	"os"
	"strconv"
	"path/filepath"
)

const (
	// context key for the data socket
	dataSock = "data-sock"
)

func endData(conn, sock net.Conn, ctx map[string]any) {
	err := sock.Close()
	delete(ctx, dataSock)

	// at least log it
	if err != nil {
		log.Println("error: failed to close socket: %s\r\n", err)
	}

	fmt.Fprintf(conn, "250 transfer complete\r\n")
}

// TODO: support LIST <arg>
func list(conn net.Conn, arg string, ctx map[string]any) {
	// would be better to test those !ok
	isock, _ := ctx[dataSock]
	sock, _ := isock.(net.Conn)
	iwd, _ := ctx["wd"]
	wd, _ := iwd.(string)

	xs, err := os.ReadDir(wd)
	if err != nil {
		fmt.Fprintf(conn, "500 internal error: failed to open '%s': %s\r\n", err)
		return
	}

	_, err = fmt.Fprintf(conn, "125 Data connection already open; transfer starting\r\n")

	if err != nil {
		log.Println(err)
		return
	}
	for _, x := range xs {
		_, err := fmt.Fprintf(sock, "%s\r\n", x.Name())
		if err != nil {
			fmt.Fprintf(conn, "500 failed to write filename '%s': %s\r\n", x.Name(), err)
			return
		}
	}

	endData(conn, sock, ctx)
}

func pwd(conn net.Conn, _ string, ctx map[string]any) {
	iwd, _ := ctx["wd"]
	wd, _ := iwd.(string)
	fmt.Fprintf(conn, "257 \"%s\" is current directory\r\n", wd)
}

func retr(conn net.Conn, fn string, ctx map[string]any) {
	if fn == "" {
		fmt.Fprintf(conn, "501 No file specified")
		return
	}

	// would be better to test those !ok
	isock, _ := ctx[dataSock]
	sock, _ := isock.(net.Conn)
	iwd, _ := ctx["wd"]
	wd, _ := iwd.(string)

	_, err := fmt.Fprintf(conn, "125 Data connection already open; transfer starting\r\n")

	fd, err := os.Open(filepath.Join(wd, fn))
	if err != nil {
		fmt.Fprintf(conn, "550 failed: %s\r\n", err)
		return
	}

	if _, err := io.Copy(sock, fd); err != nil {
		fmt.Fprintf(conn, "500 failed: %s\r\n", err)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}

	endData(conn, sock, ctx)
}

func cwd(conn net.Conn, arg string, ctx map[string]any) {
	if arg == "" {
		fmt.Fprintf(conn, "250 requested file action okay, completed\r\n")
		return
	}

	if arg[0] == '/' {
		ctx["wd"] = arg
	} else {

		iwd, _ := ctx["wd"]
		wd, _ := iwd.(string)
		wd = filepath.Join(wd, arg)

		fi, err := os.Stat(wd)
		if err != nil {
			fmt.Fprintf(conn, "500 failed: %s\r\n", err)
			return
		}
		if !fi.IsDir() {
			fmt.Fprintf(conn, "500 failed: %s is not a directory\r\n", wd)
			return
		}
		ctx["wd"] = wd
	}

	fmt.Fprintf(conn, "250 requested file action okay, completed\r\n")
}

func user(conn net.Conn, user string, ctx map[string]any) {
	ctx["user"] = user
	fmt.Fprintf(conn, "230 Login successful\r\n")
}

func syst(conn net.Conn, _ string, ctx map[string]any) {
	fmt.Fprintf(conn, "215 ftpd.go\r\n")
}

func quit(conn net.Conn, _ string, ctx map[string]any) {
	conn.Close()
}

func port(conn net.Conn, args string, ctx map[string]any) {
	xs := strings.Split(args, ",")
	if len(xs) != 6 {
		fmt.Fprintf(conn, "500 unexpected PORT value %s\r\n", args)
		return
	}

	ip := strings.Join(xs[0:4], ".")

	ph, err := strconv.Atoi(xs[4])
	if err != nil {
		fmt.Fprintf(conn, "500 %s is not an integer %s\r\n", xs[4])
		return
	}
	pl, err := strconv.Atoi(xs[5])
	if err != nil {
		fmt.Fprintf(conn, "500 %s is not an integer %s\r\n", xs[5])
		return
	}

	port := ph * 256 + pl
	println(port)

	c, err := net.Dial("tcp", xs[1]+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Fprintf(conn, "425 failed to reach %s:%d: %s\r\n", ip, port, err)
		fmt.Println(err)
		return
	}
	ctx[dataSock] = c

	fmt.Fprintf(conn, "225 listening on %s:%d\r\n", ip, port)
}

var cmds = map[string](func(net.Conn, string, map[string]any) ){
	"LIST" : list,
	"RETR" : retr,
	"CWD"  : cwd,
	"USER" : user,
	"SYST" : syst,
	"PWD"  : pwd,
	"XPWD" : pwd,
	"QUIT" : quit,
	"PORT" : port,
}

func serve(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	// say hello
	_, err := fmt.Fprintf(conn, "200 OK\r\n")
	if err != nil {
		log.Println(err)
		return
	}

	// To store user, data socket, working directory
	ctx := make(map[string]any)

	ctx["wd"] = "."

	for {
		s, err := r.ReadString('\n')
		// could be finer
		if err != nil {
			log.Println(err)
			return
		}
		s = strings.TrimRight(s, "\r\n")
		println(s)

		xs := strings.SplitN(s, " ", 2)
		cmd, ok := cmds[xs[0]]
		if !ok {
			_, err := fmt.Fprintf(conn, "502 Command not implemented\r\n")
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			if len(xs) == 1 {
				xs = append(xs, "")
			}
			cmd(conn, xs[1], ctx)
		}
	}
}

func main() {
	var port string

//	flag.StringVar(&root, "root", ".",     "root directory")
	flag.StringVar(&port, "port", ":8000", "listening port")
	flag.Parse()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on "+port)

	for {
		conn, err := lis.Accept()
		log.Println(conn)
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go serve(conn)
	}
}
