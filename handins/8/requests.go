package main

import "fmt"

type PresentDto struct {
	Type string
	// name of the server
	Data string
}

func MakePresent(address string) PresentDto {
	return PresentDto{Type: "present", Data: address}
}

type PeersListDto struct {
	Type string
	// names of the servers
	Data []string
}

func MakePeersList(data []string) PeersListDto {
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

func MakeTransactionDto(transaction *SignedTransaction) TransactionDto {
	return TransactionDto{Type: "transaction", Data: *transaction}
}

type SignedTransaction struct {
	ID        string // Any string
	From      string // A verification key coded as a string
	To        string // A verification key coded as a string
	Amount    int    // Amount to transfer
	Signature string // Potential signature coded as string
}

func (l *Ledger) DoSignedTransaction(t *SignedTransaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if !t.IsSignatureCorrect() {
		PrintStatus("Transaction " + t.ID + " has incorrect signature!")
	} else if t.Amount < 0 {
		PrintStatus("Transaction " + t.ID + " has incorrect negative amount!")
	} else if t.Amount > l.Accounts[t.From] {
		PrintStatus("Transaction " + t.ID + ": the transaction amount is higher than the account balance of the sender! ")
	} else {
		PrintStatus(fmt.Sprintf("Sender Balance: %d", l.Accounts[t.From]))
		PrintStatus(fmt.Sprintf("Recipient Balance: %d", l.Accounts[t.To]))

		l.Accounts[t.From] -= t.Amount
		l.Accounts[t.To] += t.Amount

		PrintStatus(fmt.Sprintf("Transaction %s performed - amount: %d", t.ID, t.Amount))
		PrintStatus(fmt.Sprintf("Sender Balance: %d", l.Accounts[t.From]))
		PrintStatus(fmt.Sprintf("Recipient Balance: %d", l.Accounts[t.To]))
	}
}
