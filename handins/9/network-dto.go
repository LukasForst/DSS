package main

type PresentDto struct {
	Type string
	// name of the server
	Data string
}

func MakePresent(address string) PresentDto {
	return PresentDto{Type: "present", Data: address}
}

type PeerList struct {
	Peers []string
}

type PeersListDto struct {
	Type string
	// names of the servers
	Data PeerList
}

func MakePeersList(data *PeerList) PeersListDto {
	return PeersListDto{Type: "peers-list", Data: *data}
}

type PeersRequestDto struct {
	Type string
}

func MakePeersRequest() PeersRequestDto {
	return PeersRequestDto{Type: "peers-request"}
}

type SignedBlockDto struct {
	Type string
	Data SignedBlock
}

func MakeSignedBlockDto(data *SignedBlock) SignedBlockDto {
	return SignedBlockDto{Type: "signed-block", Data: *data}
}

type AccountSetupDto struct {
	Type string
	Data AccountSetup
}

func MakeAccountSetupDto(setup *AccountSetup) AccountSetupDto {
	return AccountSetupDto{Type: "account-setup", Data: *setup}
}

type AccountSetup struct {
	AccountId string
	Amount    int
}

// data transfer object for transaction
type TransactionDto struct {
	Type string
	Data SignedTransaction
}

func MakeTransactionDto(transaction *SignedTransaction) TransactionDto {
	return TransactionDto{Type: "transaction", Data: *transaction}
}
