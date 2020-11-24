package main

import (
	"crypto/rsa"
	"net"
	"sync"
)

type OperationCache struct {
	hasBroadcastTransaction map[string]bool
	hasBroadcastBlock       map[string]bool
	hasBroadcastPeerJoined  map[string]bool
	hasExecutedTransaction  map[string]bool
}

type PeerLedger struct {
	Accounts map[string]int
	lock     sync.Mutex
}

type Network struct {
	connections map[string]net.Conn
	cMutex      sync.RWMutex

	myIpPort string

	peersList map[string]bool
	pMutex    sync.RWMutex
}

type PeerModel struct {
	network *Network

	blockChain *BlockChain

	ledger *PeerLedger
	// block that is currently on the top of the ledger
	ledgerBlockId string

	cache *OperationCache

	// transactions to be put into block, if not seen
	waitingTransactions map[string]*SignedTransaction

	peerKey *rsa.PrivateKey
}
