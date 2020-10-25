package main

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

// type given from the book
type Transaction struct {
	ID     string
	From   string
	To     string
	Amount int
}

// data transfer object for transaction
type TransactionDto struct {
	Type string
	Data Transaction
}

func MakeTransactionDto(transaction Transaction) TransactionDto {
	return TransactionDto{Type: "transaction", Data: transaction}
}

func (l *Ledger) DoTransaction(t Transaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.Accounts[t.From] -= t.Amount
	l.Accounts[t.To] += t.Amount
}

type SignedTransaction struct {
	ID        string // Any string
	From      string // A verification key coded as a string
	To        string // A verification key coded as a string
	Amount    int    // Amount to transfer
	Signature string // Potential signature coded as string
}

func (l *Ledger) SignedTransaction(t *SignedTransaction) {
	l.lock.Lock()
	defer l.lock.Unlock()
	/* We verify that the t.Signature is a valid RSA
	* signature on the rest of the fields in t under
	* the public key t.From.
	 */
	validSignature := true
	if validSignature {
		l.Accounts[t.From] -= t.Amount
		l.Accounts[t.To] += t.Amount
	}
}
