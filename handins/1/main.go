package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func OnNewConnection(conn net.Conn, cons *Model) {
	cons.AddConn(conn)
	defer cons.RemoveAndCloseConn(conn)

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}

		if !cons.WasProcessed(msg) {
			cons.MessageProcessed(msg)
			PrintIncoming(strings.TrimSpace(msg))
			go cons.BroadCastString(msg)
		}
	}
}

func RunReader(cons *Model) {
	for {
		msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if strings.TrimSpace(msg) == "quit" {
			break
		}
		cons.MessageProcessed(msg)
		go cons.BroadCastString(msg)
		fmt.Print("> ")
	}
}

func RunServer(cons *Model) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	ipAddress := GetOutboundIP()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	ipAndPort := ipAddress + ":" + port
	PrintStatus("Listening for connections on IP:port " + ipAndPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			PrintStatus("Connection failed - " + err.Error())
			continue
		}
		go OnNewConnection(conn, cons)
	}
}

func Startup() {
	cons := MakeModel()

	PrintStatus("Enter IP address and the port of the peer.")
	stdinReader := bufio.NewReader(os.Stdin)
	ipPort, _ := stdinReader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))

	if err != nil {
		PrintStatus("It was not possible to connect - creating own network.")
	} else {
		go OnNewConnection(conn, &cons)
	}

	go RunServer(&cons)
	RunReader(&cons)
}

func main() {
	Startup()
}
