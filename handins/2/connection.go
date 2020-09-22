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
			var peers []string
			if err := json.Unmarshal(payload, &peers); err != nil {
				log.Fatal("It was not possible to unmarshall the peers list.")
			}
			OnPeersList(peers, model)
		// notification about presence
		// need to check uniqueness, propagated on the network
		case "present":
			OnPresent(objmap, UnmarshalString(payload), model)
		// new transaction in the system
		case "transaction":
			var transaction Transaction
			if err := json.Unmarshal(payload, &transaction); err != nil {
				log.Fatal("It was not possible to parse transaction object.")
			}
			OnTransactionReceived(transaction, model)
		}
	}
}

func OnTransactionReceived(transaction Transaction, model *Model) {
	//check whether we already did the transaction (check ID in transactionSeen map)

	transactionID := transaction.ID

	// If we already did the transaction, return
	if model.transactionsSeen[transactionID] == true {
		return
	} else {
		// lock the ledger
		model.ledger.lock.Lock()
		defer model.ledger.lock.Unlock()

		// perform the transaction
		model.ledger.DoTransaction(transaction)

		//register Transaction as seen
		model.transactionsSeen[transactionID] = true

		// propagate transaction
		model.BroadCastJson(MakeTransactionDto(transaction))

	}
}

func OnPresent(
	objmap map[string]json.RawMessage,
	newPeer string,
	model *Model,
) {
	added := model.AddNetworkPeer(newPeer)
	if added {
		PrintStatus("Peer joined: " + newPeer)
		model.PrintPeers()
		// broadcast the presence further in the network
		model.BroadCastJson(objmap)
	}
}

func OnPeersList(peers []string, model *Model) {
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
