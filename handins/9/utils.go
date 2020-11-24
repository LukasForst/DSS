package main

import "hash"

func WriteStringToHashSafe(h *hash.Hash, str string) {
	WriteBytesToHashSafe(h, []byte(str))
}

func WriteBytesToHashSafe(h *hash.Hash, bytes []byte) {
	if _, err := (*h).Write(bytes); err != nil {
		panic(err)
	}
}
