package main

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
