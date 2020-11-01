package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func HashKey(key *SecretKey) []byte {
	// compute hash of the private key of 32 bits
	hash := sha256.New()
	hash.Write(key.D.Bytes())
	hash.Write(key.N.Bytes())
	skHash := hash.Sum(nil)
	return skHash
}

func Generate(filename string, password string) string {
	key := KeyGen(2048)
	publicKey, err := encodeToBase64(key.PK)
	if err != nil {
		panic(err)
	}

	aesKey, salt, err := DeriveKey([]byte(password), nil)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// hash plaintext secret key
	skHash := HashKey(&key.SK)
	// write 32 bytes hash of the key
	_, err = file.Write(skHash)
	if err != nil {
		panic(err)
	}
	// write 32 bytes of salt
	_, err = file.Write(salt)
	if err != nil {
		panic(err)
	}
	// create array for the key
	materialToEncrypt := make([]byte, 256*2)
	copy(materialToEncrypt[0:256], key.SK.D.Bytes())
	copy(materialToEncrypt[256:512], key.SK.N.Bytes())
	// encrypt the file
	encryptedData, err := Encrypt(aesKey, materialToEncrypt)
	if err != nil {
		panic(err)
	}
	// write encrypted data to the file
	_, err = file.Write(encryptedData)
	if err != nil {
		panic(err)
	}

	return publicKey
}

//
type Signature struct {
	Signature []byte
}

func Sign(filename string, password string, msg []byte) Signature {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read hash-32-byte, salt-32-byte, ciphertext-D-256-bytes, ciphertext-N-256-bytes
	readBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	skHash := readBytes[0:32]
	salt := readBytes[32:64]
	ciphertext := readBytes[64:]

	aesKey, _, err := DeriveKey([]byte(password), salt)
	if err != nil {
		panic(err)
	}
	plaintext, err := Decrypt(aesKey, ciphertext)
	if err != nil {
		time.Sleep(5 * time.Second)
		panic(err)
	}

	D := BytesToInt(plaintext[0:256])
	N := BytesToInt(plaintext[256:512])

	skDecrypted := SecretKey{D: D, N: N}

	decryptedHash := HashKey(&skDecrypted)
	if !bytes.Equal(skHash, decryptedHash) {
		panic("hashes don't match!")
	}

	return Signature{skDecrypted.SignatureForData(msg)}
}

func main() {
	password := "prdel"
	filename := "secret.enc"
	pk := Generate(filename, password)
	fmt.Printf("Public Key: %s\n", pk)
	_ = Sign(filename, password, []byte("some message"))
	fmt.Println("OK")
}
