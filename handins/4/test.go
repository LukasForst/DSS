package main

import (
	"fmt"
	"math/rand"
)

func RandomBytes(size int) []byte {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(rand.Intn(255))
	}
	return data
}

func CheckSigning(keySize int, testRounds int, bytesSize int) {
	key := KeyGen(keySize)
	for i := 0; i < testRounds; i++ {
		data := RandomBytes(bytesSize)
		signature := key.GetSignature(data)
		// check signature
		if !key.CheckSignature(data, signature) {
			panic("Same data has different signature!")
		}
		// change random bit
		idx := rand.Intn(len(data))
		data[idx] = data[idx] + 1
		if key.CheckSignature(data, signature) {
			panic("Different data has same signature!")
		}
		// print progress
		if i%(testRounds/100) == 0 {
			fmt.Printf("Round %d/100 OK\n", 100*i/testRounds)
		}
	}

}

func main() {
	CheckSigning(2048, 1000, 255)
}
