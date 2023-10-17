package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

// Wallet stores a collection of priKeys
// mapping address -> priKey
type Wallet struct {
	addr map[string]ecdsa.PrivateKey
}

// NewWallet creates Wallet and fills it from a file if it exists
func NewWallet(nodeID string) (*Wallet, error) {
	wallet := Wallet{}
	wallet.addr = make(map[string]ecdsa.PrivateKey)

	err := wallet.LoadFromFile(nodeID)

	return &wallet, err
}

// CreateKey adds a key to Wallet
func (w *Wallet) CreateKey() string {
	priKey := newKey()
	address := fmt.Sprintf("%s", addressOf(pubKeyOf(priKey)))

	w.addr[address] = priKey
	// savetofile?
	return address
}

// GetAddresses returns an array of addresses stored in the Wallet
func (w *Wallet) GetAddresses() []string {
	var addresses []string

	for address := range w.addr {
		addresses = append(addresses, address)
	}

	return addresses
}

// GetKey returns a priKey by its address
func (w *Wallet) GetKey(address string) ecdsa.PrivateKey {
	return w.addr[address]
}

// LoadFromFile loads Wallet from the file
func (w *Wallet) LoadFromFile(nodeID string) error {
	walletFile := fmt.Sprintf(walletFile, nodeID) //根据nodeID拿到文件名
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var encPriKeys []string
	// gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&encPriKeys)

	if err != nil {
		log.Panic(err)
	}
	// var wallet Wallet
	for _, encPriKey := range encPriKeys {
		priKey := decodePrivateKey(encPriKey)
		// wallet.addr[string(addressOf(pubKeyOf(priKey)))] = priKey
		w.addr[string(addressOf(pubKeyOf(priKey)))] = priKey
	}

	// w.addr = wallet.addr

	return nil
}

// SaveToFile saves wallets to a file
func (w Wallet) SaveToFile(nodeID string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(walletFile, nodeID)
	// gob.Register(elliptic.P256())
	var encPriKeys []string

	for _, priKey := range w.addr {
		encPriKey := encodePrivateKey(&priKey)
		encPriKeys = append(encPriKeys, encPriKey)
	}

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(encPriKeys)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
