//Ver3:使用bolt数据改写
package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"fmt"
)

//区块链数据库名
const dbFile = "blockchain.db"
const blockBucket = "bucket"
const lastHashKey = "LastHashKey"

//创世区块信息
const genesisInfo = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type blockChainInterface interface {
	AddBlock(txs []*Transaction)
	NewIterator() *BlockChainIterator
	FindUTXOTransactions(address string) []*Transaction
}

//区块链
type BlockChain struct {
	//Ver1 Ver2使用的版本
	//blocks []*Block //用区块切片模拟区块链

	//Ver3:使用bolt数据改写
	//数据库操作句柄
	db *bolt.DB
	//最后一个区声的哈希值
	tail []byte
}

//添加区块到区块链中
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//Ver1 Ver2版本使用
	//lastBlock := bc.blocks[len(bc.blocks)-1]
	//block := NewBlock(data, lastBlock.Hash)
	//bc.blocks = append(bc.blocks, block)

	//Ver3：使用bolt数据库实例化区块链

	var prevBlockHash []byte //前一区块的hash
	//读取本地数据库
	bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		//数据库没有找到
		if bucket == nil {
			log.Panic("数据打开失败，程序非正常退出")
			os.Exit(1)
		}
		prevBlockHash = bucket.Get([]byte(lastHashKey))
		return nil
	})

	//向本地数据库定入区块
	block := NewBlock(txs, prevBlockHash)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("添加区块时，无法找到数据库")
			os.Exit(1)
		}

		err := bucket.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic("向数据库添加区块失败1：", err)
		}
		err = bucket.Put([]byte(lastHashKey), block.Hash)
		if err != nil {
			log.Panic("向数据库添lastHash失败：", err)
		}
		bc.tail = block.Hash
		return nil
	})

	if err != nil {
		log.Panic("添加区块信息失败：", err)
	}
}

//查找指定地址能够支付的utxo的交易集合
func (bc *BlockChain) FindUTXOTransactions(address string) []Transaction {
	//包含目标utxo的交易集合
	var UTXOTransactions []Transaction
	//存储使用过的utxo的集合 map[交易id] []int64
	spentUTXO := make(map[string][]int64)

	//根据当前区块找到前一个区块
	it := bc.NewIterator()
	for {
		//遍历区块
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transactions {
			//如果不是挖矿交易
			if !tx.IsCoinbase() {
				//遍历input 目的：找到已经消耗过的utxo,
				//需要两个字段来标识使用过的utxo：1.交易id 2.outputs的索引
				for _, input := range tx.TXInputs {
					if input.CanUnlockUTXOWith(address) {
						spentUTXO[string(input.TXID)] = append(spentUTXO[string(input.TXID)], input.Vout)
					}
				}
			}

		OUTPUTS:
		//遍历output 目的：找到所有能支配的utxo
			for currIndex, output := range tx.TXOutputs {
				//检查当前output是否已经消费，如果消费过，那到就进行下一个output检验
				if spentUTXO[string(tx.TXID)] != nil {
					//非空，代表当前交易有消费的utxo
					indexes := spentUTXO[string(tx.TXID)]
					for _, index := range indexes {
						//当前的索引和消费的索引比较：相同表明output消耗了
						//跳过，进行下一个output判断
						if int64(currIndex) == index {
							continue OUTPUTS
						}
					}
				}

				//如果当前地址是这个utxo的所有者，就满足条件
				if output.CanBeUnlockedWith(address) {
					UTXOTransactions = append(UTXOTransactions, *tx)
				}
			}
		}

		if len(block.PervHash) == 0 {
			break
		}
	}

	return UTXOTransactions
}

//寻找录前地址能够使用的utxo
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXOs []TXOutput
	txs := bc.FindUTXOTransactions(address)
	//遍历可用的交易，查找可以的output =>utxo
	for _, tx := range txs {
		for _, utxo := range tx.TXOutputs {
			//当前地址可能的output的(是否能解锁)
			if utxo.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, utxo)
			}
		}
	}
	return UTXOs
}

//检查数据文件是否存在
func isDBExist() bool {
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//区块链初始化 创建数据文件
func InitBlockChain(address string) *BlockChain {
	//Ver4：改写命令行参数模式
	if isDBExist() {
		fmt.Println("数据文件已存在，无需创建")
		os.Exit(1)
	}
	//Ver3：使用数据改写
	var lastHash []byte

	//打开数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic("打开区块链数据" + dbFile + "出错")
	}

	db.Update(func(tx *bolt.Tx) error {
		//没有buket，需要创建，并且创建一个创世块
		coinbase := NewCoinbaseTx(address, genesisInfo) //todo
		genesis := NewGenesisBlock(coinbase)
		bucket, err := tx.CreateBucket([]byte(blockBucket))
		if err != nil {
			log.Panic("创建bolt数据失败：", err)
		}

		err = bucket.Put(genesis.Hash, genesis.Serialize()) //TODO
		if err != nil {
			log.Panic("向数据添加创世区块失败：", err)
		}
		err = bucket.Put([]byte(lastHashKey), genesis.Hash)
		if err != nil {
			log.Panic("向数据添加创世区块lastHashKey失败：", err)
		}

		//向内存更新创世区块的hash
		lastHash = genesis.Hash
		return nil
	})
	return &BlockChain{db, lastHash}
}

//获取区块句柄
func GetBlockChainHandler() *BlockChain {
	//Ver4：改写命令行参数模式
	if !isDBExist() {
		fmt.Println("数据文件不存，请先创建区块链文件")
		os.Exit(1)
	}
	//Ver3：使用数据改写
	var lastHash []byte

	//打开数据库
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic("打开区块链数据" + dbFile + "出错")
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil {
			//如果有，取出最后一个区块的hash值
			lastHash = bucket.Get([]byte(lastHashKey))
		}
		return nil
	})

	return &BlockChain{db, lastHash}
}

//迭代器 用于实现对区块链的遍历
type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

//创建迭代器，初始化为指向最后一个区块
func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{bc.tail, bc.db}
}

//迭代器指针，下移一个元素
func (it *BlockChainIterator) Next() (block *Block) {
	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil {
			data := bucket.Get(it.currentHash)
			//根据前一个区块的哈希值 ，获取该区块
			block = Deserialize(data)
			//迭代器的指针指向下一移 前一区块的前一区块的hash
			it.currentHash = block.PervHash
		}
		return nil
	})
	return
}

//所需要的合适的utxo集合
//validUTXOs := make(map[string][]int64)
//返回utxo的金额总和
//total float64

//从所有的交易里找到合适的utxo 和总金额
func (bc *BlockChain) FindSuitableUTXOs(address string, amount float64) (map[string][]int64, float64) {
	txs := bc.FindUTXOTransactions(address)
	validUTXOs := make(map[string][]int64)
	var total float64

	FIND:
	//遍历交易
	for _, tx := range txs {
		outputs := tx.TXOutputs
		for index, output := range outputs {
			if output.CanBeUnlockedWith(address) {
				//判断当前收集的utxo金额是否大于所需要花费进帐的金额
				if total < amount {
					total += output.Value
					validUTXOs[string(tx.TXID)] = append(validUTXOs[string(tx.TXID)], int64(index))
				} else {
					break FIND
				}
			}
		}
	}

	return validUTXOs, total
}
