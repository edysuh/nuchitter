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

type Client struct {
	conn net.Conn
	channel chan string
}


func startServer(port string) {
	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		fmt.Println("Error: Cannot start server:", err)
	}

	messages := make(chan string)
	addClient := make(chan Client)
	go handleChannels(messages, addClient)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error: Cannot accept connection:", err)
		}

		go handleConnection(conn, messages, addClient)
	}
}


func handleChannels(messagesChan chan string, clientsChan chan Client) {
	clientList := make([]Client, 0)

	for {
		select {
		case cl := <-clientsChan:
			fmt.Println("clientsChan")
			clientList = append(clientList, cl)

			fmt.Println("clientList:", clientList)
			// for cl := range clientsChan {
			// 	fmt.Println("client channeled here:", cl)
			// }
		case msg := <-messagesChan:
			fmt.Println("messagesChan", msg)
			// for msg := range messagesChan {
			// 	fmt.Println("mischief managed:", msg)
			// }
		}
	}
}


func handleConnection(conn net.Conn, msg chan string, addClient chan Client) {
	ch := make(chan string)
	addClient <- Client{conn, ch}

	for {
		recvdMsg, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			fmt.Println("Error: Connection Interrupted:", err)
			conn.Close()
			break
		}

		msg <- recvdMsg

		// // print to server
		// fmt.Printf("%s", recvdMsg)
		// // reply to client
		// conn.Write([]byte("> " + recvdMsg))
	}
}


func startClient(addrport string) {
	conn, err := net.Dial("tcp", addrport)
	reader := bufio.NewReader(os.Stdin)

	if err != nil {
		fmt.Println("Error: Cannot start client:", err)
	}

	for {
		sendMsg, err := reader.ReadString('\n')
		if err != nil { fmt.Println("Error: Cannot read input:", err) }

		// send message to server
		conn.Write([]byte("< " + sendMsg))
	}
}
