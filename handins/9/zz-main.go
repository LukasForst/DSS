package main

import (
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"strconv"
)

func main() {
	privateKey, _ := rsa.GenerateKey(cryptoRand.Reader, 2048)
	for i := 0; i < 1000; i++ {
		draw := GenerateSignedDraw(i, privateKey)
		if len(draw.Signature) != 256 {
			log.Fatal("Its: " + strconv.Itoa(len(draw.Signature)))
		}
	}
	fmt.Println("hellow world")
}
