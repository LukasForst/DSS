package main

import (
	"fmt"
	"log"
	"net"
)

func PrintIncoming(str string) {
	fmt.Printf("< %s\n> ", str)
}

func PrintStatus(str string) {
	fmt.Printf("- %s\n> ", str)
}

// GetOutboundIP preferred outbound ip of this machine
// based on code taken from https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go/37382208#37382208
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	hostip, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		log.Fatal(err)
	}

	return hostip
}

// ....
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
