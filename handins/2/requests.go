package main

type Present struct {
	Type string
	Data string
}

func MakePresent(address string) Present {
	return Present{Type: "present", Data: address}
}

type PeersList struct {
	Type string
	Data []string
}

func MakePeersList(data []string) PeersList {
	return PeersList{Type: "peers-list", Data: data}
}

type PeersRequest struct {
	Type string
}

func MakePeersRequest() PeersRequest {
	return PeersRequest{Type: "peers-request"}
}

type Transaction struct {
}
