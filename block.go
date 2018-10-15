package main

import (
	"time"
	"bytes"
	"encoding/binary"
	"log"
	"crypto/sha256"
)

type blockInterface interface {
	NewBlock(data string, pervHash []byte) *Block
	SetHash()
}

//区块结构
type Block struct {
	Version    uint64 //版本号
	PervHash   []byte //前一个区块的哈希值
	MerkelRoot []byte //梅克尔根（暂时不实现）
	TimeStamp  uint64 //时间戳
	Difficulty uint64 //难度值（挖矿生成哈希值的难度 ）
	Nonce      uint64 //随机数

	Hash []byte //哈希值
	Data []byte //交易数据
}

//创建区块(构造)
func (bc *Block) NewBlock(data string, pervHash []byte) *Block {
	block := Block{
		Version:    00,       //版本为00正式网络 01测试网络
		PervHash:   pervHash, //前一区块的哈希(传值过来)
		MerkelRoot: []byte{}, //先设为空
		//设为当前系统的时间
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,        //开发与测试设为0
		Nonce:      0,        //开发与测试设为0
		Hash:       []byte{}, //初始化为空
		Data:       []byte(data),
	}
	return &block
}

//将Uint64转化为字符切片
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	//任意的数据转换成byte字节流
	err := binary.Write(&buffer, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}
	return nil
}

//生成哈希方式：
//	1.将所有属性拼装起来
//	2.再进行hash运算
func (bc *Block) SetHash() {
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
}
