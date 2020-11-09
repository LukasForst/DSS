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
			var transaction SignedTransaction
			if err := json.Unmarshal(payload, &transaction); err != nil {
				log.Fatal("It was not possible to parse transaction object.")
			}
			OnTransactionReceived(&transaction, model)
		}
	}
}

func OnTransactionReceived(transaction *SignedTransaction, model *Model) {
	//check whether we already did the transaction by checking the transaction ID

	transactionID := transaction.ID
	PrintStatus("Transaction " + transactionID + " received.")
	// If we already did the transaction or it is already stored, return
	model.mpMutex.Lock()
	defer model.mpMutex.Unlock()
	if model.transactionsSeen[transactionID] == true || model.transactionsWait[transactionID] == true {
		return
	} else {
		// store the transaction
		model.transactionsWait[transactionID] = true

		// Wait for block from sequencer
		//OnBlockReceived: check if block is valid/to be accepted
		// If valid, execute transactions as in order from block

		//if valid: do te below for every transaction in the block in order

		// perform the transaction
		model.ledger.DoSignedTransaction(transaction)

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

func OnBlockReceived(block *Block, model *Model) {
	//check if block is valid/to be accepted

	blocknumber := block.number

	PrintStatus("Block " + blocknumber + " received.")

	//check if block number was seen before
	if model.blocksSeen[blocknumber] == true {
		return
	} else {
		//check if blocknumber is last seen block +1
		if model.blocksSeen[blocknumber-1] == false {
			return
		} else {
			//check block signature

		}
	}
}
