package BLC

import(
	"time"
	"fmt"
	"encoding/hex"
	"bytes"
	"encoding/gob"
	"log"
	"github.com/boltdb/bolt"
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

// 将区块序列化成字节数组
func (block *Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func DeserializeBlock(blockBytes []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewBlock(Data []byte, blockchain *BlockChain){


	err := blockchain.BlockDB.Update(func(tx *bolt.Tx) error{

		//1. 获取表
		b := tx.Bucket([]byte(blockTableName))
		//2. 创建新区块
		if b != nil {

			// ⚠️，先获取最新区块
			blockBytes := b.Get(blockchain.Tip)

			// 反序列化
			block := DeserializeBlock(blockBytes)
			//fmt.Println("it's in New Block------------------")
			//PrintBlock(block)

			//3. 将区块序列化并且存储到数据库中
			newBlock :=&Block{block.Height+1, block.BlockHash,1,Data,time.Now().Unix(),nil}

			Pow:=ProofOfWork(newBlock)
			nonce, hash := Pow.run()

			newBlock.BlockHash=hash[:]
			//fmt.Println("it's hash+++++++++++,%s,%s",hex.EncodeToString(hash),hex.EncodeToString(newBlock.BlockHash))
			newBlock.nonce=nonce
			//fmt.Println("it's nonce+++++++++++,%d",nonce)
			//fmt.Println("it's newBlock nonce+++++++++++,%d",newBlock.nonce)
			err := b.Put(newBlock.BlockHash,newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//4. 更新数据库里面"l"对应的hash
			err = b.Put([]byte("l"),newBlock.BlockHash)
			if err != nil {
				log.Panic(err)
			}
			//5. 更新blockchain的Tip
			blockchain.Tip = newBlock.BlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

func NewGenesisBlock(Data []byte) * Block{

	block := &Block{1,IntToHex(0),1,Data,time.Now().Unix(),nil}

	Pow:=ProofOfWork(block)
	nonce, hash := Pow.run()

	block.BlockHash=hash[:]
	block.nonce=nonce

	//PrintBlock(block)
	return block
}