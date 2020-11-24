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

	//check if previous block exists
	if previousBlock != nil {
		previousBlock.NextBlocksHashes = append(previousBlock.NextBlocksHashes, previousBlock.Hash)
	} else {
		PrintStatus("The previous block does not exist!")
	}
	// todo maybe lock that?
	// todo maybe ensure that nextblocks are empty
	previousBlock.NextBlocksHashes = append(previousBlock.NextBlocksHashes, previousBlock.Hash)

}
