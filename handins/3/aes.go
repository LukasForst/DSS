package main

import (
	// 	"crypto/aes"
	"crypto/rand"
)

type File struct {
	filename string
	content  string
}

//Write ciphertext to given file
func EncryptToFile(key []byte, filename File) File {

	return filename
}

//Decrypt ciphertext from file and output plaintext
func DecryptFromFile() {

}

func AesKeyGen() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		//panic(fmt.Sprintf())
	}
}
