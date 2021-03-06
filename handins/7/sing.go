package main

import (
	"crypto/sha256"
	"math/big"
)

func Sha256(data []byte) []byte {
	sha := sha256.Sum256(data)
	return sha[:]
}

// convert sha to bigint
func BytesToInt(data []byte) *big.Int {
	i := new(big.Int)
	i.SetBytes(data)
	return i
}

// compute sha
func Sha256AsInt(data []byte) *big.Int {
	return BytesToInt(Sha256(data))
}

// generates fixed size signature
// size of the signature is size in bytes if the Key.n
func (k *SecretKey) SignatureForData(data []byte) []byte {
	return k.SignatureForHash(Sha256(data))
}

// generates fixed size signature
// size of the signature is size in bytes if the Key.n
func (k *SecretKey) SignatureForHash(hash []byte) []byte {
	sha := BytesToInt(hash)
	// as we're using RSA for signing messages
	// we must use secret key to produce signature
	encrypted := k.Decrypt(sha)
	encryptedBytes := encrypted.Bytes()
	// determine maximal size for the signature
	// the biggest output from RSA can be N-1, thus taking size of N
	signatureSize := len(k.N.Bytes())
	// pad RSA output with zeros to achieve correct size
	if len(encryptedBytes) < signatureSize {
		padding := make([]byte, signatureSize-len(encryptedBytes))
		encryptedBytes = append(padding, encryptedBytes...)
	}

	return encryptedBytes
}

// verifies signature
func (k *PublicKey) CheckSignature(data []byte, signature []byte) bool {
	// check signature size
	if len(k.N.Bytes()) != len(signature) {
		return false
	}
	// converts signature to the big int
	s := BytesToInt(signature)
	// using RSA for the signing,
	// we must encrypt the the signature
	// decrypt the signature
	ds := k.Encrypt(s)
	// get sha of the data
	sha := Sha256AsInt(data)
	// compare hashes
	return sha.Cmp(ds) == 0
}
