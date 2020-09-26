package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type File struct {
	filename string
	content  string
}

//Encrypt plaintext in given input file and write encrypted ciphertext to encrypted file

func EncryptToFile(key []byte, inputFilename string, encFilename string) {

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		panic(fmt.Sprintf("Input file could not be openend: ", inputFilename))
	}
	defer inputFile.Close()
	plaintext, err := ioutil.ReadAll(inputFile)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("Cipher could not be created"))
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	outputFile, err := os.Create(encFilename)
	if err != nil {
		panic(fmt.Sprintf("Encoded file could not be created: ", encFilename))
	}
	defer outputFile.Close()
	_, err = outputFile.Write(ciphertext)
}

//Decrypt cyphertext in given encrypted file and write decrypted plaintext to decrypted file

func DecryptFromFile(key []byte, encFilename string, decFilename string) {
	inputFile, err := os.Open(encFilename)
	if err != nil {
		panic(fmt.Sprintf("Encrypted file could not be opened: ", encFilename))
	}
	defer inputFile.Close()

	ciphertext, err := ioutil.ReadAll(inputFile)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("Cipher could not be created"))
	}

	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	outputFile, err := os.Create(decFilename)
	if err != nil {
		panic(fmt.Sprintf("Encoded file could not be created: ", encFilename))
	}
	defer outputFile.Close()
	_, err = outputFile.Write(plaintext)
}

//generate 32 byte random key
func AesKeyGen() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		//panic(fmt.Sprintf())
	}
	return key
}

//geerate the IV nonce
func IVGen() []byte {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(fmt.Sprintf("IV could not be created."))
	}
	return iv

}

func main() {

	key, _ := hex.DecodeString("6368616e676520746869732070617373")

	inputFilename := "hello.txt"
	encFilename := "hello.enc.txt"
	decFilename := "hello.dec.txt"

	EncryptToFile(key, inputFilename, encFilename)
	//fmt.Println("Encrypted ciphertext: ", ciphertext)

	DecryptFromFile(key, encFilename, decFilename)
	//fmt.Println("Decrypted Plaintext: ", decryptedText)

}
