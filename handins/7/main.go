package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

// Generates RSA key, private/secret key is stored in the file
// encrypted by AES-256 GCM, under given password
// returns base64 encoded PublicKey
func Generate(filename string, password string) string {
	EnsurePasswordComplexity(password)
	// generate new RSA key
	key := KeyGen(2048)
	// encode PK in base64, it is then returned to the user
	publicKey, err := encodeToBase64(key.PK)
	if err != nil {
		panic(err)
	}
	// store SK in the file
	StoreSKInFile(key.SK, filename, password)
	return publicKey
}

// Signs given msg by the private/secret key from the keychain
func Sign(filename string, password string, msg []byte) Signature {
	sk := ReadSKFromFile(filename, password)
	// sign the data
	return Signature{sk.SignatureForData(msg)}
}

type Signature struct {
	Signature []byte
}

// minimum eight characters,
// at least one uppercase letter,
// one lowercase letter and one number
func EnsurePasswordComplexity(password string) {
	passwordLen := len(password) >= 8
	numberMatch, _ := regexp.MatchString("[0-9]", password)
	capitalCaseMatch, _ := regexp.MatchString("[A-Z]", password)
	lowerCaseMatch, _ := regexp.MatchString("[a-z]", password)

	if !(passwordLen && numberMatch && capitalCaseMatch && lowerCaseMatch) {
		panic("Weak password used! " +
			"Minimum eight characters, at least one uppercase letter, " +
			"one lowercase letter and one number.")
	}
}

func StoreSKInFile(key SecretKey, filename string, password string) {
	// derive aes key from the given password
	aesKey, salt, err := DeriveKey([]byte(password), nil)
	if err != nil {
		panic(err)
	}
	// create keychain file
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// write 32 bytes of salt
	_, err = file.Write(salt)
	if err != nil {
		panic(err)
	}
	// create array for the key
	materialToEncrypt := make([]byte, 256*2)
	// store firstly D, then N
	copy(materialToEncrypt[0:256], key.D.Bytes())
	copy(materialToEncrypt[256:512], key.N.Bytes())
	// encrypt the data
	encryptedData, err := Encrypt(aesKey, materialToEncrypt)
	if err != nil {
		panic(err)
	}
	// write encrypted data to the file
	_, err = file.Write(encryptedData)
	if err != nil {
		panic(err)
	}
}

func ReadSKFromFile(filename string, password string) SecretKey {
	// open the file
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// read salt-32-byte, ciphertext
	readBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	// prepare data
	salt := readBytes[0:32]
	ciphertext := readBytes[32:]
	// derive actual AES key from the password and the salt
	aesKey, _, err := DeriveKey([]byte(password), salt)
	if err != nil {
		panic(err)
	}
	// decrypt ciphertext
	plaintext, err := Decrypt(aesKey, ciphertext)
	// if the AES failed, wrong password was used
	if err != nil {
		// slowdown potential bruteforce attack
		// wait 5s, before returning, thus reducing
		// the bruteforce attack speed
		time.Sleep(5 * time.Second)
		panic(err)
	}
	// get data from the plaintext
	D := BytesToInt(plaintext[0:256])
	N := BytesToInt(plaintext[256:512])
	// build secret key
	return SecretKey{D: D, N: N}
}
