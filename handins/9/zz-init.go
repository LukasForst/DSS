package main

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"sync"
)

type InitialAccounts struct {
	PrivateKeys  []rsa.PrivateKey
	InitialState []int
}

func InitModel(
	peerKey *rsa.PrivateKey,
	accounts *InitialAccounts,
) PeerModel {
	cache := OperationCache{
		hasBroadcastTransaction: make(map[string]bool),
		hasBroadcastBlock:       make(map[string]bool),
		hasBroadcastPeerJoined:  make(map[string]bool),
		hasExecutedTransaction:  make(map[string]bool),
	}

	network := Network{
		connections: make(map[string]net.Conn),
		cMutex:      sync.RWMutex{},
		myIpPort:    "",
		peersList:   make(map[string]bool),
		pMutex:      sync.RWMutex{},
	}

	seed := 1234
	hardness := big.NewInt(100000000)

	initialAccountsState := make(map[string]int)
	for i, key := range accounts.PrivateKeys {
		initialAccountsState[FromRsaPubToAccount(&key.PublicKey)] = accounts.InitialState[i]
	}

	genesis := GenesisBlock{
		Hash:                 "",
		Seed:                 seed,
		InitialAccountStates: initialAccountsState,
		NextBlocksHashes:     make([]string, 0, 0),
	}
	genesis.Hash = genesis.ComputeHash()

	blockChain := BlockChain{
		Blocks:            make(map[string]Block),
		GenesisBlock:      genesis,
		Hardness:          hardness,
		Seed:              seed,
		SlotLengthSeconds: 10,
		lock:              sync.Mutex{},
	}

	ledger := PeerLedger{
		Accounts: make(map[string]int),
		lock:     sync.Mutex{},
	}

	for account, amount := range initialAccountsState {
		ledger.Accounts[account] = amount
	}

	return PeerModel{
		network:             &network,
		blockChain:          &blockChain,
		ledger:              &ledger,
		ledgerBlockId:       genesis.Hash,
		cache:               &cache,
		waitingTransactions: make(map[string]*SignedTransaction),
		peerKey:             peerKey,
	}
}

func LoadFromFile(input string) InitialAccounts {
	dat, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}
	var initialAccounts InitialAccounts
	err = json.Unmarshal(dat, &initialAccounts)
	if err != nil {
		log.Fatal(err)
	}
	return initialAccounts
}

func GenerateAndDump(accountsCount int, output string) {
	init := InitialAccounts{
		PrivateKeys:  make([]rsa.PrivateKey, 0, accountsCount),
		InitialState: make([]int, 0, accountsCount),
	}

	for i := 0; i < accountsCount; i++ {
		privateKey, _ := rsa.GenerateKey(cryptoRand.Reader, 2048)
		init.PrivateKeys = append(init.PrivateKeys, *privateKey)
		init.InitialState = append(init.InitialState, 10000000)
	}
	jsons, err := json.Marshal(init)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(output, jsons, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
