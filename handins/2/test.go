package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

func StartTest() {
	stdinReader := bufio.NewReader(os.Stdin)
	PrintStatus("Enter unique test peer ID")
	testId, _ := stdinReader.ReadString('\n')

	PrintStatus("Enter IP address and the port of the peer.")
	ipPort, _ := stdinReader.ReadString('\n')

	// connect to the peer
	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))
	if err != nil {
		log.Fatal("It was not possible to connect.")
	}

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)
	// ask for peers
	if err := enc.Encode(MakePeersRequest()); err != nil {
		log.Fatal("It was not possible to connect to the first peer! ->" + err.Error())
	}
	// receive peers
	var peers PeersListDto
	if err := dec.Decode(&peers); err != nil {
		log.Fatal("It was not possible to connect to the first peer! ->" + err.Error())
	}

	PrintStatus("Pres enter to start the execution.")
	_, _ = stdinReader.ReadString('\n')

	transactions := 100
	for i := 0; i < transactions; i++ {
		from := rand.Intn(len(peers.Data))
		to := rand.Intn(len(peers.Data))

		transaction := Transaction{
			ID:     testId + strconv.Itoa(i),
			From:   peers.Data[from],
			To:     peers.Data[to],
			Amount: rand.Intn(100),
		}

		dto := MakeTransactionDto(transaction)
		if err := enc.Encode(dto); err != nil {
			log.Fatal("Error while sending data -> " + err.Error())
		}
	}
}

//func main() {
//	StartTest()
//}
