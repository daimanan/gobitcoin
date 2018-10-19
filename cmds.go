package main

import (
	"fmt"
)

type cliInterface interface {
	CreateChain(address string)
	AddBlock(data string)
	PrintChain()
}

//命令行添加区块
func (cli *CLI) AddBlock(data string) {
	//bc := GetBlockChainHandler()
	//bc.AddBlock(data) //TODO
}

//打印数据
func (cli *CLI) PrintChain() {
	bc := GetBlockChainHandler()
	defer bc.db.Close()
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
		//fmt.Printf("Data : %s\n", block.Data) //TODO
		fmt.Printf("IsValid : %v\n", NewProofOfWork(block).IsValid())
		if len(block.PervHash) == 0 {
			fmt.Println(" ==============  print over  ============== ")
			break
		}
	}
}

//创建区块链
func (cli *CLI) CreateChain(address string) {
	bc := InitBlockChain(address)
	defer bc.db.Close()
	fmt.Println("Create blockchain successfully!")
}

//获取余额
func (cli *CLI) GetBalance(address string) {
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	utxos := bc.FindUTXOs(address)

	//总金额
	var total float64 = 0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	//return total
	fmt.Printf("用户%s的余额是%f", address, total)
}
