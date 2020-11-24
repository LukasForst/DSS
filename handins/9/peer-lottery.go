package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"math/big"
	"time"
)

func StartLotteryProcess(model *PeerModel) {
	sleepTime := model.blockChain.SlotLengthSeconds
	for {
		slot := model.blockChain.GetSlotNumber()
		peerAccountId := FromRsaPubToAccount(&model.peerKey.PublicKey)
		tokens := model.ledger.Accounts[peerAccountId]

		draw := RunLottery(
			slot,
			model.blockChain.Seed,
			big.NewInt(int64(tokens)),
			model.blockChain.Hardness,
			model.peerKey,
		)

		// won lottery
		if draw != nil {
			signedBlock := model.BuildAndExecuteSignedBlock(draw)
			transactionsCount := len(signedBlock.Block.Transactions)
			// reward for the block
			model.ledger.Accounts[peerAccountId] += transactionsCount + 10

			dto := MakeSignedBlockDto(signedBlock)
			model.network.BroadCastJson(dto)
		}

		time.Sleep(time.Duration(int64(sleepTime) * int64(time.Second)))
	}
}

func (pm *PeerModel) BuildAndExecuteSignedBlock(draw *Draw) *SignedBlock {
	block := pm.CreateAndExecuteBlock()
	blockHash := block.ComputeHash()

	signature, err := rsa.SignPSS(rand.Reader, pm.peerKey, crypto.SHA256, blockHash, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &SignedBlock{Block: *block, Draw: *draw, Signature: ToBase64(signature)}
}

func (bc *BlockChain) GetSlotNumber() int {
	return int(time.Now().Unix() / int64(bc.SlotLengthSeconds))
}
