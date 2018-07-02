package BLC

import(
	"time"
	"fmt"
	"encoding/hex"
)
type Block struct{
	Height int64

	PrevBlockHash []byte

	nonce int64

	Data []byte

	Timestamp int64

	BlockHash []byte
}

func PrintBlock(block *Block){
	fmt.Println("===============")
	fmt.Println(block.Height)
	fmt.Println(hex.EncodeToString(block.PrevBlockHash))
	fmt.Println(block.nonce)
	fmt.Println(string(block.Data))
	fmt.Println(block.Timestamp)
	fmt.Println(hex.EncodeToString(block.BlockHash))
	fmt.Println("===============")
}

func NewBlock(Height int64, PrevBlockHash []byte, Data []byte) * Block{
	block := &Block{Height,PrevBlockHash,0,Data,time.Now().Unix(),nil}

	Pow:=ProofOfWork(block)
	nonce, hash := Pow.run()

	block.BlockHash=hash[:]
	block.nonce=nonce

	PrintBlock(block)

	return block
}

func NewGenesisBlock(Data []byte) * Block{

	block := &Block{1,IntToHex(0),1,Data,time.Now().Unix(),nil}

	Pow:=ProofOfWork(block)
	nonce, hash := Pow.run()

	block.BlockHash=hash[:]
	block.nonce=nonce
	PrintBlock(block)
	return block
}