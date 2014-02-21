package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func GenerateSecretAesKey() (key []byte, err error) {
	keySize := 32
	key = make([]byte, keySize)

	io.ReadFull(rand.Reader, key)
	return
}

func main() {
	aesSecret, _ := GenerateSecretAesKey()
	fmt.Println("binary: ", aesSecret)
	hexStr := fmt.Sprintf("%x", aesSecret)
	fmt.Println("hex: " + hexStr)

	aesSecret2, _ := hex.DecodeString(hexStr)
	fmt.Println("binary: ", aesSecret2)
}
