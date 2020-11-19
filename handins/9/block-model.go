package main

type SignedBlock struct {
	Block     Block
	Signature []byte
}

type Block struct {
	PreviousBlockHash string
	Transactions      []string
	NextBlocksHashes  []string
}

type BlockChain struct {
	Blocks map[string]Block
}

type Chain struct {
}
