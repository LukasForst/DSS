package main

import (
	"net"
	"sync"
)

type Model struct {
	connections  map[string]net.Conn
	cMutex       sync.RWMutex
	messagesSent map[string]bool
	mpMutex      sync.RWMutex
}

func MakeModel() Model {
	return Model{
		connections:  make(map[string]net.Conn),
		messagesSent: make(map[string]bool),
	}
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
	for k, v := range m.connections {
		_, err := v.Write(msg)
		if err != nil {
			PrintStatus("It was not possible to send something to " + k + " -> " + err.Error())
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

func (m *Model) BroadCastString(msg string) {
	m.BroadCast(msg)
}
