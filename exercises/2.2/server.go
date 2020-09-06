package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// GetOutboundIP preferred outbound ip of this machine
// based on code taken from https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go/37382208#37382208
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	hostip, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		log.Fatal(err)
	}

	return hostip
}

// handleconnection receives incoming connections, prints their ip, then closes them.
func handleConnection(c *Connections, conn net.Conn) {
	defer conn.Close()
	defer c.RemoveConn(conn)

	log.Printf("Received a connection from %s\n", conn.RemoteAddr().String())
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Error: " + err.Error())
			break
		} else {
			log.Print("From Client:", msg)
			fmt.Fprintf(conn, "You wrote: %s", msg)
		}
	}
}

type Connections struct {
	m map[string]net.Conn
}

func (c *Connections) AddConn(conn net.Conn) {
	c.m[conn.RemoteAddr().String()] = conn
}

func (c *Connections) RemoveConn(conn net.Conn) {
	delete(c.m, conn.RemoteAddr().String())
}

func (c *Connections) SedToAll(msg []byte) {
	for k, v := range c.m {
		_, err := v.Write(msg)
		if err != nil {
			log.Printf("It was not possibel to send something to " + k + " -> " + err.Error())
		} else {
			log.Printf("Data sent to " + k)
		}
	}

}

//Run the "server" functionality
func runServer() {
	// leaving the port as ":0" allows go to choose an available port on the machine
	ln, err := net.Listen("tcp", ":0")
	//standard boilerplate for catching errors
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	//get outbound IP address
	ipAddress := GetOutboundIP()
	//get the port the listener is currently listening on
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	ipAndPort := ipAddress + ":" + port
	log.Println("Listening for connections on IP:port " + ipAndPort)

	cons := Connections{m: make(map[string]net.Conn)}

	//Loop to accept incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Connection failed - " + err.Error())
			continue
		}
		cons.SedToAll([]byte("New Client: " + conn.RemoteAddr().String()))
		cons.AddConn(conn)

		go handleConnection(&cons, conn)
	}
}

func main() {
	runServer()
}
