package main

import "fmt"

//程序执行入口
func main() {
	//创建区块链
	bc := NewBlockChain()
	bc.AddBlock("创建人获取100枚比特币")
	bc.AddBlock("创建人获取500枚比特币")


	//遍历区块链
	for index,block:=range bc.blocks{
		fmt.Println(" ============== current block index :", index)
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PervHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerkelRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Difficuty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)
	}

}
