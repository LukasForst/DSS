package main

import (
	"bufio"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ConnectToNeighborhood(model *Model) {
	// connect to the peers in the neighborhood
	peers := model.SelectTopNAfterMe(10)
	for _, peer := range peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			log.Fatal("It was not possible to connect to peer - " + err.Error())
		}
		go OnNewConnection(conn, model)
	}
}

// should be executed once we have some peers ready
func RunServer(model *Model) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	ipAddress := GetOutboundIP()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	ipAndPort := ipAddress + ":" + port

	PrintStatus("Listening for connections on IP:port " + ipAndPort)
	// register my own address in the address book
	model.RegisterMyAddress(ipAndPort)
	// connect to neighborhood
	ConnectToNeighborhood(model)
	// broadcast presence
	go model.BroadCastJson(MakePresent(ipAndPort))

	// listen
	for {
		conn, err := ln.Accept()
		if err != nil {
			PrintStatus("Connection failed -> " + err.Error())
			continue
		}
		go OnNewConnection(conn, model)
	}
}

// pull data from the remote peer
func InitialConnection(conn net.Conn, model *Model) {
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
	// register the peers in the application
	model.AddNetworkPeers(peers.Data.Peers)
	// print peers
	model.PrintPeers()
	model.sequencerPublicKey = &peers.Data.SequencerPk
	// drop this initial connection
	_ = conn.Close()
}

func StartSequencer(model *Model) {
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	model.amISequencer = true
	model.pk = privateKey
	model.sequencerPublicKey = &privateKey.PublicKey

	go StartSequencerWork(model)
}

func StartSequencerWork(model *Model) {
	for {
		time.Sleep(10 * time.Second)
		PrintStatus("Creating new block " + strconv.Itoa(model.currentBlockNumber))
		transactionsToSend := make([]string, 0, 0)
		for transactionId, transaction := range model.transactionsWaiting {
			// check if it was already processed in different block
			if transaction == nil {
				continue
			}
			// add to the current block
			transactionsToSend = append(transactionsToSend, transactionId)
		}
		PrintStatus("Block created with " + strconv.Itoa(len(transactionsToSend)) + " transactions.")

		// check if there are any transactions to send
		if len(transactionsToSend) == 0 {
			continue
		}
		// sort the transactions
		sort.Strings(transactionsToSend)
		// process the transactions
		for _, transactionId := range transactionsToSend {
			transaction, _ := model.transactionsWaiting[transactionId]
			model.ledger.DoSignedTransaction(transaction)
			model.transactionsWaiting[transactionId] = nil
			model.transactionsProcessed[transactionId] = true
		}

		// create block
		block := SequencerBlock{BlockNumber: model.currentBlockNumber, TransactionIds: transactionsToSend}
		// increase the block number
		model.currentBlockNumber++
		// sign block
		signedBlock := block.SignBlock(model.pk)
		// broadcast the block
		model.BroadCastJson(MakeSignedSequencerBlockDto(signedBlock))
	}
}

func Startup() {
	model := MakeModel()

	PrintStatus("Enter IP address and the port of the peer.")
	stdinReader := bufio.NewReader(os.Stdin)
	ipPort, _ := stdinReader.ReadString('\n')

	conn, err := net.Dial("tcp", strings.TrimSpace(ipPort))
	// if the peer does not exist, don't connect and start the network
	if err != nil {
		PrintStatus("It was not possible to connect - creating own network & starting the sequencer mode.")
		StartSequencer(&model)
	} else {
		// perform initial sync with the remote peer
		InitialConnection(conn, &model)
	}

	RunServer(&model)
}

func main() {
	Startup()
}
