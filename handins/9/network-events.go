package main

import (
	"encoding/json"
	"log"
	"net"
)

func OnAccountSetup(setup AccountSetup, model *PeerModel) {
	model.ledger.lock.Lock()
	defer model.ledger.lock.Unlock()

	model.ledger.Accounts[setup.AccountId] = setup.Amount
}

func OnTransactionReceived(transaction *SignedTransaction, model *PeerModel) {
	//check whether we already did the transaction by checking the transaction ID

	transactionID := transaction.ID
	PrintStatus("Transaction " + transactionID + " received.")

	if model.cache.ShouldBroadcastTransaction(transaction) {

		model.QueueTransactionIfNew(transaction)

		model.network.BroadCastJson(MakeTransactionDto(transaction))
	}
}

func OnPresent(
	objmap map[string]json.RawMessage,
	newPeer string,
	model *PeerModel,
) {

	added := model.network.AddNetworkPeer(newPeer)
	if added {
		PrintStatus("Peer joined: " + newPeer)
		model.network.PrintPeers()
	}

	if model.cache.ShouldBroadcastPeerJoined(newPeer) {
		model.network.BroadCastJson(objmap)
	}
}

func OnPeersList(peers PeerList, model *PeerModel) {
	model.network.AddNetworkPeers(peers.Peers)
	model.network.PrintPeers()
}

func OnPeersRequest(conn net.Conn, model *PeerModel) {
	list := PeerList{Peers: model.network.GetPeersList()}
	peers := MakePeersList(&list)
	// encode as json and send to other party
	if err := json.NewEncoder(conn).Encode(peers); err != nil {
		log.Println("It was not possible to send data.")
	}
}

func OnBlockReceived(block *SignedBlock, model *PeerModel) {
	// send the block to other nodes
	if model.cache.ShouldBroadcastBlock(block) {
		// if not seen previously, process the block
		model.ProcessBlock(block)

		model.network.BroadCastJson(MakeSignedBlockDto(block))
	}
}
