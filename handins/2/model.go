package main

import (
	"encoding/json"
	"log"
	"net"
	"sort"
	"sync"
)

type Model struct {
	connections map[string]net.Conn
	cMutex      sync.RWMutex

	messagesSent map[string]bool
	mpMutex      sync.RWMutex

	// key: ip & port, value: is me?
	peersList map[string]bool
	pMutex    sync.RWMutex
}

func MakeModel() Model {
	return Model{
		connections:  make(map[string]net.Conn),
		messagesSent: make(map[string]bool),
		peersList:    make(map[string]bool),
	}
}

func (m *Model) AddNetworkPeer(address string) {
	value, ok := m.peersList[address]
	m.peersList[address] = ok && value
}

// during the protocol warmup, nobody knows my own address
func (m *Model) AddNetworkPeers(addresses []string) {
	for _, address := range addresses {
		m.AddNetworkPeer(address)
	}
}

func (m *Model) RegisterMyAddress(address string) {
	m.peersList[address] = true
}

func (m *Model) GetPeersList() []string {
	// create sort list
	peers := make([]string, 0, len(m.peersList))
	for address, _ := range m.peersList {
		peers = append(peers, address)
	}
	sort.Strings(peers)
	return peers
}

func (m *Model) SelectTopNAfterMe(n int) []string {
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

	//TODO verify this
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

func (m *Model) BroadCastBytes(msg []byte) {
	m.cMutex.RLock()
	defer m.cMutex.RUnlock()

	for peerAddress, connection := range m.connections {
		_, err := connection.Write(msg)
		if err != nil {
			PrintStatus("It was not possible to send something to " + peerAddress + " -> " + err.Error())
		}
	}
}

func (m *Model) BroadCastJson(v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Fatal("It was not possible to serialize data.")
	}
	go m.BroadCastBytes(bytes)

}

func (m *Model) BroadCast(text string) {
	m.BroadCastBytes([]byte(text))
}

func (m *Model) WasProcessed(msg string) bool {
	m.mpMutex.RLock()
	defer m.mpMutex.RUnlock()

	value, ok := m.messagesSent[msg]
	return ok && value
}

func (m *Model) MessageProcessed(msg string) {
	m.mpMutex.Lock()
	defer m.mpMutex.Unlock()
	m.messagesSent[msg] = true
}
