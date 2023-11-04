package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"
)

// proof of work：工作量证明机制，通过预先设定一个难度 targetBits，例如本段代码设置的 24，
// 当比特币网络节点对数据进行sha256加密得到一个哈希值，数据包括：
// 1. 区块头：区块的头部信息，包括版本号、前一区块的哈希值、Merkle 根、时间戳和难度目标
// 2. Nonce：挖矿者通过不断更改 Nonce的值，重新计算区块头的哈希，直到找到符合挖矿难度要求的哈希值
// 3. 时间戳：挖矿操作开始的时间
// 4. merkle 根：吧区块内所有的交易打包进行哈希
// 5. 难度目标：例如：targetBits 为 24，难度目标就是 1 << 2^(256-24)，判断哈希值是否有效，就看得到的哈希值是否前 24 位为0（也就是小于难度目标）

// 24 指的是计算出来的哈希值前 24 位必须是 0
const targetBits = 24
const maxNonce = math.MaxInt64

type Block struct {
	Timestamp    int64  // 矿工开始计算 header 哈希的时间点
	PreBlockHash []byte // 前一个区块的 header 哈希
	Hash         []byte // 根哈希，由当前区块中包含的所有交易的哈希值运算得出
	Data         []byte
	Nonce        int
}

type Blockchain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewGenesisBlock 创建一个创世块
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlockchain 创建一个有创世块的链条
func NewBlockchain() *Blockchain {
	return &Blockchain{blocks: []*Block{NewGenesisBlock()}}
}

// NewBlock 创建新的区块，要用 pow 找到有效的哈希，也就是挖到矿的人，才有资格创建新区块
func NewBlock(data string, preBlockHash []byte) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		PreBlockHash: preBlockHash,
		Hash:         []byte{},
		Data:         []byte(data),
		Nonce:        0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (c *Blockchain) AddBlock(data string) {
	preBlock := c.blocks[len(c.blocks)-1]
	nowBlock := NewBlock(data, preBlock.Hash)
	c.blocks = append(c.blocks, nowBlock)
}

// NewProofOfWork 创建一个工作量证明对象，并设置难度目标值
// block: 要进行工作量证明的区块
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// lsh: 左移，将 target 左移 256-targetBits 位，也就是 1*2^232
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{block: block, target: target}
	return pow
}

func (p *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		p.block.PreBlockHash,
		p.block.Data,
		IntToHex(p.block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

// Run 寻找有效哈希，也就是挖矿的过程
func (p *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	// sha256 生成 256 位的二进制哈希值，一个字节等于 8 比特，所以用 32 字节
	var hash [32]byte
	nonce := 0
	fmt.Printf("Mining the block containing \"%s\"\n", p.block.Data)
	for nonce < maxNonce {
		data := p.prepareData(nonce)
		hash = sha256.Sum256(data)
		// 将字节数组转为 big.Int
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 {
			// 如果 hashInt 小于 p.target
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")
	return nonce, hash[:]
}

// Validate 证明工作量，只要得到的哈希小于 target 就有效
func (p *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(p.target) == -1
}

// IntToHex 将 int64 转为字节数组
func IntToHex(num int64) []byte {
	buff := &bytes.Buffer{}
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func main() {
	blockchain := NewBlockchain()
	blockchain.AddBlock("foo send 1 BTC to bar")
	blockchain.AddBlock("bar send 0.15 BTC to bob")
	for _, block := range blockchain.blocks {
		fmt.Printf("Prev hash: %x\n", block.PreBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
