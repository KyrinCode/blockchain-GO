package main

import "fmt"

func (cli *CLI) createWallet(nodeID string) {
	wallet, _ := NewWallet(nodeID)
	address := wallet.CreateKey()
	wallet.SaveToFile(nodeID)

	fmt.Printf("Your new address: %s\n", address)
}
