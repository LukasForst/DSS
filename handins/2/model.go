package main

import (
	"net"
	"sort"
	"strings"
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

func (m *Model) FormatMyPeers() string {
	// create sort list
	peers := make([]string, 0, len(m.peersList))
	for address, _ := range m.peersList {
		peers = append(peers, address)
	}
	sort.Strings(peers)
	return strings.Join(peers, ",")
}

func (m *Model) SelectTopNAfterMe(n int) []string {
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
	//TODO verify this
	fSize := Min(len(peers)-idx, idx+n)
	return peers[idx:fSize]
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

func (m *Model) BroadCast(text string) {
	msg := []byte(text)
	m.cMutex.RLock()
	defer m.cMutex.RUnlock()

	for peerAddress, connection := range m.connections {
		_, err := connection.Write(msg)
		if err != nil {
			PrintStatus("It was not possible to send something to " + peerAddress + " -> " + err.Error())
		}
	}
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
