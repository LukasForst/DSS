package main

import (
	"crypto/sha256"
	"math/rand"
	"strconv"
)

func MakeGenesis(initialAccountStates *map[string]int) GenesisBlock {
	genesis := GenesisBlock{
		Hash:                 "",
		Seed:                 rand.Int(),
		InitialAccountStates: *initialAccountStates,
		NextBlocksHashes:     make([]string, 0, 1),
	}
	genesis.Hash = genesis.ComputeHash()
	return genesis
}

func (g *GenesisBlock) ComputeHash() string {
	msgHash := sha256.New()

	WriteToHashSafe(&msgHash, strconv.Itoa(g.Seed))
	for accountKey, value := range g.InitialAccountStates {
		WriteToHashSafe(&msgHash, accountKey)
		WriteToHashSafe(&msgHash, strconv.Itoa(value))
	}

	msgHashSum := msgHash.Sum(nil)
	return string(msgHashSum)
}

func (g *GenesisBlock) GetLongestChainLeaf(bc *BlockChain) (int, string) {
	currentMaxDepth := 0
	currentMaxHash := g.Hash

	for _, hash := range g.NextBlocksHashes {
		block := bc.Blocks[hash]
		depth, foundHash := block.GetLongestChainLeaf(bc, currentMaxDepth+1)
		if depth > currentMaxDepth {
			currentMaxDepth = depth
			currentMaxHash = foundHash
		}
	}

	return currentMaxDepth, currentMaxHash
}
