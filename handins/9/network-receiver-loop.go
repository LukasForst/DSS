package main

import (
	"encoding/json"
	"log"
	"net"
)

// process the connection
func OnNewConnection(conn net.Conn, model *PeerModel) {
	model.network.AddConn(conn)
	defer model.network.RemoveAndCloseConn(conn)

	decoder := json.NewDecoder(conn)
	for {
		PrintStatus("Waiting on message...")

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
			var block SignedBlock
			if err := json.Unmarshal(payload, &block); err != nil {
				log.Fatal("It was not possible to parse signed block object.")
			}
			OnBlockReceived(&block, model)
		case "account-setup":
			var setup AccountSetup
			if err := json.Unmarshal(payload, &setup); err != nil {
				log.Fatal("It was not possible to parse setup account object.")
			}
			OnAccountSetup(setup, model)
		}
	}
}
