package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"
)

// 序列化时，使用了encoding/gob，切记必须要大些
type Block struct {
	Version       int64
	PrevBlockHash []byte
	MerkleRoot    []byte
	Timestamp     int64
	Nbits         int64
	Nonce         int64
	Transactions  []*Transaction
	Hash          []byte
	Height        int
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// create MerkleRoot though transactions
func (b *Block) createMerkelTreeRoot(Transactions []*Transaction) {
	var transactions [][]byte

	for _, tx := range Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	b.MerkleRoot = mTree.RootNode.Data
}

func TestMerkleTree() {

	//https://www.blockchain.com/btc/block/00000000000090ff2791fe41d80509af6ffbd6c5b10294e29cdf1b603acab92c

	data1, _ := hex.DecodeString("6b6a4236fb06fead0f1bd7fc4f4de123796eb51675fb55dc18c33fe12e33169d")
	data2, _ := hex.DecodeString("2af6b6f6bc6e613049637e32b1809dd767c72f912fef2b978992c6408483d77e")
	data3, _ := hex.DecodeString("6d76d15213c11fcbf4cc7e880f34c35dae43f8081ef30c6901f513ce41374583")
	data4, _ := hex.DecodeString("08c3b50053b010542dca85594af182f8fcf2f0d2bfe8a806e9494e4792222ad2")
	data5, _ := hex.DecodeString("612d035670b7b9dad50f987dfa000a5324ecb3e08745cfefa10a4cefc5544553")
	data6 := reverse2(data1)
	data7 := reverse2(data2)
	data8 := reverse2(data3)
	data9 := reverse2(data4)
	data10 := reverse2(data5)

	hehe := [][]byte{
		data6,
		data7,
		data8,
		data9,
		data10,
	}

	result := (*NewMerkleTree(hehe).RootNode).Data
	rev := reverse2(result)
	fmt.Printf("result=%x\n", rev)
}

// create MerkleRoot though transactions
func TestcreateMerkelTreeRoot(Transactions []*Transaction) []byte {
	var transactions [][]byte

	for _, tx := range Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
	//b.MerkleRoot= mTree.RootNode.Data
}

// 产生初始区块,传入了第一笔coinbase交易
func NewGenesisBlock(Transactions []*Transaction) *Block {
	block := &Block{int64(2), []byte{}, []byte("abc"), time.Now().Unix(), 111111, 100, Transactions, []byte{}, 0}

	var transactions [][]byte
	for _, tx := range Transactions {
		transactions = append(transactions, tx.Hash())
	}
	block.MerkleRoot = (*NewMerkleTree(transactions).RootNode).Data

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	//fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Printf("Prev. version: %s\n", strconv.FormatInt(block.Version, 10))
	fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
	fmt.Printf("merkleroot: %x\n", block.MerkleRoot)
	fmt.Printf("time: %s\n", strconv.FormatInt(block.Timestamp, 10))
	fmt.Printf("nbits: %s\n", strconv.FormatInt(block.Nbits, 10))
	fmt.Printf("nonce: %s\n", strconv.FormatInt(block.Nonce, 10))
	fmt.Printf("Hash: %x\n", block.Hash)
	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Printf("------------------------------------------------------------\n")
	fmt.Println()
	return block
}

func NewBlock(Transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{2, prevBlockHash, []byte("dfg"), time.Now().Unix(), 111111, 0, Transactions, []byte{}, height}

	var transactions [][]byte
	for _, tx := range Transactions {
		transactions = append(transactions, tx.Hash())
	}
	block.MerkleRoot = (*NewMerkleTree(transactions).RootNode).Data

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block

}
