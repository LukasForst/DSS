package main

func (bc *BlockChain) GetLongestChainLeaf() (int, string) {
	return bc.GenesisBlock.GetLongestChainLeaf(bc)
}

func (bc *BlockChain) AppendBlock(block *Block) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	bc.Blocks[block.Hash] = block
	// todo what if that one does not exist?
	previousBlock := bc.Blocks[block.PreviousBlockHash]

	if block.PreviousBlockHash != "" {
		previousBlock.NextBlocksHashes = append(previousBlock.NextBlocksHashes, previousBlock.Hash)
	} else {
		PrintStatus("The previous block hash is empty!")
	}
	// todo maybe lock that?
	// todo maybe ensure that nextblocks are empty
	//if len(previousBlock.NextBlocksHashes) == 0 {
	previousBlock.NextBlocksHashes = append(previousBlock.NextBlocksHashes, previousBlock.Hash)
	// } else {
	// 	PrintStatus("")
	// }
}
