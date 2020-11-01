package main

import "fmt"

func HappyScenario() {
	password := "NicePassword124!"
	filename := "secret.enc"
	dataToSign := []byte("some message in bytes")
	// generate key
	encodedPk := Generate(filename, password)
	fmt.Printf("Public Key: %s\n", encodedPk)
	// sign key
	signature := Sign(filename, password, dataToSign)

	// test that it was indeed signed correctly
	// decode PK from the base64
	var pk PublicKey
	err := decodeFromBase64(&pk, encodedPk)
	if err != nil {
		panic(err)
	}

	// check the signature
	if !pk.CheckSignature(dataToSign, signature.Signature) {
		panic("Signature is incorrect!")
	}

	fmt.Println("OK")
}

func WrongPasswordScenario() {
	password := "NicePassword124!"
	filename := "secret.enc"
	dataToSign := []byte("some message in bytes")
	// generate key
	encodedPk := Generate(filename, password)
	fmt.Printf("Public Key: %s\n", encodedPk)
	// this should fail on panic: cipher: message authentication failed
	Sign(filename, password+"no!", dataToSign)
	panic("The code should not end up here! It means that Sign was successful.")
}

func WeakPasswordScenario() {
	password := "Password"
	filename := "secret.enc"
	// generate key
	Generate(filename, password)
	panic("The code should not end up here! It means that Generate was successful.")
}

func main() {
	HappyScenario()
	WeakPasswordScenario()
	WrongPasswordScenario()
}
