package main

import (
	"crypto/sha256"
	"strconv"
)

func (d *Draw) ComputeHash() []byte {
	hash := sha256.New()
	WriteStringToHashSafe(&hash, "lottery"+strconv.Itoa(d.Slot))
	return hash.Sum(nil)
}

func GenerateDrawHash(slot int) []byte {
	hash := sha256.New()
	WriteStringToHashSafe(&hash, "lottery"+strconv.Itoa(slot))
	return hash.Sum(nil)
}
