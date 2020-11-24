package main

import (
	"crypto/sha256"
	"strconv"
)

func (d *Draw) ComputeHash(seed int) []byte {
	hash := sha256.New()
	WriteStringToHashSafe(&hash, "lottery")
	WriteStringToHashSafe(&hash, strconv.Itoa(d.Slot))
	WriteStringToHashSafe(&hash, strconv.Itoa(seed))
	return hash.Sum(nil)
}

func GenerateDrawHash(slot int, seed int) []byte {
	hash := sha256.New()
	WriteStringToHashSafe(&hash, "lottery")
	WriteStringToHashSafe(&hash, strconv.Itoa(slot))
	WriteStringToHashSafe(&hash, strconv.Itoa(seed))
	return hash.Sum(nil)
}
