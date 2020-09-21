package main

import (
	"sync"
)

type Present struct {
	Type string
	Data string
}

func MakePresent(address string) Present {
	return Present{Type: "present", Data: address}
}

type PeersList struct {
	Type string
	Data []string
}

func MakePeersList(data []string) PeersList {
	return PeersList{Type: "peers-list", Data: data}
}

type PeersRequest struct {
	Type string
}

func MakePeersRequest() PeersRequest {
	return PeersRequest{Type: "peers-request"}
}

type Ledger struct {
	Accounts map[string]int
	lock     sync.Mutex
}

func MakeLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]int)
	return ledger
}

type Transaction struct {
	ID     string
	From   string
	To     string
	Amount int
}

func (l *Ledger) DoTransaction(t *Transaction) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.Accounts[t.From] -= t.Amount
	l.Accounts[t.To] += t.Amount

}
