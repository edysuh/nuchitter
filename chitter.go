package main

import (
	"os"
	"fmt"
	"net"
	// "io"
	"bufio"
	// "strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Arguments missings. Use format \"$ go run nuchitter.go [-c] PORT\"")
		os.Exit(1)
	} else if len(os.Args) < 3 {
			fmt.Println("SERVER listening on port", os.Args[1])
			startServer(os.Args[1])
	} else if len(os.Args) < 4 {
		if os.Args[1] == "-c" {
			fmt.Println("runnning CLIENT on port", os.Args[2])
			startClient(os.Args[2])
		}
	} else {
		fmt.Println("Error: Invalid command")
		os.Exit(1)
	}
}

func startServer(port string) {
	ln, err := net.Listen("tcp", ":" + port)

	if err != nil {
		fmt.Println("Error: Cannot start server:", err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Error: Cannot accept connection:", err)
		}

		go acceptConnection(conn)
	}
}

// func acceptConnection(conn net.Conn) {
// 	buf := make([]byte, 1024)

// 	for {
// 		_, err := conn.Read(buf)

// 		if err != nil {
// 			fmt.Println("error reading")
// 		}

// 		conn.Write([]byte("message received\n"))
// 	}
// }

func acceptConnection(conn net.Conn) {
	bufr := bufio.NewReader(conn)
	buf := make([]byte, 1024)

	for {
		readBytes, err := bufr.Read(buf)
		if err != nil {
			conn.Close()
			return
		}
		fmt.Printf("%s", buf[:readBytes])
		// conn.Write(">" + []byte(string(buf[:readBytes])))
	}
}

func startClient(addrport string) {
	conn, err := net.Dial("tcp", addrport)

	if err != nil {
		fmt.Println("Error: Cannot start client:", err)
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
}
