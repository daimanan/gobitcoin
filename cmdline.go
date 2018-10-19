//命令行支持
package main

import (
	"os"
	"fmt"
	"flag"
	"log"
)

const usage = `
	createChain --address ADDRESS "创建一个区块"
	addBlock --data DATA "添加一个区块到区块链中"
	send --from FROM --to TO --amount AMOUNT "由 FROM 给 TO 转款 AMOUNT"
	printChain           "打印所有区块信息"
`
const CreateBlockCmdString = "createChain"
const AddBlockCmdString = "addBlock"
const PrintChainCmdString = "printChain"

//命令行
type CLI struct {
	//bc *BlockChain
}

//Usage信息打印
func (cli *CLI) printUsage() {
	fmt.Println("输入法的参数错误")
	fmt.Println(usage)
	os.Exit(1)
}

//参数检查
func (cli *CLI) parameterCheck() {
	if len(os.Args) < 2 {
		cli.printUsage()
	}
}

//命令行运行
func (cli *CLI) Run() {
	cli.parameterCheck()
	createBlockCmd := flag.NewFlagSet(CreateBlockCmdString, flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet(AddBlockCmdString, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PrintChainCmdString, flag.ExitOnError)

	//解析命令行参数 name:参数命令 value:默认值 usage:命令说明
	createBlockPara := createBlockCmd.String("address", "", "address info!")
	addBlockPara := addBlockCmd.String("data", "", "block transaction info!")

	switch os.Args[1] {
	case CreateBlockCmdString:
		//创建区块数据库
		//1.解析命令行
		err := createBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic("解析CreateBlockCmdString出错\n", err)
		}

		//2.解析命行成功
		if createBlockCmd.Parsed() {
			//命令参数不为空
			if *createBlockPara == "" {
				fmt.Println("创建区块数据不能为空")
				cli.printUsage()
			}
			cli.CreateChain(*createBlockPara)
		}
	case AddBlockCmdString:
		//添加区块的操作
		//1.解析命令行
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic("解析AddBlockCmdString出错\n", err)
		}

		//2.解析命行成功
		if addBlockCmd.Parsed() {
			//命令参数不为空
			if *addBlockPara == "" {
				fmt.Println("添加区块数据不能为空")
				cli.printUsage()
			}
			cli.AddBlock(*addBlockPara)
		}

	case PrintChainCmdString:
		//打印输出区块信息操作
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic("解析PrintChainCmdString出错\n", err)
		}

		if printChainCmd.Parsed() {
			cli.PrintChain()
		}

	default:
		cli.printUsage()
	}
}
