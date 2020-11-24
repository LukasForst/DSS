package main

func (pm *PeerModel) GetBlock(blockId string) *Block {
	return pm.blockChain.Blocks[blockId]
}

// returns true if the transaction should be broadcast
func (pm *PeerModel) QueueTransactionIfNew(transaction *SignedTransaction) bool {
	_, exist := pm.hasBroadcastedTransaction[transaction.ID]
	if exist {
		return false
	} else {
		executed, exists := pm.hasExecutedTransaction[transaction.ID]
		if exists && executed {
			// ignore
		} else {
			pm.waitingTransactions[transaction.ID] = transaction
		}
		return false
	}
}

func (pm *PeerModel) ProcessBlock(signedBlock *SignedBlock) {
	// todo some locking

	// todo verify signature
	block := &signedBlock.Block
	pm.blockChain.AppendBlock(block)
	// happy scenario
	if block.PreviousBlockHash == pm.ledgerBlockId {
		pm.DoBlocksTransactions(block)
	} else {
		PrintStatus("Previous block is different than the current one.")
		_, longestPathBlockId := pm.blockChain.GetLongestChainLeaf()

		if longestPathBlockId == block.Hash {
			// reset all transactions until parent
			currentParentBlockId := block.PreviousBlockHash
			wrongId := pm.ledgerBlockId

			// undo all transactions
			for wrongId != currentParentBlockId {
				blockToUndo := pm.GetBlock(wrongId)
				// undo transactions
				pm.UndoBlocksTransactions(blockToUndo)
				// and run next iteration
				wrongId = pm.GetBlock(blockToUndo.PreviousBlockHash).Hash
			}

			// execute transactions from new leaf
			pm.DoBlocksTransactions(block)
		} else {
			PrintStatus("New block is not longest path, not executing transactions.")
		}
	}
}

func (pm *PeerModel) DoBlocksTransactions(block *Block) {
	for _, transaction := range block.Transactions {
		pm.DoLedgerTransaction(&transaction)
	}
	// set ledger identification
	pm.ledgerBlockId = block.Hash
}

func (pm *PeerModel) DoLedgerTransaction(transaction *SignedTransaction) {
	// todo locking
	pm.ledger.DoTransaction(transaction)

	pm.hasExecutedTransaction[transaction.ID] = true
	// delete transaction from waiting list
	_, exists := pm.waitingTransactions[transaction.ID]
	if exists {
		delete(pm.waitingTransactions, transaction.ID)
	}
}

func (pm *PeerModel) UndoBlocksTransactions(blockToUndo *Block) {
	// undo the transactions in reverse order
	for i := len(blockToUndo.Transactions) - 1; i >= 0; i-- {
		pm.UndoLedgerTransaction(&blockToUndo.Transactions[i])
	}
	// set ledger identification
	pm.ledgerBlockId = blockToUndo.PreviousBlockHash
}

func (pm *PeerModel) UndoLedgerTransaction(transaction *SignedTransaction) {
	// todo locking
	pm.ledger.UndoTransaction(transaction)

	pm.hasExecutedTransaction[transaction.ID] = false
	pm.waitingTransactions[transaction.ID] = transaction
}

func (pm *PeerModel) CreateAndExecuteBlock() *Block {
	// todo locking
	transactionsInBlock := make([]SignedTransaction, 0, 0)
	for _, transaction := range pm.waitingTransactions {
		transactionsInBlock = append(transactionsInBlock, *transaction)
	}

	_, previousBlock := pm.blockChain.GetLongestChainLeaf()
	block := Block{
		Hash:              "",
		PreviousBlockHash: previousBlock,
		Transactions:      transactionsInBlock,
		NextBlocksHashes:  make([]string, 0, 0),
	}
	// compute hash ~= id of the block
	block.Hash = block.ComputeBase64Hash()

	// execute transactions, removing them from the waiting list
	pm.blockChain.AppendBlock(&block)
	// execute transactions
	pm.DoBlocksTransactions(&block)

	return &block
}
