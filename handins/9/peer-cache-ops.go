package main

func (c *OperationCache) ShouldBroadcastTransaction(transaction *SignedTransaction) bool {
	return RunCacheOps(&c.hasBroadcastTransaction, transaction.ID)
}

func (c *OperationCache) ShouldBroadcastPeerJoined(peerIdJoined string) bool {
	return RunCacheOps(&c.hasBroadcastPeerJoined, peerIdJoined)
}

func (c *OperationCache) ShouldBroadcastBlock(block *SignedBlock) bool {
	return RunCacheOps(&c.hasBroadcastBlock, block.Signature)
}

func RunCacheOps(m *map[string]bool, key string) bool {
	_, exists := (*m)[key]
	if !exists {
		(*m)[key] = true
		return true
	} else {
		return false
	}
}
