package main

import (
	"crypto/rsa"
	"net"
	"sync"
)

type PeerLedger struct {
	Accounts map[string]int
	lock     sync.Mutex
}

type PeerModel struct {
	blockChain *BlockChain

	connections map[string]*net.Conn
	cMutex      sync.RWMutex

	//key: id, value transaction
	// transactions to be put into block, if not seen
	waitingTransactions map[string]*SignedTransaction

	hasBroadcastedTransaction map[string]bool
	hasExecutedTransaction    map[string]bool

	// broadcas to other parts of the system
	hasSeenTransaction map[string]bool

	//key: transaction ID, value: seen before?
	transactionsProcessed map[string]bool
	mpMutex               sync.RWMutex

	// key: ip & port, value: is me?
	peersList map[string]bool
	pMutex    sync.RWMutex

	ledger *PeerLedger
	// block that is currently on the top of the ledger
	ledgerBlockId string

	pk *rsa.PrivateKey
}
