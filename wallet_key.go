package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"

	"fmt"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "wallet_%s.dat"
const addressChecksumLen = 4

// Wallet stores private and public keys
// type Wallet struct {
// 	PrivateKey ecdsa.PrivateKey
// 	PublicKey  []byte
// }

// NewWallet creates and returns a Wallet
// func NewWallet() *Wallet {
// 	private, public := newKeyPair()
// 	wallet := Wallet{private, public}

// 	return &wallet
// }

func newKey() ecdsa.PrivateKey {
	curve := elliptic.P256()
	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	return *priKey
}

func pubKeyOf(priKey ecdsa.PrivateKey) []byte {
	pubKey := append(priKey.PublicKey.X.Bytes(), priKey.PublicKey.Y.Bytes()...)
	return pubKey
}

// GetAddress returns wallet address
func addressOf(pubKey []byte) []byte {
	pubKeyHash := HashPubKey(pubKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// func newKeyPair() (ecdsa.PrivateKey, []byte) {
// 	curve := elliptic.P256()
// 	private, err := ecdsa.GenerateKey(curve, rand.Reader)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
// 	fmt.Printf("publkey = %x\n", pubKey)
// 	fmt.Printf("prikey  = %x\n", private.D.Bytes())
// 	return *private, pubKey
// }

func encodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return string(pemEncoded)
}

func decodePrivateKey(pemEncoded string) ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return *privateKey
}

func encodeKey(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func decodeKey(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}

func testEncodeAndDecodeKey() {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey

	fmt.Println(privateKey)
	fmt.Println(publicKey)

	encPriv, encPub := encodeKey(privateKey, publicKey)

	fmt.Println(encPriv)
	fmt.Println(encPub)

	decPriv, decPub := decodeKey(encPriv, encPub)

	fmt.Println(decPriv)
	fmt.Println(decPub)
}
