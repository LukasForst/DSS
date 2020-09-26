package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

func EncryptToFile(key []byte, inputFilename string, encFilename string, iv []byte) *os.File {

	//var fileContent string //read filecontent from file

	inputFile, err := os.Open(inputFilename)
	data, err := ioutil.ReadAll(inputFile)

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
		panic(fmt.Sprintf("Creation of encrypted output file was not successful."))
	}

	//fileContent = hex.EncodeToString(ciphertext)

	return outputFile
}

//Decrypt cyphertext in given encrypted file and write decrypted plaintext to decrypted file

func DecryptFromFile(key []byte, encFilename string, decFilename string, iv []byte) *os.File {

	//var fileContent string //read filecontent from file

	inputFile, err := os.Open(encFilename)

	if err != nil {
		panic(fmt.Sprintf("Encrypted file could not be opened: ", encFilename))
	}

	defer inputFile.Close()

	//ciphertext, _ := hex.DecodeString(fileContent)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(fmt.Sprintf("Cipher could not be created"))
	}

	//ciphertext = ciphertext[aes.BlockSize:]
	// if len(ciphertext)%aes.BlockSize != 0 {
	// 	panic("Ciphertext length is not a multiple of the cipher block size.")
	// }

	stream := cipher.NewCTR(block, iv)
	//stream.XORKeyStream(ciphertext, ciphertext)
	//fileContent = string(ciphertext[:])
	outputFile, err := os.OpenFile(decFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(fmt.Sprintf("Decoded file could not be created: ", decFilename))
	}
	defer outputFile.Close()

	reader := &cipher.StreamReader{S: stream, R: inputFile}
	if _, err := io.Copy(outputFile, reader); err != nil {
		panic(fmt.Sprintf("Creation of decrypted output file was not successful."))
	}

	return outputFile
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

	key := AesKeyGen()

	inputFilename := ""
	encFilename := ""
	decFilename := ""

	iv := IVGen()

	// if len(plaintext)%aes.BlockSize != 0 {
	// 	//panic or do padding?
	// }

	//fmt.Println("Original Plaintext:", plaintext)

	EncryptToFile(key, inputFilename, encFilename, iv)
	//fmt.Println("Encrypted ciphertext: ", ciphertext)

	DecryptFromFile(key, encFilename, decFilename, iv)
	//fmt.Println("Decrypted Plaintext: ", decryptedText)

}
