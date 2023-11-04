package main

import (
	"fmt"
	"strconv"
)

func (cli *CLI) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("add block success!")
}

func (cli *CLI) printChain() {
	iter := cli.blockchain.Iterator()
	for {
		next := iter.Next()
		fmt.Printf("Prev hash:%x\n", next.PreBlockHash)
		fmt.Printf("Data: %s\n", next.Data)
		fmt.Printf("Hash: %x\n", next.Hash)
		pow := NewProofOfWork(next)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(next.PreBlockHash) == 0 {
			break
		}
	}
}
