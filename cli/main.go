package main

func main() {
	blockChain := NewBlockChain()
	defer blockChain.db.Close()
	cli := CLI{blockchain: blockChain}
	cli.Run()
}
