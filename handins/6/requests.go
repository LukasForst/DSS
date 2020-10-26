package main

import (
	"crypto/rsa"
)

type PresentDto struct {
	Type string
	Data PeerId
}

type PeerId struct {
	// name of the server
	Address   string
	PublicKey rsa.PublicKey
}

func MakePresent(id PeerId) PresentDto {
	return PresentDto{Type: "present", Data: id}
}

type PeersListDto struct {
	Type string
	// names of the servers
	Data []PeerId
}

func MakePeersList(data []PeerId) PeersListDto {
	return PeersListDto{Type: "peers-list", Data: data}
}

type PeersRequestDto struct {
	Type string
}

func MakePeersRequest() PeersRequestDto {
	return PeersRequestDto{Type: "peers-request"}
}

// data transfer object for transaction
type TransactionDto struct {
	Type string
	Data SignedTransaction
}

func MakeTransactionDto(transaction SignedTransaction) TransactionDto {
	return TransactionDto{Type: "transaction", Data: transaction}
}

type SignedTransaction struct {
	ID        string // Any string
	From      string // A verification key coded as a string
	To        string // A verification key coded as a string
	Amount    int    // Amount to transfer
	Signature string // Potential signature coded as string
}

func (m *Model) SignedTransaction(t *SignedTransaction) {
	l := &m.ledger
	l.lock.Lock()
	defer l.lock.Unlock()

	if t.IsSignatureCorrect() {
		l.Accounts[t.From] -= t.Amount
		l.Accounts[t.To] += t.Amount
	}
}
