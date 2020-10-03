package main

import (
	"fmt"
	"math/rand"
	"time"
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
		signature := key.SignatureForData(data)
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

func PrintOperationSpeed(data []byte, elapsed time.Duration) {
	// data size to bits divided by elapsed milliseconds transformed to seconds
	bps := ((int64(len(data)) * 8) * 1e9) / elapsed.Nanoseconds()
	fmt.Printf("Operation speed: %d bps\n", bps)
}

func A1() {
	// #1
	CheckSigning(2048, 1000, 255)
}

func A2() {
	// #2
	// generate 10MB of the data
	megaBytes := 50
	data := RandomBytes(megaBytes * 1024 * 1024)
	start := time.Now()
	Sha256(data)
	elapsed := time.Since(start)
	fmt.Printf("%dMB of data hashed in: %s\n", megaBytes, elapsed)
	// 50MB of data hashed in: 118.986757ms

	PrintOperationSpeed(data, elapsed)
	// Operation speed: 3525017494 bps
}

func A3() {
	data := RandomBytes(255)
	// will always be 32 bytes
	hash := Sha256(data)
	key := KeyGen(2000)

	start := time.Now()
	key.SignatureForHash(hash)
	elapsed := time.Since(start)

	fmt.Printf("Data signed with 2000 bits key in: %s\n", elapsed)
	// Data signed with 2000 bits key in: 6.038µs

	PrintOperationSpeed(hash, elapsed)
	// Operation speed: 42398145 bps
}

func A4() {
	// hash is 32 bytes long == 256 bits
	// with 2000 bits key, we were able to hash it in cca 6 microseconds
	// hashing speed is 3525017494 bps, signing is 42398145 bps

	// it is much more efficient (~ 100x) to hash the message and then
	// sign the hash

	// numbers:
	// 3525017494
	// 0042398145
}

func main() {
	A1()
	fmt.Println("---------------------")
	A2()
	fmt.Println("---------------------")
	A3()
	fmt.Println("---------------------")
	A4()
}
