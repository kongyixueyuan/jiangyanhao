package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"time"
	"math/big"
)

type BlockChain struct{
	Tip []byte //最新的区块的Hash
	BlockDB  *bolt.DB
}


func (blockchain *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip,blockchain.BlockDB}
}

// 遍历输出所有区块的信息
func (blc *BlockChain) Printchain()  {

	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()

		fmt.Printf("Height：%d\n",block.Height)
		fmt.Printf("PrevBlockHash：%x\n",block.PrevBlockHash)
		fmt.Printf("Timestamp：%s\n",time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n",block.BlockHash)
		fmt.Printf("Nonce：%d\n",block.nonce)
		fmt.Printf("Nonce：%x\n",IntToHex(block.nonce))

		fmt.Println("Txs:")
		for _, tx := range block.Transaction {

			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.TXIn {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.TXOut {
				fmt.Println(out.Value)
				fmt.Println(out.ScriptPubKey)
			}
		}


		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

	//	time.Sleep(1 * time.Second)
		if big.NewInt(0).Cmp(&hashInt) == 0{
			break;
		}
	}

}

// 数据库名字
const dbName  = "blockchain.db"

// 表的名字
const blockTableName  = "blocks"

func BlockchainObject() *BlockChain {
	db,err:= bolt.Open(dbName, 0600,nil)
	Tip := *new([]byte)
	if err!=nil{
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b==nil{
			log.Println("no tables")
		}
		Tip = b.Get([]byte("l"))
		if(Tip==nil){
			log.Println("no Tip")
		}
		return nil
	})


	return &BlockChain{Tip,db}
}

func NewBlockChain(block * Block) * BlockChain{

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		//  获取表
		b := tx.Bucket([]byte(blockTableName))

		if b == nil {
			// 创建数据库表
			b,err = tx.CreateBucket([]byte(blockTableName))

			if err != nil {
				log.Panic(err)
			}
		}

		err = b.Put(block.BlockHash, block.Serialize())
		if err != nil{
			log.Panic(err)
		}

		err = b.Put([]byte("l"), block.BlockHash)
		if err != nil{
			log.Panic(err)
		}

		return nil
	})
	blockchain := &BlockChain{[]byte(block.BlockHash),db}

	return blockchain
}

func AddToChain(block *Block, blockchain *BlockChain) *BlockChain{

	return blockchain
}