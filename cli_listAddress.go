package main

import (
	"fmt"
	"log"
)

func (cli *CLI) listAddresses(nodeID string) {
	wallet, err := NewWallet(nodeID)
	if err != nil {
		log.Panic(err)
	}
	addresses := wallet.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
