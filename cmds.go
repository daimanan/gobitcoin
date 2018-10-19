package main

import "fmt"

type cliInterface interface {
	CreateChain(address string)
	AddBlock(data string)
	PrintChain()
}

//命令行添加区块
func (cli *CLI) AddBlock(data string) {
	bc := GetBlockChainHandler()
	bc.AddBlock(data)
}

//打印数据
func (cli *CLI) PrintChain() {
	bc := GetBlockChainHandler()
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Println(" ============== current block ============== ")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PervHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerkelRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Difficuty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("IsValid : %v\n", NewProofOfWork(block).IsValid())
		if len(block.PervHash) == 0 {
			fmt.Println(" ==============  print over  ============== ")
			break
		}
	}
}

//创建区块链
func (cli *CLI) CreateChain(address string) {
	bc := InitBlockChain()
	defer bc.db.Close()
	fmt.Println("Create blockchain successfully!")
}
