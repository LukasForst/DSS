package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"math/big"
)

func VerifyWon(
	draw *Draw,
	tickets *big.Int,
	hardness *big.Int,
	key *rsa.PublicKey,
) bool {
	drawHash := draw.ComputeHash()
	err := rsa.VerifyPSS(key, crypto.SHA256, drawHash, draw.Signature, nil)
	if err != nil {
		return false
	}

	return DidWin(draw, tickets, hardness)
}

// returns draw if won, nil otherwise
func RunLottery(
	slot int,
	tickets *big.Int,
	hardness *big.Int,
	key *rsa.PrivateKey,
) *Draw {
	draw := GenerateSignedDraw(slot, key)
	if DidWin(&draw, tickets, hardness) {
		return &draw
	}
	return nil
}

func DidWin(
	draw *Draw,
	tickets *big.Int,
	hardness *big.Int,
) bool {
	hDraw := ByteHash(draw.Signature)

	hDrawValue := big.NewInt(0).SetBytes(hDraw)

	res := big.NewInt(0)
	res.Mul(tickets, hDrawValue)

	return res.Cmp(hardness) >= 0
}

func GenerateSignedDraw(slot int, key *rsa.PrivateKey) Draw {
	drawHash := GenerateDrawHash(slot)
	// produces signature of 256 bytes
	signature, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, drawHash, nil)
	if err != nil {
		log.Fatal(err)
	}

	return Draw{Slot: slot, Signature: signature}
}
