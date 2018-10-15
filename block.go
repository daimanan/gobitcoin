package main

//区块结构
type Block struct {
	Version    uint64 //版本号
	PervHash   []byte //前一个区块的哈希值
	MerkelRoot []byte //梅克尔根（暂时不实现）
	TimeStamp  uint64 //时间戳
	Difficulty uint64 //难度值（挖矿生成哈希值的难度 ）
	Nonce      uint64 //随机数

	Hash       []byte //哈希值
	Data       []byte //交易数据
}

//生成哈希
func (bc * Block) SetHash()  {

}