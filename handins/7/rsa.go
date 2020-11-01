package main

import (
	"crypto/rand"
	"math/big"
)

const (
	// public exponent e should be 3 (the smallest possible value, which gives the fastest possible encryption)
	e = 3
)

type SecretKey struct {
	D *big.Int
	N *big.Int
}

type PublicKey struct {
	E *big.Int
	N *big.Int
}

type Key struct {
	SK SecretKey
	PK PublicKey
}

func (k *SecretKey) Decrypt(number *big.Int) *big.Int {
	// this can be actually really optimized using Chinese remainder theorem
	// and working in (p-1) and (q-1) - not necessary for the hand in though
	return big.NewInt(0).Exp(number, k.D, k.N)
}

func (k *PublicKey) Encrypt(number *big.Int) *big.Int {
	return big.NewInt(0).Exp(number, k.E, k.N)
}

//an integer k, such that the bit length of the generated modulus n = pq is k
func KeyGen(k int) Key {
	e := big.NewInt(e)
	p, pp := GeneratePrime(k / 2)
	q, qq := GeneratePrime(k / 2)
	n := big.NewInt(0).Mul(p, q)

	nn := big.NewInt(0).Mul(pp, qq)
	d := big.NewInt(0).ModInverse(e, nn)
	return Key{SK: SecretKey{D: d, N: n}, PK: PublicKey{E: e, N: n}}
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
