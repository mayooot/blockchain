package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const targetBits = 24
const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{
		block:  b,
		target: target,
	}
}

func (p *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join([][]byte{
		p.block.PreBlockHash,
		p.block.Data,
		Int64ToBytes(p.block.Timestamp),
		Int64ToBytes(int64(targetBits)),
		Int64ToBytes(int64(nonce)),
	}, []byte{})
}

func (p *ProofOfWork) Run() (nonce int, hash [sha256.Size]byte) {
	var hashInt big.Int
	for nonce < maxNonce {
		p.block.Nonce = nonce
		data := p.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(p.target) == -1 {
			fmt.Printf("mining success, nonce %d hash %x\n", nonce, hash)
			return
		}
		nonce++
	}
	return
}

func (p *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := p.prepareData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(p.target) == -1
}
