package main

import (
	"crypto/sha256"
	"math/big"
)

func Sha256(data []byte) []byte {
	sha := sha256.Sum256(data)
	return sha[:]
}

// compute sha
func Sha256AsInt(data []byte) *big.Int {
	return BytesToInt(Sha256(data))
}

// convert sha to bigint
func BytesToInt(data []byte) *big.Int {
	i := new(big.Int)
	i.SetBytes(data)
	return i
}

// generates fixed size signature
// size of the signature is size in bytes if the Key.n
func (k *Key) SignatureForData(data []byte) []byte {
	return k.SignatureForHash(Sha256(data))
}

// generates fixed size signature
// size of the signature is size in bytes if the Key.n
func (k *Key) SignatureForHash(hash []byte) []byte {
	sha := BytesToInt(hash)
	// encrypt sha
	encrypted := k.Encrypt(sha)
	encryptedBytes := encrypted.Bytes()
	// determine maximal size for the signature
	// the biggest output from RSA can be N-1, thus taking size of N
	signatureSize := len(k.n.Bytes())
	// pad RSA output with zeros to achieve correct size
	if len(encryptedBytes) < signatureSize {
		padding := make([]byte, signatureSize-len(encryptedBytes))
		encryptedBytes = append(padding, encryptedBytes...)
	}

	return encryptedBytes
}

// verifies signature
func (k *Key) CheckSignature(data []byte, signature []byte) bool {
	// check signature size
	if len(k.n.Bytes()) != len(signature) {
		return false
	}
	// converts signature to the big int
	s := BytesToInt(signature)
	// decrypt signature
	ds := k.Decrypt(s)
	// get sha of the data
	sha := Sha256AsInt(data)
	// compare hashes
	return sha.Cmp(ds) == 0
}
