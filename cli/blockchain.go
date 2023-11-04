package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"

// BlockChain tail 尾部，存储的是最后一个块的 hash 值，存储最后的 tail 就能推到出整条链
type BlockChain struct {
	tail []byte
	db   *bolt.DB
}

func NewBlockChain() *BlockChain {
	var tail []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	_ = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		// 如果数据库中不存在链就创建一个，否则返回最后一个块的哈希
		if bucket != nil {
			tail = bucket.Get([]byte("1"))
		} else {
			// 创建 blockchain
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err)
			}

			err = bucket.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = bucket.Put([]byte("1"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tail = genesis.Hash
		}

		return nil
	})

	return &BlockChain{
		tail: tail,
		db:   db,
	}
}

func (c *BlockChain) AddBlock(data string) {
	var lastHash []byte
	var err error
	// 获取最后一个块的哈希值，用于生成新的哈希值
	_ = c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		lastHash = bucket.Get([]byte("1"))
		return nil
	})

	newBlock := NewBlock(data, lastHash)

	_ = c.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		err = bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		c.tail = newBlock.Hash
		return nil
	})
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (c *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: c.tail,
		db:          c.db,
	}
}

// Next 返回链中的下一个块
func (i *BlockchainIterator) Next() *Block {
	var block *Block
	_ = i.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		block = Deserialize(bucket.Get(i.currentHash))
		return nil
	})
	i.currentHash = block.PreBlockHash
	return block
}
