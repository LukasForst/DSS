package main

import (
	"log"
	"net"
)

type Model struct {
	connections       map[string]net.Conn
	messagesProcessed map[string]bool
}

func MakeModel() Model {
	return Model{
		connections:       make(map[string]net.Conn),
		messagesProcessed: make(map[string]bool),
	}
}

func (m *Model) AddConn(conn net.Conn) {
	log.Println("Adding new connection: " + conn.RemoteAddr().String())
	m.connections[conn.RemoteAddr().String()] = conn
}

func (m *Model) RemoveAndCloseConn(conn net.Conn) {
	log.Println("Removing connection: " + conn.RemoteAddr().String())
	delete(m.connections, conn.RemoteAddr().String())
	conn.Close()
}

func (m *Model) BroadCast(text string) {
	msg := []byte(text)
	for k, v := range m.connections {
		_, err := v.Write(msg)
		if err != nil {
			log.Printf("It was not possibel to send something to " + k + " -> " + err.Error())
		}
	}
}

func (m *Model) WasProcessed(msg string) bool {
	value, ok := m.messagesProcessed[msg]
	return ok && value
}

func (m *Model) MessageProcessed(msg string) {
	m.messagesProcessed[msg] = true
}

func (m *Model) BroadCastString(msg string) {
	m.BroadCast(msg)
}
