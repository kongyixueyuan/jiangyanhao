package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"time"
	"math/big"
	"os"
	"encoding/hex"
	"crypto/ecdsa"
	"bytes"
)

type BlockChain struct{
	Tip []byte //最新的区块的Hash
	BlockDB  *bolt.DB
}

//1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG
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

		fmt.Println("Txs:")
		for _, tx := range block.Transaction {

			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.TXIn {
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%x\n", in.Signature)
				fmt.Println("------")
				fmt.Printf("%x\n", in.PubKey)
				fmt.Println("------")
				//fmt.Printf("%s\n", in.ScriptSig)
			}

			fmt.Println("Vouts:")
			for _, out := range tx.TXOut {
				fmt.Println(out.Value)
				fmt.Printf("%x\n", out.Ripemd160Hash)
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

func (bclockchain *BlockChain) SignTransaction(tx *Transaction,privKey ecdsa.PrivateKey)  {

	if tx.IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.TXIn {
		prevTX, err := bclockchain.FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	tx.Sign(privKey, prevTXs)

}


func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {

	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transaction {
			if bytes.Compare(tx.TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)


		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return Transaction{},nil
}


// 验证数字签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {


	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.TXIn {
		prevTX, err := bc.FindTransaction(vin.TxHash)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	return tx.Verify(prevTXs)
}

// 遍历输出所有区块的信息
func (blc *BlockChain) getBalance()  {

	blockchainIterator := blc.Iterator()


	//把可花费的utxo算出来，都是output
	//[Txhash][0,2]
	spendableUTXO := FindSpendableUTXO(blc)
	for tx, nums := range spendableUTXO {
		fmt.Printf(":txout:%x----", tx)
		for i := range nums {
			fmt.Printf("%d---", nums[i])
		}
		fmt.Printf("\n")
	}
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

		if big.NewInt(0).Cmp(&hashInt) == 0{
			break;
		}
	}

	for _, Txs := range spendableMap {
		for _, txOut := range Txs {
			balances[string(txOut.Ripemd160Hash)]=balances[string(txOut.Ripemd160Hash)]+txOut.Value
		}
	}

	for name, balance := range balances{
		x:=[]byte(name)
		fmt.Printf("%x",x)
		fmt.Println(":",balance)
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