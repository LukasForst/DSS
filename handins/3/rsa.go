package main

import (
	"crypto/rand"
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

func main() {
	keyBits := 2048

	// generate key
	key := KeyGen(keyBits)

	// assert key length
	if key.n.BitLen() != keyBits {
		panic(fmt.Sprintf("Key bits are wrong! Expecting %d, was %d", keyBits, key.n.BitLen()))
	} else {
		fmt.Printf("Key has correct size of %d\n", keyBits)
	}

	// run the tests
	max := big.NewInt(100000)
	times := 1000
	fmt.Printf("Executing enc/dec tests, rounds: %d\n", times)
	for i := 0; i < times; i++ {
		num, _ := rand.Int(rand.Reader, max)
		numE := key.Encrypt(num)
		numD := key.Decrypt(numE)
		// assert num == numD
		if num.Cmp(numD) != 0 {
			panic(fmt.Sprintf("Decryption failed! Expecting %s, was %s", num.String(), numD.String()))
		}
	}
	fmt.Println("All tests passed!")
}
