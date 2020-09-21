package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

func OnNewConnection(conn net.Conn, model *Model) {
	model.AddConn(conn)
	defer model.RemoveAndCloseConn(conn)
	for {
		decoder := json.NewDecoder(conn)

		var objmap map[string]json.RawMessage
		if err := decoder.Decode(&objmap); err != nil {
			PrintStatus("Connection dropped - " + err.Error())
			break
		}

		mType, payload := MessageTypeAndRest(objmap)

		switch mType {
		// just send my own peers
		//not propagated in the network, no need to check for uniqueness
		case "peers-request":
			PrintStatus("Peers request received!")
			peers := MakePeersList(model.GetPeersList())
			if err := json.NewEncoder(conn).Encode(peers); err != nil {
				log.Println("It was not possible to send data.")
			}
		// notification about presence
		// need to check uniqueness, propagated on the network
		case "present":
			PrintStatus("Present status received from: " + UnmarshalString(payload))
			added := model.AddNetworkPeer(UnmarshalString(payload))
			if added {
				model.PrintPeers()
				// broadcast the presence further in the network
				model.BroadCastJson(objmap)
			}
		case "transaction":
			var transaction Transaction
			var ledger Ledger
			if err := json.Unmarshal(payload, &transaction); err != nil {
				log.Fatal("It was not possible to parse transaction.")
			}
			OnTransactionReceived(transaction, ledger, model)
		}
	}
}

func OnTransactionReceived(transaction Transaction, ledger Ledger, model *Model) {
	//TODO: make transaction, broadcast Transaction object, update local Ledger object

	//ledger.DoTransaction(transaction)
	model.BroadCastTransaction(transaction)

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
	go model.BroadCastJson(MakePresent(ipAndPort))

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
	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	PrintStatus("Asking for the peers.")
	// ask for peers
	if err := enc.Encode(MakePeersRequest()); err != nil {
		log.Fatal("It was not possible to connect to the first peer! ->" + err.Error())
	}
	PrintStatus("Receiving peers.")
	// receive peers
	var peers PeersList
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
	model := MakeModel()

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

func main() {
	Startup()
}
