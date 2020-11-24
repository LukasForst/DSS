package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
)

func FromRsaPubToAccount(key *rsa.PublicKey) string {
	res, err := json.Marshal(key)
	if err != nil {
		log.Fatal(err)
	}
	return string(res)
}

func (l *PeerLedger) DoTransaction(t *SignedTransaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if !t.IsSignatureCorrect() {
		PrintStatus("Transaction " + t.ID + " has incorrect signature!")
	} else if t.Amount < 1 {
		PrintStatus("Transaction " + t.ID + " has incorrect amount below 1!")
	} else if t.Amount > l.Accounts[t.From] {
		PrintStatus("Transaction " + t.ID + ": the transaction amount is higher than the account balance of the sender! ")
	} else {
		PrintStatus(fmt.Sprintf("Sender Balance: %d", l.Accounts[t.From]))
		PrintStatus(fmt.Sprintf("Recipient Balance: %d", l.Accounts[t.To]))

		l.Accounts[t.From] -= t.Amount
		l.Accounts[t.To] += t.Amount - 1

		PrintStatus(fmt.Sprintf("Transaction %s performed - amount: %d", t.ID, t.Amount))
		PrintStatus(fmt.Sprintf("Sender Balance: %d", l.Accounts[t.From]))
		PrintStatus(fmt.Sprintf("Recipient Balance: %d", l.Accounts[t.To]))
	}
}

func (l *PeerLedger) UndoTransaction(t *SignedTransaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if !t.IsSignatureCorrect() {
		PrintStatus("Transaction " + t.ID + " has incorrect signature!")
	} else if t.Amount < 1 {
		PrintStatus("Transaction " + t.ID + " has incorrect amount below 1!")
	} else if t.Amount > l.Accounts[t.From] {
		PrintStatus("Transaction " + t.ID + ": the transaction amount is higher than the account balance of the sender! ")
	} else {
		PrintStatus(fmt.Sprintf("Sender Balance: %d", l.Accounts[t.From]))
		PrintStatus(fmt.Sprintf("Recipient Balance: %d", l.Accounts[t.To]))

		l.Accounts[t.From] += t.Amount
		l.Accounts[t.To] -= t.Amount + 1

		PrintStatus(fmt.Sprintf("Transaction %s performed - amount: %d", t.ID, t.Amount))
		PrintStatus(fmt.Sprintf("Sender Balance: %d", l.Accounts[t.From]))
		PrintStatus(fmt.Sprintf("Recipient Balance: %d", l.Accounts[t.To]))
	}
}
