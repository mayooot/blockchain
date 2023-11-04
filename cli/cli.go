package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	blockchain *BlockChain
}

const usage = `
Usage:
	addblock -data BLOCK_DATA add a block to the blockchain
	printchain print all the blocks of the blockchain
`

func (cli *CLI) printUsage() {
	fmt.Print(usage)
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 将 data 的值存储在 addBlockData 中
	addBlockData := addBlockCmd.String("data", "", "Block Data")

	// os.Args[0] 是程序本身, [1] 是 addblock、printchain [2] 是 -data 后的数据
	switch os.Args[1] {
	case "addblock":
		// 解析从第三个参数开始的命令行参数，并将结果存储在 addBlockData 中
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
