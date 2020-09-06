package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func dial(addr string) net.Conn {
	conn, _ := net.Dial("tcp", addr)
	return conn
}

func sayHello(conn net.Conn) {
	n, err := conn.Write([]byte("Hello from the client!\n"))
	if err != nil {
		log.Fatal("It was not possible to write something! " + err.Error())
	} else {
		log.Printf("%d bytes sent\n", n)
	}
}

func readTcpLine(conn net.Conn) {
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal("Error during reading: " + err.Error())
		} else {
			log.Println("Server respond with: " + msg)
		}
	}
}

func sendFromStdio(conn net.Conn) {
	for {
		fmt.Print("> ")
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if strings.TrimSpace(text) == "quit" {
			log.Print("Ending transmission!")
			break
		}
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Fatal("Sending failed! " + err.Error())
		}
	}
}

func runClient() {
	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Println(">")
	ipPort, _ := stdinReader.ReadString('\n')

	conn := dial(strings.TrimSpace(ipPort))
	defer conn.Close()

	sayHello(conn)

	go readTcpLine(conn)
	go sendFromStdio(conn)
	for{}
}

func main() {
	runClient()
}
