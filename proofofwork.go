package main

import (
	"math/big"
	"fmt"
	"crypto/sha256"
)

type proofOfWorkInterface interface {
	Run() (hash []byte, nonce uint64)
	IsValid() bool
}

//工作量证明
type ProofOfWork struct {
	block  *Block   //区块
	target *big.Int //目标值(控制挖矿的难度)
}

//运行工作证明(挖矿方法)
func (pow *ProofOfWork) Run() (hash []byte, nonce uint64) {
	var tmpBlock = pow.block

	fmt.Println("开始挖矿......")
	for {
		tmpBlock.Nonce = nonce
		//将区块转化为hash，将hash转化为big.Int
		hash = BlockToHash(tmpBlock)
		var tmpInt big.Int
		tmpInt.SetBytes(hash)
		//判断生成的hash值与挖矿的难度值之间的大小
		//如果当前的hash值 < 挖矿的难度值 找到了
		if tmpInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！hash : %x, nonce : %d\n", hash, nonce)
			break
		} else { //循环继教找
			nonce++
		}
	}

	return hash, nonce
}

//哈希值有效性的校验
func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int
	//将当前的区块随机数转换为[]byte
	data := Uint64ToByte(pow.block.Nonce)
	//对data进行hash计算
	hash := sha256.Sum256(data)
	//将hash值转换为hastInt
	hashInt.SetBytes(hash[:])

	//校验生成当前区块的哈希是否为效(小于挖矿的目标值)
	return hashInt.Cmp(pow.target) == -1
}

//区块工作证明构造
func NewProofOfWork(block *Block) *ProofOfWork {
	var pow ProofOfWork
	pow.block = block

	//难度值，先写成固定值
	targetString := "0000100000000000000000000000000000000000000000000000000000000000"
	var bigIntTemp big.Int
	//将字符转化为bigInt型
	bigIntTemp.SetString(targetString, 16)

	pow.target = &bigIntTemp

	return &pow
}
