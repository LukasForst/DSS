package main

import (
	"crypto/rsa"
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

	mpMutex sync.RWMutex

	// key: ip & port, value: is me?
	// key: ip & port, value: PeerId
	peersList map[string]PeerId
	pMutex    sync.RWMutex

	myAddress string

	privateKey *rsa.PrivateKey

	ledger Ledger
}

func MakeModel(pk *rsa.PrivateKey) Model {
	return Model{
		connections:      make(map[string]net.Conn),
		transactionsSeen: make(map[string]bool),
		peersList:        make(map[string]PeerId),
		ledger:           MakeLedger(),
		myAddress:        "",
		privateKey:       pk,
	}
}

// returns true if it was added, or already existed
func (m *Model) AddNetworkPeer(peerId PeerId) bool {
	m.pMutex.Lock()
	defer m.pMutex.Unlock()

	_, ok := m.peersList[peerId.Address]
	// if does not exist
	if !ok {
		// add new zero record to the ledger
		m.ledger.lock.Lock()
		defer m.ledger.lock.Unlock()
		// store key as json of public key of the peer
		peerPKBytes, err := json.Marshal(peerId.PublicKey)
		if err != nil {
			panic(err)
		}
		// set 0 at the beginning
		m.ledger.Accounts[string(peerPKBytes)] = 0
		m.peersList[peerId.Address] = peerId
	}
	return !ok
}

// during the protocol warmup, nobody knows my own address
func (m *Model) AddNetworkPeers(peerIds []PeerId) {
	for _, peerId := range peerIds {
		m.AddNetworkPeer(peerId)
	}
}

// store server address in the peers list
func (m *Model) RegisterMyAddress(address string) {
	// this is sequential, no need for locking
	m.myAddress = address
	m.AddNetworkPeer(PeerId{address, m.privateKey.PublicKey})
}

// returns sorted list of all peers
func (m *Model) GetPeersList() []PeerId {
	m.pMutex.RLock()
	defer m.pMutex.RUnlock()
	// create sort list
	peers := make([]PeerId, 0, len(m.peersList))
	for _, peerId := range m.peersList {
		peers = append(peers, peerId)
	}
	// TODO sort by address
	//sort.Strings(peers)
	return peers
}

// return array of all peers that are in the sorted list behind
// current instance
func (m *Model) SelectTopNAfterMe(n int) []PeerId {
	m.pMutex.RLock()
	defer m.pMutex.RUnlock()

	if len(m.peersList) == 1 {
		return make([]PeerId, 0, 0)
	}
	// remember who am I
	me := m.myAddress
	// create sort list
	peers := make([]string, 0, len(m.peersList))
	for address := range m.peersList {
		peers = append(peers, address)
	}
	sort.Strings(peers)
	// find me
	idx := sort.SearchStrings(peers, me)
	// select 10 after me, maybe less
	fSelection := make([]PeerId, 1, Min(len(peers)-idx, n))
	for i := 0; i < len(fSelection); i++ {
		fSelection[i] = m.peersList[peers[(i+idx+1)%len(peers)]]
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
	selection := make([]string, 0, len(m.peersList))
	for _, peerId := range m.GetPeersList() {
		selection = append(selection, peerId.Address)
	}
	PrintStatus("Peers: " + strings.Join(selection, ", "))
}
