package bitverse

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func encodePayload(payload string) (encodedPayload string) {
	encodedPayload = base64.StdEncoding.EncodeToString([]byte(payload))
	return
}

func decodePayload(encodedPayload string) (payload string) {
	payloadBuffer, err := base64.StdEncoding.DecodeString(encodedPayload)

	if err != nil {
		fatal("msg: failed to decode payload")
	}

	return string(payloadBuffer[0:len(payloadBuffer)])
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		err := errors.New("failed to decrypt payload")
		return nil, err
	}
	return data, nil
}

// Generate a symmetric key. This is suitable for session keys
// and other short-term key material.
func GenerateSecretKey() (key []byte, err error) {
	keySize := 32
	key = make([]byte, keySize)

	io.ReadFull(rand.Reader, key)
	return
}

func encrypt(key string, text string) string {
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

func decrypt(key string, ciphertext string) (string, error) {
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
