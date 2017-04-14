package main

import (
	"os"
	"fmt"
	"net"
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

	messages := make(chan string)
	clients := make(chan net.Conn)
	go manageClients(messages, clients)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error: Cannot accept connection:", err)
		}

		fmt.Println("-")

		clients <- conn

		go handleConnection(conn, messages)
	}
}


func manageClients(clntchan chan net.Conn, msgchan chan string) {
	select {
	case <-clntchan:
			for cl := range clntchan {
				fmt.Println("client channeled here:", cl)
			}
	case <-msgchan:
		for msg := range msgchan {
			fmt.Println("mischief managed:", msg)
		}
	}
}


func handleConnection(conn net.Conn, msg chan string) {
	for {
		recvdMsg, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Println("Error: Connection Interrupted:", err)
			conn.Close()
			break
		}

		// print to server
		fmt.Printf("%s", recvdMsg)
		// reply to client
		conn.Write([]byte("> " + recvdMsg))

		msg <- recvdMsg
	}
}


func startClient(addrport string) {
	conn, err := net.Dial("tcp", addrport)
	reader := bufio.NewReader(os.Stdin)
	if err != nil { fmt.Println("Error: Cannot start client:", err) }

	for {
		sendMsg, err := reader.ReadString('\n')
		if err != nil { fmt.Println("Error: Cannot read input:", err) }

		// send message to server
		conn.Write([]byte("< " + sendMsg))
	}
}
