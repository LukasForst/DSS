package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

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
	go model.BroadCastJson(MakePresent(PeerId{ipAndPort, model.privateKey.PublicKey}))

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
func InitialConnection(conn net.Conn, model *Model) {
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
	model.AddNetworkPeers(peers.Data)
	// print peers
	model.PrintPeers()
	// drop this initial connection
	_ = conn.Close()
}

func Startup() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	model := MakeModel(privateKey)

	PrintStatus("Enter IP address and the port of the peer.")
	stdinReader := bufio.NewReader(os.Stdin)
	ipPort, _ := stdinReader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))
	// if the peer does not exist, don't connect
	if err != nil {
		PrintStatus("It was not possible to connect - creating own network.")
	} else {
		// perform initial sync with the remote peer
		InitialConnection(conn, &model)
	}

	RunServer(&model)
}

//
//func main() {
//	Startup()
//}
