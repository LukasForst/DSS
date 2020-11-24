package main

import (
	"math/big"
	"sync"
)

type Draw struct {
	Slot      int
	Signature []byte
}

type SignedBlock struct {
	Block     Block
	Draw      Draw
	Signature string
}

type Block struct {
	Hash              string
	Epoch             int
	CreatorAccount    string
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
	Blocks       map[string]Block
	GenesisBlock GenesisBlock

	Hardness          *big.Int
	Seed              int
	SlotLengthSeconds int

	lock sync.Mutex
}
