package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"hash"
	"strconv"
)

func (t *SignedTransaction) ComputeHash() []byte {
	msgHash := sha256.New()

	WriteToHashSafe(&msgHash, t.ID)
	WriteToHashSafe(&msgHash, t.From)
	WriteToHashSafe(&msgHash, t.To)
	// intentionally hashing int as a string
	WriteToHashSafe(&msgHash, strconv.Itoa(t.Amount))
	msgHashSum := msgHash.Sum(nil)
	return msgHashSum
}

func (t *SignedTransaction) ComputeAndSetSignature(key *rsa.PrivateKey) {
	transactionHash := t.ComputeHash()
	signature, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, transactionHash, nil)
	if err != nil {
		panic(err)
	}

	t.Signature = base64.StdEncoding.EncodeToString(signature)
}

func (t *SignedTransaction) IsSignatureCorrect() bool {
	transactionHash := t.ComputeHash()
	// get public key of the from
	var fromPk rsa.PublicKey
	if err := json.Unmarshal([]byte(t.From), &fromPk); err != nil {
		panic(err)
	}

	payloadSignature, _ := base64.StdEncoding.DecodeString(t.Signature)
	err := rsa.VerifyPSS(&fromPk, crypto.SHA256, transactionHash, payloadSignature, nil)
	if err != nil {
		PrintStatus("Could not verify signature: " + err.Error())
	}
	return err == nil
}

func WriteToHashSafe(h *hash.Hash, str string) {
	if _, err := (*h).Write([]byte(str)); err != nil {
		panic(err)
	}
}
