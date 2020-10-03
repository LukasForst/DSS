package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const (
	// public exponent e should be 3 (the smallest possible value, which gives the fastest possible encryption)
	e = 3
)

type Key struct {
	e *big.Int
	d *big.Int
	n *big.Int
}

func (k *Key) Decrypt(number *big.Int) *big.Int {
	// this can be actually really optimized using Chinese remainder theorem
	// and working in (p-1) and (q-1) - not necessary for the hand in though
	return big.NewInt(0).Exp(number, k.d, k.n)
}

func (k *Key) Encrypt(number *big.Int) *big.Int {
	return big.NewInt(0).Exp(number, k.e, k.n)
}

//an integer k, such that the bit length of the generated modulus n = pq is k
func KeyGen(k int) Key {
	e := big.NewInt(e)
	p, pp := GeneratePrime(k / 2)
	q, qq := GeneratePrime(k / 2)
	n := big.NewInt(0).Mul(p, q)

	nn := big.NewInt(0).Mul(pp, qq)
	d := big.NewInt(0).ModInverse(e, nn)
	return Key{e: e, d: d, n: n}
}

// Generates valid (gcd(e, p-1) == 1) prime of k bits
func GeneratePrime(k int) (*big.Int, *big.Int) {
	// convert e to bigint
	e := big.NewInt(e)
	// prepare variables
	one := big.NewInt(1)
	gcd := big.NewInt(0)
	p := big.NewInt(0)
	pp := big.NewInt(0)
	// while GCD(e, p-1) != 1
	for gcd.Cmp(one) != 0 {
		// generate probably prime
		p, _ = rand.Prime(rand.Reader, k)
		// get p-1
		pp = big.NewInt(0).Sub(p, one)
		// find gcd
		gcd.GCD(nil, nil, e, pp)
	}
	return p, pp
}

func Sha256(data []byte) []byte {
	sh := sha256.Sum256(data)
	return sh[:]
}

func BytesToInt(data []byte) *big.Int {
	i := new(big.Int)
	i.SetBytes(data)
	return i
}

func (k *Key) GetSignature(data []byte) []byte {
	sha := Sha256(data)
	i := BytesToInt(sha)
	encrypted := k.Encrypt(i)
	return encrypted.Bytes()
}

func (k *Key) CheckSignature(data []byte, signature []byte) bool {
	sha := Sha256(data)
	i := BytesToInt(sha)
	s := BytesToInt(signature)
	ds := k.Decrypt(s)
	return i.Cmp(ds) == 0
}

func main() {
	keyBits := 2048

	// generate key
	key := KeyGen(keyBits)
	var signature []byte
	var data []byte
	rounds := 1000
	for i := 0; i < rounds; i++ {
		randInt, _ := rand.Int(rand.Reader, big.NewInt(1000000000000000000))
		data = randInt.Bytes()

		signature = key.GetSignature(data)
		if !key.CheckSignature(data, signature) {
			panic("Same data has different signature!")
		}
		data[0] = data[0] + 1
		if key.CheckSignature(data, signature) {
			panic("Different data has same signature!")
		}
		if i%(rounds/100) == 0 {
			fmt.Printf("Round %d/100 OK\n", 100*i/rounds)
		}
	}
}
