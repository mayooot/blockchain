package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	Timestamp    int64
	PreBlockHash []byte
	Hash         []byte
	Data         []byte
}

func NewBlock(data string, preBlockHash []byte) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Hash:         []byte{},
		Data:         []byte(data),
	}
	block.SetHash()
	return block
}

// SetHash Hash = sha256(PrevBlockHash + Data + Timestamp)
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PreBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// NewGenesisBlock 创建一个创世块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlockchain 创建一个有创世块的链条
func NewBlockchain() *Blockchain {
	return &Blockchain{blocks: []*Block{NewGenesisBlock()}}
}

func (c *Blockchain) AddBlock(data string) {
	preBlock := c.blocks[len(c.blocks)-1]
	nowBlock := NewBlock(data, preBlock.Hash)
	c.blocks = append(c.blocks, nowBlock)
}

func main() {
	blockchain := NewBlockchain()
	blockchain.AddBlock("Send 1 BTC to foo")
	blockchain.AddBlock("Send 1 BTC to bar")
	for _, block := range blockchain.blocks {
		fmt.Printf("Prev hash: %x\n", block.PreBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
