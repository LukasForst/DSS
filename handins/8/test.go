package main

import (
	"bufio"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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

func GetEncoder(id int, stdinReader *bufio.Reader) *json.Encoder {
	PrintStatus("Enter IP address and the port of the peer # " + strconv.Itoa(id))
	ipPort, _ := stdinReader.ReadString('\n')
	// connect to the peer
	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))
	if err != nil {
		log.Fatal("It was not possible to connect.")
	}
	return json.NewEncoder(conn)
}

func SetupAccounts(enc *json.Encoder, accounts []Account) {
	for _, account := range accounts {
		accountId, _ := json.Marshal(account.PK.PublicKey)
		dto := MakeAccountSetupDto(AccountSetup{AccountId: string(accountId), Amount: account.amount})

		if err := enc.Encode(dto); err != nil {
			log.Fatal("Error while sending data -> " + err.Error())
		}
	}
}

func RunTransaction(transactionId string, from *Account, to *Account, amount int, enc *json.Encoder) {
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
	if err := enc.Encode(dto); err != nil {
		log.Fatal("Error while sending data -> " + err.Error())
	}
}

func StartTest() {
	stdinReader := bufio.NewReader(os.Stdin)
	PrintStatus("Enter unique test peer ID")
	testId, _ := stdinReader.ReadString('\n')
	testId = strings.TrimSpace(testId)

	peer1 := GetEncoder(1, stdinReader)
	peer2 := GetEncoder(2, stdinReader)

	A := GenerateAccount(1000)
	B := GenerateAccount(0)
	C := GenerateAccount(0)
	accounts := []Account{A, B, C}

	SetupAccounts(peer1, accounts)
	SetupAccounts(peer2, accounts)

	PrintStatus("Pres enter to start the execution.")
	_, _ = stdinReader.ReadString('\n')

	for i := 0; i < A.amount; i++ {
		go RunTransaction(testId+strconv.Itoa(i), &A, &B, 1, peer1)
		go RunTransaction(testId+strconv.Itoa(i), &A, &C, 1, peer2)
	}
	// manually check how do the amounts look like
	log.Println("end")
	PrintStatus("Pres enter to end the execution.")
	_, _ = stdinReader.ReadString('\n')
}

func main() {
	StartTest()
}
