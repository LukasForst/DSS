package main

import (
	"bufio"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Account struct {
	PK     *rsa.PrivateKey
	amount int
}

func GenerateAccount(amount int) Account {
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return Account{privateKey, amount}
}

func GetEncoder(id int, stdinReader *bufio.Reader) *net.Conn {
	PrintStatus("Enter IP address and the port of the peer # " + strconv.Itoa(id))
	ipPort, _ := stdinReader.ReadString('\n')
	// connect to the peer
	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))
	if err != nil {
		log.Fatal("It was not possible to connect.")
	}
	return &conn
}

func SetupAccounts(c *net.Conn, accounts []Account) {
	for _, account := range accounts {
		accountId, _ := json.Marshal(account.PK.PublicKey)
		dto := MakeAccountSetupDto(
			&AccountSetup{AccountId: string(accountId), Amount: account.amount},
		)
		SendJson(c, dto)
		// sleep otherwise the json encoder on the other side will break
		time.Sleep(50 * time.Millisecond)
	}
}

func SendJson(c *net.Conn, d interface{}) {
	bytes, err := json.Marshal(d)
	if err != nil {
		log.Fatal("Error while sending data -> " + err.Error())
	}

	_, err = (*c).Write(bytes)
	if err != nil {
		log.Fatal("Error while sending data -> " + err.Error())
	}

}

func RunTransaction(transactionId string, from *Account, to *Account, amount int, c *net.Conn) {
	fromPk, _ := json.Marshal(from.PK.PublicKey)
	toPk, _ := json.Marshal(to.PK.PublicKey)

	transaction := SignedTransaction{
		ID:     transactionId,
		From:   string(fromPk),
		To:     string(toPk),
		Amount: amount,
	}
	transaction.ComputeAndSetSignature(from.PK)

	dto := MakeTransactionDto(&transaction)
	SendJson(c, dto)
}

func StartExecution(even bool, testId string, from *Account, to *Account, peer *net.Conn) {
	for i := 0; i < from.amount*2; i++ {
		if (i%2 == 0 && even) || (i%2 == 1 && !even) {
			transactionId := fmt.Sprintf("%s-%04d", testId, i)
			RunTransaction(transactionId, from, to, 1, peer)
		}
	}
}

func StartTest() {
	stdinReader := bufio.NewReader(os.Stdin)
	PrintStatus("Enter unique test peer ID")
	testId, _ := stdinReader.ReadString('\n')
	testId = strings.TrimSpace(testId)

	peer1 := GetEncoder(1, stdinReader)
	peer2 := GetEncoder(2, stdinReader)

	A := GenerateAccount(20)
	B := GenerateAccount(0)
	C := GenerateAccount(0)
	accounts := []Account{A, B, C}

	SetupAccounts(peer1, accounts)
	SetupAccounts(peer2, accounts)

	PrintStatus("Pres enter to start the execution.")
	_, _ = stdinReader.ReadString('\n')

	go StartExecution(true, testId, &A, &B, peer1)
	go StartExecution(false, testId, &A, &C, peer2)

	// manually check how do the amounts look like
	PrintStatus("Pres enter to end the execution.")
	_, _ = stdinReader.ReadString('\n')
}

//
//func main() {
//	StartTest()
//}
