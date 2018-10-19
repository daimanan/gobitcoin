package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

type blockInterface interface {
	SetHash()          //ver1 ver2时使用
	Serialize() []byte //Ver3:序列化
}

//区块结构
type Block struct {
	Version    uint64 //版本号
	PervHash   []byte //前一个区块的哈希值
	MerkelRoot []byte //梅克尔根（暂时不实现）
	TimeStamp  uint64 //时间戳
	Difficulty uint64 //难度值（挖矿生成哈希值的难度 ）
	Nonce      uint64 //随机数
	Hash       []byte //哈希值

	//Ver4之前使用
	//Data []byte //交易数据

	//Ver4：使用交易改写
	Transactions []*Transaction //交易信息
}

//创建区块(构造)
//Ver4：使用交易改写生成区块
func NewBlock(txs []*Transaction, pervHash []byte) *Block {
	block := Block{
		Version:    00,       //版本为00正式网络 01测试网络
		PervHash:   pervHash, //前一区块的哈希(传值过来)
		MerkelRoot: []byte{}, //先设为空
		//设为当前系统的时间
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0, //开发与测试设为0
		Nonce:      0, //开发与测试设为0

		//ver4之前使用
		//Hash:       []byte{}, //初始化为空
		//Data:       []byte(data),

		//Ver4：使用交易改写
		Transactions: txs,
	}

	//引入算力
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//序列化
func (bc *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&bc)
	if err != nil {
		log.Panic("区块编码失败：", err)
	}
	return buffer.Bytes()
}

//反序列化
func Deserialize(data []byte) *Block {
	if len(data) == 0 || data == nil {
		return nil
	}

	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("区块解码失败：", err)
	}

	return &block
}

//生成哈希方式：
//	1.将所有属性拼装起来
//	2.再进行hash运算
func (bc *Block) SetHash() {
	/*
	//将属性转化为字节切片
	tmp := [][]byte{
		Uint64ToByte(bc.Version),
		bc.PervHash,
		bc.MerkelRoot,
		Uint64ToByte(bc.TimeStamp),
		Uint64ToByte(bc.Difficulty),
		Uint64ToByte(bc.Nonce),
		bc.Hash,
		bc.Data,
	}

	//合并字节流
	blockInfo := bytes.Join(tmp, []byte{})
	//计算哈希值
	hash := sha256.Sum256(blockInfo)
	//设置当前区块的哈希
	bc.Hash = hash[:]
	*/
	bc.Hash = BlockToHash(bc)
}

//创世区块生成
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}) //创世区块的前区块是空
}
