package main

import (
	"encoding/json"
	"log"
	"net"
	"sort"
	"strings"
)

func (n *Network) AddConn(conn net.Conn) {
	n.cMutex.Lock()
	defer n.cMutex.Unlock()
	n.connections[conn.RemoteAddr().String()] = conn
}

func (n *Network) RemoveAndCloseConn(conn net.Conn) {
	delete(n.connections, conn.RemoteAddr().String())
	_ = conn.Close()
}

// send bytes to all current connections
func (n *Network) BroadCastBytes(bytes []byte) {
	n.cMutex.RLock()
	defer n.cMutex.RUnlock()

	for peerAddress, connection := range n.connections {
		_, err := connection.Write(bytes)
		if err != nil {
			PrintStatus("It was not possible to send something to " + peerAddress + " -> " + err.Error())
		}
	}
}

// encode as JSON and sent it to all connections
func (n *Network) BroadCastJson(v interface{}) {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Fatal("It was not possible to serialize data.")
	}
	n.BroadCastBytes(bytes)
}

// print all peers to the console
func (n *Network) PrintPeers() {
	PrintStatus("Peers: " + strings.Join(n.GetPeersList(), ", "))
}

// returns sorted list of all peers
func (n *Network) GetPeersList() []string {
	n.pMutex.RLock()
	defer n.pMutex.RUnlock()
	// create sort list
	peers := make([]string, 0, len(n.peersList))
	for address := range n.peersList {
		peers = append(peers, address)
	}
	sort.Strings(peers)
	return peers
}

// returns sorted list of all peers
func (n *Network) SelectNeighborhood(size int) []string {
	return n.SelectTopNAfterMe(size)
}

// returns true if it was added, or already existed
func (n *Network) AddNetworkPeer(address string) bool {
	n.pMutex.Lock()
	defer n.pMutex.Unlock()

	value, ok := n.peersList[address]
	n.peersList[address] = ok && value
	return !ok
}

// during the protocol warmup, nobody knows my own address
func (n *Network) AddNetworkPeers(addresses []string) {
	for _, address := range addresses {
		n.AddNetworkPeer(address)
	}
}

// return array of all peers that are in the sorted list behind
// current instance
func (n *Network) SelectTopNAfterMe(size int) []string {
	n.pMutex.RLock()
	defer n.pMutex.RUnlock()

	if len(n.peersList) == 1 {
		return make([]string, 0, 0)
	}
	// remember who am I
	var me string
	// create sort list
	peers := make([]string, 0, len(n.peersList))
	for address, isMe := range n.peersList {
		if isMe {
			me = address
		}
		peers = append(peers, address)
	}
	sort.Strings(peers)
	// find me
	idx := sort.SearchStrings(peers, me)
	// select 10 after me, maybe less
	fSelection := make([]string, 1, Min(len(peers)-idx, size))
	for i := 0; i < len(fSelection); i++ {
		fSelection[i] = peers[(i+idx+1)%len(peers)]
	}
	return fSelection
}
