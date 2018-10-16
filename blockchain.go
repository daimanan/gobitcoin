package main

type blockChainInterface interface {
	AddBlock(data string)
}

//区块链
type BlockChain struct {
	blocks []*Block //用区块切片模拟区块链
}

//添加区块到区块链中
func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.blocks[len(bc.blocks)-1]
	block := NewBlock(data, lastBlock.Hash)
	block.SetHash()
	bc.blocks = append(bc.blocks, block)
}

//区块链构造函数
func NewBlockChain() *BlockChain {
	block := NewBlock("创世区块", []byte{})
	block.SetHash()
	blocks := []*Block{block}
	return &BlockChain{blocks}
}
