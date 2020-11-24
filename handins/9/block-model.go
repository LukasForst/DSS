package main

import "sync"

type SignedBlock struct {
	Block     Block
	Signature []byte
}

type Block struct {
	Hash              string
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
