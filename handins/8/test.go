package main

import (
	"bufio"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

type Account struct {
	PK             *rsa.PrivateKey
	expectedAmount int
}

func GenerateAccounts(count int) []Account {
	accounts := make([]Account, 0, count)
	for i := 0; i < count; i++ {
		privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		acc := Account{privateKey, 0}
		accounts = append(accounts, acc)
	}
	return accounts
}

func StartTest() {
	stdinReader := bufio.NewReader(os.Stdin)
	PrintStatus("Enter unique test peer ID")
	testId, _ := stdinReader.ReadString('\n')
	testId = strings.TrimSpace(testId)

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

	transactions := 10
	accounts := GenerateAccounts(5)

	for i := 0; i < transactions; i++ {
		from := rand.Intn(len(accounts))
		to := rand.Intn(len(accounts))

		fromPk, _ := json.Marshal(accounts[from].PK.PublicKey)
		toPk, _ := json.Marshal(accounts[to].PK.PublicKey)

		amount := rand.Intn(100)
		accounts[from].expectedAmount = -amount
		accounts[to].expectedAmount = +amount

		transaction := SignedTransaction{
			ID:     testId + strconv.Itoa(i),
			From:   string(fromPk),
			To:     string(toPk),
			Amount: amount,
		}
		transaction.ComputeAndSetSignature(accounts[from].PK)

		dto := MakeTransactionDto(&transaction)
		if err := enc.Encode(dto); err != nil {
			log.Fatal("Error while sending data -> " + err.Error())
		}
	}
	// manually check how do the amounts look like
	log.Println("end")
}

//func main() {
//	StartTest()
//}
