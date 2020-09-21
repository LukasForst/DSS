package main

import (
	"encoding/json"
	"log"
	"net"
)

// process the connection
func OnNewConnection(conn net.Conn, model *Model) {
	model.AddConn(conn)
	defer model.RemoveAndCloseConn(conn)
	for {
		decoder := json.NewDecoder(conn)
		// read raw json and then decide what type is that
		var objmap map[string]json.RawMessage
		if err := decoder.Decode(&objmap); err != nil {
			PrintStatus("Connection dropped - " + err.Error())
			break
		}

		// obtain type of json
		mType, payload := MessageTypeAndRest(objmap)
		switch mType {
		// just send my own peers
		// not propagated in the network, no need to check for uniqueness
		case "peers-request":
			OnPeersRequest(conn, model)
		// new peers received
		case "peers-list":
			OnPeersList(payload, model)
		// notification about presence
		// need to check uniqueness, propagated on the network
		case "present":
			OnPresent(objmap, payload, model)
		// new transaction in the system
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

func OnPresent(
	objmap map[string]json.RawMessage,
	payload json.RawMessage,
	model *Model,
) {
	added := model.AddNetworkPeer(UnmarshalString(payload))
	if added {
		PrintStatus("Peer joined: " + UnmarshalString(payload))
		model.PrintPeers()
		// broadcast the presence further in the network
		model.BroadCastJson(objmap)
	}
}

func OnPeersList(payload json.RawMessage, model *Model) {
	var peers []string
	if err := json.Unmarshal(payload, &peers); err != nil {
		log.Fatal("It was not possible to unmarshall the peers list.")
	}
	model.AddNetworkPeers(peers)
	model.PrintPeers()
}

func OnPeersRequest(conn net.Conn, model *Model) {
	peers := MakePeersList(model.GetPeersList())
	// encode as json and send to other party
	if err := json.NewEncoder(conn).Encode(peers); err != nil {
		log.Println("It was not possible to send data.")
	}
}
