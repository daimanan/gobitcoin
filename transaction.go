package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
)

//挖矿交易收易
const reward = 12.5

type transactionInterface interface {
	SetTXID()
	IsCoinbase() bool
}

//比特币交易
type Transaction struct {
	TXID      []byte     //交易ID
	TXInputs  []TXInput  //输入
	TXOutputs []TXOutput //输出
}

//生成交易的ID (交易自己的Hash值)
func (tx *Transaction) SetTXID() {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic("生成交易ID转码失败\n", err)
	}

	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}

//是否挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	if len(tx.TXInputs) == 1 { //打包的区块中，至少有一个输入
		//交易的id 是 0 并且 输入索引没有
		if len(tx.TXInputs[0].TXID) == 0 && tx.TXInputs[0].Vout == -1 {
			return true
		}
	}

	return false
}

//交易输入
type TXInput struct {
	TXID []byte //引用输出的交易ID
	Vout int64  //引用输出的索引
	//解锁脚本 ，指明可以使用某个output的格格件
	ScriptSig string
}

//检查当前用䚮能否解开引用的utxo
func (input *TXInput) CanUnlockUTXOWith(unlockData string) bool {
	return input.ScriptSig == unlockData
}

//交易输出
type TXOutput struct {
	Value float64 //支付给收款方的金额
	//锁定脚要，指定收款方的地址
	ScriptPubKey string
}

//检查当瓣用䚮时候是否这个utxo的所有者
func (output *TXOutput) CanBeUnlockedWith(unlockData string) bool {
	return output.ScriptPubKey == unlockData
}

//创建coinbae交易，只有收款人，没有付款人，是矿工的挖矿交易
func NewCoinbaseTx(address string, data string) *Transaction {
	//提示：矿工挖矿信息
	if data == "" {
		data = fmt.Sprintf("奖劢矿工 %s %d btc", address, reward)
	}

	input := TXInput{[]byte{}, -1, data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetTXID()

	return &tx //todo
}
