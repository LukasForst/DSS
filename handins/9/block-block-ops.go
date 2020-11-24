package main

import (
	"crypto/sha256"
	"encoding/base64"
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
	hash := sha256.New()

	WriteStringToHashSafe(&hash, b.PreviousBlockHash)
	for _, transaction := range b.Transactions {
		WriteBytesToHashSafe(&hash, transaction.ComputeHash())
		WriteStringToHashSafe(&hash, transaction.Signature)
	}
	for _, nextBlock := range b.NextBlocksHashes {
		WriteStringToHashSafe(&hash, nextBlock)
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (b *Block) VerifyHash(hashToVerify string) bool {
	thisHash := b.ComputeBase64Hash()
	return hashToVerify == thisHash
}
