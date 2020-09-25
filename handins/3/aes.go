package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type File struct {
	filename string
	content  string
}

//Write ciphertext to given file
func EncryptToFile(key []byte, inputFilename string, encFilename string, iv []byte) *os.File {

	//var fileContent string //read filecontent from file

	inputFile, err := os.Open(inputFilename)

	if err != nil {
		panic(fmt.Sprintf("Input file could not be openend: ", inputFilename))
	}

	defer inputFile.Close()

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(fmt.Sprintf("Cipher could not be created"))
	}

	//ciphertext := make([]byte, aes.BlockSize+len(fileContent))

	stream := cipher.NewCTR(block, iv)
	//stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(fileContent))

	outputFile, err := os.OpenFile(encFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		panic(fmt.Sprintf("Encoded file could not be created: ", encFilename))
	}

	defer outputFile.Close()

	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	if _, err := io.Copy(writer, inputFile); err != nil {
		panic(fmt.Sprintf("Creation of output file and encryption was not successful."))
	}

	//fileContent = hex.EncodeToString(ciphertext)

	return outputFile
}

//Decrypt ciphertext from file and output plaintext
func DecryptFromFile(key []byte, string, iv []byte) string {

	//var fileContent string //read filecontent from file

	ciphertext, _ := hex.DecodeString(fileContent)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		//panic()
	}

	ciphertext = ciphertext[aes.BlockSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("Ciphertext length is not a multiple of the cipher block size.")
	}

	mode := cipher.NewCTR(block, iv)
	mode.XORKeyStream(ciphertext, ciphertext)
	fileContent = string(ciphertext[:])

	return fileContent
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
		//panic(fmt.Sprintf)
	}
	return iv

}

func main() {

	key := AesKeyGen()

	filename := ""  //determine filename from input?
	plaintext := "" //read plaintext from file

	iv := IVGen()

	if len(plaintext)%aes.BlockSize != 0 {
		//panic or do padding?
	}

	fmt.Println("Original Plaintext:", plaintext)

	ciphertext := EncryptToFile(key, filename, iv)
	fmt.Println("Encrypted ciphertext: ", ciphertext)

	decryptedText := DecryptFromFile(key, ciphertext, iv)
	fmt.Println("Decrypted Plaintext: ", decryptedText)

}
