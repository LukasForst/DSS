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

// wrapper adding clock to the transaction
type TransactionWithClock struct {
	Transaction Transaction
	// TODO @Hannah, please fix clock
	Clock map[string]int
}

// data transfer object for transaction
type TransactionDto struct {
	Type string
	Data TransactionWithClock
}

func MakeTransactionDto(transaction TransactionWithClock) TransactionDto {
	return TransactionDto{Type: "transaction", Data: transaction}
}

func (l *Ledger) DoTransaction(t Transaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.Accounts[t.From] -= t.Amount
	l.Accounts[t.To] += t.Amount
}
