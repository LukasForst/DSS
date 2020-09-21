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
			var transaction TransactionWithClock
			if err := json.Unmarshal(payload, &transaction); err != nil {
				log.Fatal("It was not possible to parse transaction object.")
			}
			OnTransactionReceived(transaction, model)
		}
	}
}

func OnTransactionReceived(transaction TransactionWithClock, model *Model) {
	// TODO @Hannah - check whether we already did the transaction (use clock)

	for i := 0; i < (len(transaction.Clock)); i++ {
		//check clock counters
	}

	//transaction.Transaction.From

	// transaction.Clock
	// TODO @Hannah - if we already did the transaction, return

	// lock the ledger
	model.ledger.lock.Lock()
	defer model.ledger.lock.Unlock()

	// TODO @Hannah - check whether we can perform transaction right away (diff in clock is just one)
	// if this is not the case, store transaction in some waiting queue

	// TODO @Hannah - perform the transaction if it is safe (diff in clock just one)
	model.ledger.DoTransaction(transaction.Transaction)

	// propagate transaction with the same clock
	model.BroadCastJson(MakeTransactionDto(transaction))
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
