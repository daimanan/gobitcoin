//工具方法
package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"crypto/sha256"
)

//将Uint64转化为字符切片
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	//任意的数据转换成byte字节流
	err := binary.Write(&buffer, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

//区块转Hash值
func BlockToHash(block *Block) []byte {
	//将属性转化为字节切片
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PervHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Hash,
		block.Data,
	}

	//合并字节流
	blockInfo := bytes.Join(tmp, []byte{})
	//计算哈希值
	hash := sha256.Sum256(blockInfo)
	return hash[:]
}
