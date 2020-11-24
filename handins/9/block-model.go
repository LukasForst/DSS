package main

import "sync"

type Draw struct {
	Slot      int
	Signature []byte
}

type SignedBlock struct {
	Block     Block
	Signature string
}

type Block struct {
	Hash              string
	Epoch             int
	PreviousBlockHash string
	Transactions      []SignedTransaction
	NextBlocksHashes  []string
}

type GenesisBlock struct {
	Hash                 string
	Seed                 int
	InitialAccountStates map[string]int
	NextBlocksHashes     []string
}

type BlockChain struct {
	Blocks       map[string]*Block
	GenesisBlock GenesisBlock

	lock sync.Mutex
}
