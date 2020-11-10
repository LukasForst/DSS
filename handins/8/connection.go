package main

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
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
			var peers PeerList
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
		case "signed-block":
			var block SignedSequencerBlock
			if err := json.Unmarshal(payload, &block); err != nil {
				log.Fatal("It was not possible to parse signed block object.")
			}
			OnBlockReceived(&block, model)
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

	// the transaction is already in waiting state
	if _, ok := model.transactionsWaiting[transactionID]; ok {
		return
	}
	// the transaction was already processed
	if processed, exists := model.transactionsProcessed[transactionID]; exists && processed {
		return
	}

	// store the transaction
	model.transactionsWaiting[transactionID] = transaction

	// propagate transaction
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

func OnPeersList(peers PeerList, model *Model) {
	model.AddNetworkPeers(peers.Peers)
	model.PrintPeers()
	model.sequencerPublicKey = &peers.SequencerPk
}

func OnPeersRequest(conn net.Conn, model *Model) {
	peers := MakePeersList(PeerList{Peers: model.GetPeersList(), SequencerPk: *model.sequencerPublicKey})
	// encode as json and send to other party
	if err := json.NewEncoder(conn).Encode(peers); err != nil {
		log.Println("It was not possible to send data.")
	}
}

func OnBlockReceived(block *SignedSequencerBlock, model *Model) {
	//check if block is valid/to be accepted

	blocknumber := block.Block.BlockNumber
	PrintStatus("Block " + strconv.Itoa(blocknumber) + " received.")

	if !block.IsSignatureCorrect(model.sequencerPublicKey) {
		PrintStatus("Signature is incorrect! Not processing.")
		return
	}

	//check if block number was seen before
	if model.blocksProcessed[blocknumber] == true {
		return
	}

	// check if blocknumber is last seen block +1
	if blocknumber != 0 && model.blocksProcessed[blocknumber-1] == false {
		PrintStatus("Missing block! Received block: " +
			strconv.Itoa(blocknumber) + " but the previous is missing!")
		return
	}

	// process the block
	model.blocksProcessed[blocknumber] = true
	for _, id := range block.Block.TransactionIds {
		// pick up the transaction
		transaction, _ := model.transactionsWaiting[id]
		// check if it was already processed or not
		if transaction == nil {
			continue
		}
		PrintStatus("Processing " + transaction.ID)
		// perform the transaction
		model.ledger.DoSignedTransaction(transaction)
		model.transactionsWaiting[id] = nil
		model.transactionsProcessed[id] = true
	}
	PrintStatus("Block processed.")

	// send the block to other nodes
	model.BroadCastJson(MakeSignedSequencerBlockDto(*block))
}
