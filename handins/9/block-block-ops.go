package main

import (
	"crypto/sha256"
	"strconv"
)

func (b *Block) GetLongestChainLeaf(bc *BlockChain, currentDepth int) (int, string) {
	currentMaxDepth := currentDepth
	currentMaxHash := b.Hash

	for _, hash := range b.NextBlocksHashes {
		block := bc.Blocks[hash]
		depth, foundHash := block.GetLongestChainLeaf(bc, currentDepth+1)
		if depth > currentMaxDepth {
			currentMaxDepth = depth
			currentMaxHash = foundHash
		}
	}

	return currentMaxDepth, currentMaxHash
}

func (b *Block) ComputeBase64Hash() string {
	return ToBase64(b.ComputeHash())
}

func (b *Block) ComputeHash() []byte {
	hash := sha256.New()

	WriteStringToHashSafe(&hash, b.PreviousBlockHash)
	WriteStringToHashSafe(&hash, strconv.Itoa(b.Epoch))
	for _, transaction := range b.Transactions {
		WriteStringToHashSafe(&hash, transaction.ComputeBase64Hash())
		WriteStringToHashSafe(&hash, transaction.Signature)
	}
	for _, nextBlock := range b.NextBlocksHashes {
		WriteStringToHashSafe(&hash, nextBlock)
	}
	WriteStringToHashSafe(&hash, b.CreatorAccount)

	return hash.Sum(nil)
}

func (b *Block) VerifyHash(hashToVerify string) bool {
	thisHash := b.ComputeBase64Hash()
	return hashToVerify == thisHash
}
