package main

import (
	"encoding/json"
	"log"
	"net"
	"sort"
	"strings"
	"sync"
)

type Ledger struct {
	Accounts map[string]int
	lock     sync.Mutex
}

func MakeLedger() Ledger {
	return Ledger{Accounts: make(map[string]int)}
}

type Model struct {
	connections map[string]net.Conn
	cMutex      sync.RWMutex

	//key: transaction ID, value: seen before?
	transactionsSeen map[string]bool
	mpMutex          sync.RWMutex

	// key: ip & port, value: is me?
	peersList map[string]bool
	pMutex    sync.RWMutex

	ledger Ledger

	blocksSeen       map[string]bool
	transactionsWait map[string]bool
}

func MakeModel() Model {
	return Model{
		connections:      make(map[string]net.Conn),
		transactionsWait: make(map[string]bool),
		transactionsSeen: make(map[string]bool),
		peersList:        make(map[string]bool),
		ledger:           MakeLedger(),
		blocksSeen:       make(map[string]bool),
	}
}

// returns true if it was added, or already existed
func (m *Model) AddNetworkPeer(address string) bool {
	m.pMutex.Lock()
	defer m.pMutex.Unlock()

	value, ok := m.peersList[address]
	m.peersList[address] = ok && value
	return !ok
}

// during the protocol warmup, nobody knows my own address
func (m *Model) AddNetworkPeers(addresses []string) {
	for _, address := range addresses {
		m.AddNetworkPeer(address)
	}
}

// store server address in the peers list
func (m *Model) RegisterMyAddress(address string) {
	// this is sequential, no need for locking
	m.peersList[address] = true
}

// returns sorted list of all peers
func (m *Model) GetPeersList() []string {
	m.pMutex.RLock()
	defer m.pMutex.RUnlock()
	// create sort list
	peers := make([]string, 0, len(m.peersList))
	for address := range m.peersList {
		peers = append(peers, address)
	}
	sort.Strings(peers)
	return peers
}

// return array of all peers that are in the sorted list behind
// current instance
func (m *Model) SelectTopNAfterMe(n int) []string {
	m.pMutex.RLock()
	defer m.pMutex.RUnlock()

	if len(m.peersList) == 1 {
		return make([]string, 0, 0)
	}
	// remember who am I
	var me string
	// create sort list
	peers := make([]string, 0, len(m.peersList))
	for address, isMe := range m.peersList {
		if isMe {
			me = address
		}
		peers = append(peers, address)
	}
	sort.Strings(peers)
	// find me
	idx := sort.SearchStrings(peers, me)
	// select 10 after me, maybe less
	fSelection := make([]string, 1, Min(len(peers)-idx, n))
	for i := 0; i < len(fSelection); i++ {
		fSelection[i] = peers[(i+idx+1)%len(peers)]
	}
	return fSelection
}

func (m *Model) AddConn(conn net.Conn) {
	m.cMutex.Lock()
	defer m.cMutex.Unlock()
	m.connections[conn.RemoteAddr().String()] = conn
}

func (m *Model) RemoveAndCloseConn(conn net.Conn) {
	delete(m.connections, conn.RemoteAddr().String())
	_ = conn.Close()
}

// send bytes to all current connections
func (m *Model) BroadCastBytes(bytes []byte) {
	m.cMutex.RLock()
	defer m.cMutex.RUnlock()

	for peerAddress, connection := range m.connections {
		_, err := connection.Write(bytes)
		if err != nil {
			PrintStatus("It was not possible to send something to " + peerAddress + " -> " + err.Error())
		}
	}
}

// encode as JSON and sent it to all connections
func (m *Model) BroadCastJson(v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Fatal("It was not possible to serialize data.")
	}
	m.BroadCastBytes(bytes)
}

// print all peers to the console
func (m *Model) PrintPeers() {
	PrintStatus("Peers: " + strings.Join(m.GetPeersList(), ", "))
}
