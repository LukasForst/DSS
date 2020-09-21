package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

func OnNewConnection(conn net.Conn, model *Model) {
	model.AddConn(conn)
	defer model.RemoveAndCloseConn(conn)
	// first phase
	for {

		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}

		if model.WasProcessed(msg) {
			continue
		}
		model.MessageProcessed(msg)

		// todo maybe remove this
		PrintIncoming(strings.TrimSpace(msg))

		// todo this is ugly
		split := strings.Split(strings.TrimSpace(msg), "|")
		command := split[0]

		if command == "present" {
			payload := split[1]
			model.AddNetworkPeer(payload)
			go model.BroadCast(msg)
		} else if command == "send-peers" {
			peers := model.FormatMyPeers()
			_, err := conn.Write([]byte(peers))
			if err != nil {
				log.Fatal("Sending peers failed.")
			}
		} else if command == "transactions" {
			PrintStatus("Ending connection phase.")
			break
		}
	}
}

func ConnectToNeighborhood(model *Model) {
	// connect to the peers in the neighborhood
	peers := model.SelectTopNAfterMe(10)
	for _, peer := range peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			log.Fatal("It was not possible to connect to peer - " + err.Error())
		}
		go OnNewConnection(conn, model)
	}
}

// should be executed once we have some peers ready
func RunServer(model *Model) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	ipAddress := GetOutboundIP()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	ipAndPort := ipAddress + ":" + port

	PrintStatus("Listening for connections on IP:port " + ipAndPort)
	// register my own address in the address book
	model.RegisterMyAddress(ipAndPort)
	// connect to neighborhood
	ConnectToNeighborhood(model)
	// broadcast presence
	go model.BroadCast("present|" + ipAndPort + "\n")

	// listen
	for {
		conn, err := ln.Accept()
		if err != nil {
			PrintStatus("Connection failed -> " + err.Error())
			continue
		}
		go OnNewConnection(conn, model)
	}
}

func InitialConnection(conn net.Conn, model *Model) {
	defer conn.Close()

	// ask for peers
	_, err := conn.Write([]byte("send-peers\n"))
	if err != nil {
		log.Fatal("It was not possible to connect to the first peer! ->" + err.Error())
	}
	// receive peers
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal("It was not possible to connect to the first peer! ->" + err.Error())
	}
	// register the peers in the application
	peers := strings.Split(strings.TrimSpace(msg), ",")
	model.AddNetworkPeers(peers)
	PrintStatus("Peers list registered.")
}

func Startup() {
	model := MakeModel()

	PrintStatus("Enter IP address and the port of the peer.")
	stdinReader := bufio.NewReader(os.Stdin)
	ipPort, _ := stdinReader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))

	if err != nil {
		PrintStatus("It was not possible to connect - creating own network.")
	} else {
		InitialConnection(conn, &model)
	}

	go RunServer(&model)
	for {
	}
}

func main() {
	Startup()
}
