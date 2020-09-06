package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func printIncoming(str string) {
	fmt.Printf("< %s\n> ", str)
}

func printStatus(str string) {
	fmt.Printf("- %s\n> ", str)
}

func newConnection(conn net.Conn, cons *Model) {
	cons.AddConn(conn)
	defer cons.RemoveAndCloseConn(conn)

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Client disconnected.")
			break
		}

		if !cons.WasProcessed(msg) {
			cons.MessageProcessed(msg)
			printIncoming(strings.TrimSpace(msg))
			go cons.BroadCastString(msg)
		}
		fmt.Print("> ")
	}
	log.Printf("Exiting newConnection")
}

func runReader(cons *Model) {
	for {
		msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		text := msg
		// TODO maybe check whether was processed?
		cons.MessageProcessed(text)
		go cons.BroadCastString(text)
	}
}

func runServer(cons *Model) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	ipAddress := GetOutboundIP()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	ipAndPort := ipAddress + ":" + port
	printStatus("Listening for connections on IP:port " + ipAndPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Connection failed - " + err.Error())
			continue
		}
		go newConnection(conn, cons)
	}
}

func startup() {
	cons := MakeModel()

	printStatus("Enter IP address and the port of the peer.")
	stdinReader := bufio.NewReader(os.Stdin)
	ipPort, _ := stdinReader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))

	if err != nil {
		printStatus("It was not possible to connect - creating own network.")
	} else {
		go newConnection(conn, &cons)
	}

	go runServer(&cons)
	go runReader(&cons)
	for {
	}
}

func main() {
	startup()
}
