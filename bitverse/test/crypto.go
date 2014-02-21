package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// https://leanpub.com/gocrypto/read#leanpub-auto-chapter-two-symmetric-ciphers

const RSAKeySize = 3072

var defaultLabel = []byte{}

func GenerateSecretAesKey() (key []byte, err error) {
	keySize := 32
	key = make([]byte, keySize)

	io.ReadFull(rand.Reader, key)
	return
}

func maxMessageLength(key *rsa.PublicKey) int {
	if key == nil {
		return 0
	}
	return (key.N.BitLen() / 8) - (2 * sha256.Size) - 2
}

/// PRIVATE

func decryptRsa(prv *rsa.PrivateKey, ct []byte) (pt []byte, err error) {
	hash := sha256.New()
	pt, err = rsa.DecryptOAEP(hash, rand.Reader, prv, ct, defaultLabel)
	return
}

func encryptRsa(pub *rsa.PublicKey, pt []byte) (ct []byte, err error) {
	if len(ct) > maxMessageLength(pub) {
		err = fmt.Errorf("message is too long")
		return
	}

	hash := sha256.New()
	ct, err = rsa.EncryptOAEP(hash, rand.Reader, pub, pt, defaultLabel)
	return
}

func generatePrivatePem(prv *rsa.PrivateKey) (prvPem string, err error) {
	cert := x509.MarshalPKCS1PrivateKey(prv)
	blk := new(pem.Block)
	blk.Type = "RSA PRIVATE KEY"
	blk.Bytes = cert

	var b bytes.Buffer
	err = pem.Encode(&b, blk)
	if err != nil {
		return
	}

	prvPem = b.String()
	return
}

func generatePublicPem(pub *rsa.PublicKey) (pubPem string, err error) {
	cert, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return
	}

	blk := new(pem.Block)
	blk.Type = "RSA PUBLIC KEY"
	blk.Bytes = cert

	var b bytes.Buffer
	err = pem.Encode(&b, blk)
	if err != nil {
		return
	}

	pubPem = b.String()
	return
}

func generatePemKeys() (prvPem string, pubPem string, err error) {
	key, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		return
	}

	prvPem, err = generatePrivatePem(key)
	if err != nil {
		return
	}

	pubPem, err = generatePublicPem(&key.PublicKey)
	if err != nil {
		return
	}

	return
}

func exportPem(filename string, prvPem string, pubPem string) (err error) {
	privateKeyFile, err := os.Create(filename)
	if err != nil {
		return
	}

	privateKeyFile.WriteString(prvPem)
	privateKeyFile.Sync()

	publicKeyFile, err := os.Create(filename + ".pub")
	if err != nil {
		return
	}

	publicKeyFile.WriteString(pubPem)
	publicKeyFile.Sync()
	return
}

func importKeyFromPem(filename string) (prv *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	cert, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	for {
		var blk *pem.Block
		blk, cert = pem.Decode(cert)
		if blk == nil {
			break
		}
		switch blk.Type {
		case "RSA PRIVATE KEY":
			prv, err = x509.ParsePKCS1PrivateKey(blk.Bytes)
			return
		case "RSA PUBLIC KEY":
			var in interface{}
			in, err = x509.ParsePKIXPublicKey(blk.Bytes)
			if err != nil {
				return
			}
			pub = in.(*rsa.PublicKey)
			return
		}
		if cert == nil || len(cert) == 0 {
			break
		}
	}
	return
}

func Sign(prv *rsa.PrivateKey, m []byte) (sig []byte, err error) {
	h := sha256.New()
	h.Write(m)
	d := h.Sum(nil)
	sig, err = rsa.SignPSS(rand.Reader, prv, crypto.SHA256, d, nil)
	return
}

func Verify(pub *rsa.PublicKey, m, sig []byte) (err error) {
	h := sha256.New()
	h.Write(m)
	d := h.Sum(nil)
	return rsa.VerifyPSS(pub, crypto.SHA256, d, sig, nil)
}

/// PUBLIC

func GeneratePem(filename string) (err error) {
	prvPem, pubPem, err := generatePemKeys()
	if err != nil {
		return
	}
	err = exportPem(filename, prvPem, pubPem)
	return
}

func ImportPem(filename string) (prv *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	prv, _, err = importKeyFromPem(filename)
	if err != nil {
		return
	}

	_, pub, err = importKeyFromPem(filename + ".pub")
	if err != nil {
		return
	}

	return
}

func main() {
	aesSecret, _ := GenerateSecretAesKey()
	fmt.Println("binary: ", aesSecret)
	hexStr := fmt.Sprintf("%x", aesSecret)
	fmt.Println("hex: " + hexStr)

	aesSecret2, _ := hex.DecodeString(hexStr)
	fmt.Println("binary: ", aesSecret2)

	GeneratePem("cert")
	prv, pub, err := ImportPem("cert")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(prv)
	fmt.Println(pub)
}
