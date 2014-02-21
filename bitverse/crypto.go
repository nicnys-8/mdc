package bitverse

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io"
)

/// PUBLIC

func GenerateAesSecret() (secret string, err error) {
	keySize := 32
	key := make([]byte, keySize)

	io.ReadFull(rand.Reader, key)
	secret = encodeHex(key)
	return
}

func UniqueHashkey() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	// calculate sha-1 hash
	hasher := sha1.New()
	hasher.Write([]byte(u.String()))

	return encodeHex(hasher.Sum(nil))
}

func HashkeyFromString(str string) string {
	// calculate sha-1 hash
	hasher := sha1.New()
	hasher.Write([]byte(str))

	return encodeHex(hasher.Sum(nil))
}

/// PRIVATE

func hex2Bin(hexStr string) (bytes []byte, err error) {
	bytes, err = hex.DecodeString(hexStr)
	return
}

func encodeHex(bytes []byte) string {
	return fmt.Sprintf("%x", bytes)
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		err := errors.New("failed to base64 decode payload")
		return nil, err
	}
	return data, nil
}

// aes stuff

func encryptAES(hexKey string, text string) string {
	key, err := hex2Bin(hexKey)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	b := encodeBase64([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return encodeBase64(ciphertext)
}

func decryptAES(hexKey string, ciphertext string) (string, error) {
	key, err := hex2Bin(hexKey)
	if err != nil {
		panic(err)
	}

	text, err := decodeBase64(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	if len(text) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)

	temp, err := decodeBase64(string(text))
	if err != nil {
		return "", err
	}

	return string(temp), nil
}

// rsa stuff

//key, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
