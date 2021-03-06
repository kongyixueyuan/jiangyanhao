package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"time"
	"math/big"
	"os"
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
		fmt.Printf("Nonce：%d\n",block.Nonce)
		//fmt.Println(block.Nonce)
		//fmt.Printf("Nonce：%x\n",IntToHex(block.nonce))

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

// 遍历输出所有区块的信息
func (blc *BlockChain) getBalance()  {

	blockchainIterator := blc.Iterator()


	//把可花费的utxo算出来，都是output
	//[Txhash][0,2]
	spendableUTXO := FindSpendableUTXO(blc)
	//var spendableTxout []*TXOutput
	spendableMap := make(map[string][]*TXOutput)

	balances := make(map[string]int64)
	for {
		block := blockchainIterator.Next()

		for _, tx := range block.Transaction {
			//1、计算block-spendable
			//TXOut []*TXOutput
			for i, txOut := range tx.TXOut {

				for _, value := range spendableUTXO[string(tx.TxHash)] {
					if (i == value) {
						//	spendableTxout=append(spendableTxout, txOut)
						spendableMap[string(tx.TxHash)] = append(spendableMap[string(tx.TxHash)], txOut)
					}
				}
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

	for _, Txs := range spendableMap {
		for _, txOut := range Txs {
			balances[txOut.ScriptPubKey]=balances[txOut.ScriptPubKey]+txOut.Value
		}
	}

	for name, balance := range balances{
		fmt.Println(name,":",balance)
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
		//fmt.Println("new blockchain:%d",block.Nonce)
		ss:=block.Serialize()
		//fmt.Println("in new chain",DeserializeBlock(ss).Nonce)
		err = b.Put(block.BlockHash, ss)
		//fmt.Println("ss deserialize",DeserializeBlock(ss).Nonce)
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

// 判断数据库是否存在
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}