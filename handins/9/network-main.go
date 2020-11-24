package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

func ConnectToNeighborhood(model *PeerModel) {
	// connect to the peers in the neighborhood
	peers := model.network.SelectNeighborhood(10)
	for _, peer := range peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			log.Fatal("It was not possible to connect to peer - " + err.Error())
		}
		go OnNewConnection(conn, model)
	}
}

// should be executed once we have some peers ready
func RunServer(model *PeerModel) {
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
	model.network.myIpPort = ipAndPort
	// connect to neighborhood
	ConnectToNeighborhood(model)
	// broadcast presence
	go model.network.BroadCastJson(MakePresent(ipAndPort))

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

// pull data from the remote peer
func InitialConnection(conn net.Conn, model *PeerModel) {
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	// ask for peers
	if err := enc.Encode(MakePeersRequest()); err != nil {
		log.Fatal("It was not possible to connect to the first peer! ->" + err.Error())
	}
	// receive peers
	var peers PeersListDto
	if err := dec.Decode(&peers); err != nil {
		log.Fatal("It was not possible to connect to the first peer! ->" + err.Error())
	}
	// register the peers in the application
	model.network.AddNetworkPeers(peers.Data.Peers)
	// print peers
	model.network.PrintPeers()
	// drop this initial connection
	_ = conn.Close()
}

func StartupServer() {
	// todo initialize model
	model := PeerModel{}

	PrintStatus("Enter IP address and the port of the peer.")
	stdinReader := bufio.NewReader(os.Stdin)
	ipPort, _ := stdinReader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))
	// if the peer does not exist, don't connect and start the network
	if err != nil {
		PrintStatus("It was not possible to connect - creating own network & starting the sequencer mode.")
	} else {
		// perform initial sync with the remote peer
		InitialConnection(conn, &model)
	}

	RunServer(&model)
}
