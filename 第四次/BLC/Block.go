package BLC

import(
	"time"
	"fmt"
	"encoding/hex"
	"bytes"
	"encoding/gob"
	"log"
	"github.com/boltdb/bolt"
	"crypto/sha256"
)
type Block struct{
	Height int64

	PrevBlockHash []byte

	nonce int64

	Transaction []*Transaction

	Timestamp int64

	BlockHash []byte
}

func PrintBlock(block *Block){
	fmt.Println("===============")
	fmt.Println(block.Height)
	fmt.Println(hex.EncodeToString(block.PrevBlockHash))
	fmt.Println(block.nonce)
	fmt.Println(block.Timestamp)
	fmt.Println(hex.EncodeToString(block.BlockHash))
	fmt.Println("===============")
}

// 需要将Txs转换成[]byte
func (block *Block) HashTransactions() []byte  {


	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transaction {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]

}

// 将区块序列化成字节数组
func (block *Block) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	fmt.Println("in serialize, %d",block.nonce)
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
	fmt.Println("in deserialize, %d",block.nonce)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewBlock(Transaction []byte, blockchain *BlockChain){


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
			newBlock :=&Block{block.Height+1, block.BlockHash,1,nil,time.Now().Unix(),nil}

			Pow:=ProofOfWork(newBlock)
			nonce, hash := Pow.run()

			newBlock.BlockHash=hash[:]
			//fmt.Println("it's hash+++++++++++,%s,%s",hex.EncodeToString(hash),hex.EncodeToString(newBlock.BlockHash))
			newBlock.nonce=nonce
			//fmt.Println("it's nonce+++++++++++,%d",nonce)
			fmt.Println("it's newBlock nonce+++++++++++,%d",newBlock.nonce)
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

func NewGenesisBlock(address string) * Block{
	GenesisTransaction:=SetCoinbBaseTransaction("jiang")
	block := &Block{1,IntToHex(0),1,GenesisTransaction,time.Now().Unix(),nil}

	Pow:=ProofOfWork(block)
	nonce, hash := Pow.run()

	block.BlockHash=hash[:]
	block.nonce=nonce

	//PrintBlock(block)
	return block
}