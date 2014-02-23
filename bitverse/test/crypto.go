package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
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

func sign(prv *rsa.PrivateKey, msg string) (signature string, err error) {
	h := sha256.New()
	h.Write([]byte(msg))
	d := h.Sum(nil)
	sigBin, _ := rsa.SignPSS(rand.Reader, prv, crypto.SHA256, d, nil)
	signature = encodeBase64(sigBin)
	return
}

func verify(pub *rsa.PublicKey, msg string, signature string) (err error) {
	sig, _ := decodeBase64(signature)
	h := sha256.New()
	h.Write([]byte(msg))
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
	//aesSecret, _ := GenerateSecretAesKey()
	//fmt.Println("binary: ", aesSecret)
	//hexStr := fmt.Sprintf("%x", aesSecret)
	//fmt.Println("hex: " + hexStr)

	//aesSecret2, _ := hex.DecodeString(hexStr)
	//fmt.Println("binary: ", aesSecret2)

	GeneratePem("cert")
	prv, _, err := ImportPem("cert")
	if err != nil {
		fmt.Println(err)
	}

	GeneratePem("cert2")
	_, pub2, err := ImportPem("cert2")
	if err != nil {
		fmt.Println(err)
	}

	str := "hejsan hopp"
	signature, _ := sign(prv, str)
	fmt.Println(signature)

	err2 := verify(pub2, str, signature)

	if err2 != nil {
		fmt.Printf("failed to very signature")
	} else {
		fmt.Printf("signature ok")
	}

}
