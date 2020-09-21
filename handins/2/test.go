package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

// should be executed once we have some peers ready
func Server() {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	ipAddress := GetOutboundIP()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	ipAndPort := ipAddress + ":" + port

	PrintStatus("Listening for connections on IP:port " + ipAndPort)

	// listen
	for {
		conn, err := ln.Accept()
		if err != nil {
			PrintStatus("Connection failed -> " + err.Error())
			continue
		}
		decoder := json.NewDecoder(conn)
		type B struct {
			Type string
			Rest json.RawMessage
		}
		var objmap map[string]json.RawMessage
		err = decoder.Decode(&objmap)

		if err != nil {
			log.Fatal("err -> " + err.Error())
		}
		value, ok := objmap["Type"]
		if !ok {
			log.Fatal("No type present!")
		}
		var messageType string
		err = json.Unmarshal(value, &messageType)
		if err != nil {
			log.Fatal("No type!")
		}

		if messageType == "transaction" {
			println("Transaction!")
		} else if messageType == "command" {
			println("Command!")
		}
	}
}

func client() {
	PrintStatus("Enter IP address and the port of the peer.")
	stdinReader := bufio.NewReader(os.Stdin)
	ipPort, _ := stdinReader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()

	type M struct {
		Type string
		Rest int
	}
	m := M{"transaction", 10}
	log.Println(m)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(m)

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("sent!")
}

//
//func main() {
//	//Server()
//	client()
//}
