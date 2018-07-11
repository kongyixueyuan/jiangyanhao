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
	"strconv"
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
		block := blockchainIterator.jyh_Next()

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
				fmt.Println("------")
				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
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

func (bclockchain *BlockChain) SignTransaction(tx *Transaction,privKey ecdsa.PrivateKey,txs []*Transaction)  {

	if tx.jyh_IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.TXIn {
		prevTX, err := bclockchain.FindTransaction(vin.TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	tx.jyh_Sign(privKey, prevTXs)

}


func (bc *BlockChain) FindTransaction(ID []byte,txs []*Transaction) (Transaction, error) {

	for _,tx := range txs  {
		if bytes.Compare(tx.TxHash, ID) == 0 {
			return *tx, nil
		}
	}


	bci := bc.Iterator()

	for {
		block := bci.jyh_Next()

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
func (bc *BlockChain) VerifyTransaction(tx *Transaction,txs []*Transaction) bool {


	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.TXIn {
		prevTX, err := bc.FindTransaction(vin.TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.TxHash)] = prevTX
	}

	return tx.jyh_Verify(prevTXs)
}

/*
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

*/
// 查询余额
func (blockchain *BlockChain) GetBalance(address string) int64 {

	utxos := blockchain.UnUTXOs(address,[]*Transaction{})

	var amount int64

	for _, utxo := range utxos {

		amount = amount + utxo.Output.Value
	}

	return amount
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *BlockChain) UnUTXOs(address string,txs []*Transaction) []*UTXO {



	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range txs {

		if tx.jyh_IsCoinbaseTransaction() == false {
			for _, in := range tx.TXIn {
				//是否能够解锁
				publicKeyHash := Base58Decode([]byte(address))

				ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]
				if in.IsLockedPubkeyTxIn(ripemd160Hash) {

					key := hex.EncodeToString(in.TxHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}

			}
		}
	}


	for _,tx := range txs {

	Work1:
		for index,out := range tx.TXOut {

			if out.IsLockedPubkeyTxOut(address) {
				fmt.Println("看看是否是俊诚...")
				fmt.Println(address)

				fmt.Println(spentTXOutputs)

				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}


	blockIterator := blockchain.Iterator()

	for {

		block := blockIterator.jyh_Next()

		fmt.Println(block)
		fmt.Println()

		for i := len(block.Transaction) - 1; i >= 0 ; i-- {

			tx := block.Transaction[i]
			// txHash
			// Vins
			if tx.jyh_IsCoinbaseTransaction() == false {
				for _, in := range tx.TXIn {
					//是否能够解锁
					publicKeyHash := Base58Decode([]byte(address))

					ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]

					if in.IsLockedPubkeyTxIn(ripemd160Hash) {

						key := hex.EncodeToString(in.TxHash)

						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.TXOut {

				if out.IsLockedPubkeyTxOut(address) {

					fmt.Println(out)
					fmt.Println(spentTXOutputs)

					//&{2 zhangqiang}
					//map[]

					if spentTXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

						if len(spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {

								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)

							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}

		}

		fmt.Println(spentTXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
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
//1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM
//13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ
//1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG
//./main createBlockchain "13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"
func NewBlockChain(block * Block) * BlockChain{
	// 判断数据库是否存在
	if DBExists() {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

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
		//fmt.Println("in new chain",DeserializeBlock(ss).Nonce)
		err = b.Put(block.BlockHash, block.jyh_Serialize())
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

//get UTXO from total tables
func (blc *BlockChain)jyh_FindUTXOMap() map[string]*TXOutputs{
	blcIterator := blc.Iterator()

	// 存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*TXInput)


	utxoMaps := make(map[string]*TXOutputs)
	for {
		block := blcIterator.jyh_Next()


		for i := len(block.Transaction) - 1; i >= 0 ;i-- {

			txOutputs := &TXOutputs{[]*UTXO{}}

			tx := block.Transaction[i]


			// coinbase
			if tx.jyh_IsCoinbaseTransaction() == false {
				for _,txInput := range tx.TXIn {

					txHash := hex.EncodeToString(txInput.TxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash],txInput)

				}
			}



			txHash := hex.EncodeToString(tx.TxHash)

		WorkOutLoop:
			for index,out := range tx.TXOut  {

				if tx.jyh_IsCoinbaseTransaction() {

					fmt.Println("IsCoinbaseTransaction")
					fmt.Println(out)
					fmt.Println(txHash)
				}

				txInputs := spentableUTXOsMap[txHash]

				if len(txInputs) > 0 {

					isSpent := false

					for _,in := range  txInputs {

						outPublicKey := out.Ripemd160Hash
						inPublicKey := in.PubKey
						fmt.Println("outPubkey:%s",out.Ripemd160Hash)
						fmt.Println("inPubkey:%s",in.PubKey)
						if bytes.Compare(outPublicKey,HashPubKey(inPublicKey)) == 0{
							if index == in.Vout {
								isSpent = true
								continue WorkOutLoop
							}
						}

					}

					if isSpent == false {
						utxo := &UTXO{tx.TxHash,index,out}
						txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
					}

				} else {
					utxo := &UTXO{tx.TxHash,index,out}
					txOutputs.UTXOS = append(txOutputs.UTXOS,utxo)
				}

			}

			// 设置键值对
			utxoMaps[txHash] = txOutputs

		}


		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return utxoMaps
}
	/*./main createBlockchain "1JnMQYJaXFkWXCjFrCokMqfk66J7FYmWbn"

func (blc *BlockChain)FindUTXOMap() map[string]*TXOutputs{
	UTXOMap := make(map[string]*TXOutputs)
	spentableUTXOsMap := make(map[string][]*TXInput)

	blockchainIterator := blc.Iterator()

	for {
		block := blockchainIterator.Next()
		for _, tx := range block.Transaction {
			for i, TxOutput := range tx.TXOut {
				UTXOMap[string(tx.TxHash)].UTXOS=append(UTXOMap[string(tx.TxHash)].UTXOS,&UTXO{tx.TxHash, i, TxOutput })
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)


		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}

	}

	blockchainIterator = blc.Iterator()
	//var temp []int
	for {
		block := blockchainIterator.Next()
		////input:[Txhash][1]------>[Txhash][0,2]
		for _, tx := range block.Transaction {
			//fmt.Println("it's typing in:")
			for _, in := range tx.TXIn {
				for RemoveIndex , TxOutput := range UTXOMap[string(in.TxHash)].UTXOS{
					if(in.Vout==TxOutput.Index){
						UTXOMap[string(in.TxHash)] = append(UTXOMap[string(in.TxHash)][:RemoveIndex], UTXOMap[string(in.TxHash)][RemoveIndex+1:]...)
						//spendableUTXO[string(in.TxHash)]=RemoveInt(spendableUTXO[string(in.TxHash)], value)
					}
				}
			}
		}
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return UTXOMap*/

// 转账时查找可用的UTXO
func (blockchain *BlockChain) jyh_FindSpendableUTXOS(from string, amount int,txs []*Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO

	utxos := blockchain.UnUTXOs(from,txs)

	spendableUTXO := make(map[string][]int)

	//2. 遍历utxos

	var value int64

	for _, utxo := range utxos {

		value = value + utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)

		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {

		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}


//./main send -from '["1JnMQYJaXFkWXCjFrCokMqfk66J7FYmWbn"]' -to '["14y6Dxyyo6Eh3JYFMVRVG6wK3M3Sm6m4Qo"]' -amount '["2"]'
// 挖掘新的区块
func (blockchain *BlockChain) jyh_MineNewBlock(from []string, to []string, amount []string) {
	//	$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//1.建立一笔交易


	utxoSet := &jyh_UTXOSet{blockchain}

	var txs []*Transaction

	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := jyh_NewSimpleTransaction(address, to[index], int64(value), utxoSet,txs)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	//奖励
	tx := jyh_SetCoinbBaseTransaction(from[0])
	txs = append(txs,tx)


	//1. 通过相关算法建立Transaction数组
	var block *Block

	blockchain.BlockDB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = jyh_DeserializeBlock(blockBytes)

		}

		return nil
	})
	//	./main send -from '["183UkSrFV72YZ53FS5MAxiSPgNrgWv9uU8"]' -to '["1Ap11ovDK7aorb1Cmkb3L6JxhLZDCC2sbN"]' -amount '["3"]'

	// 在建立新区块之前对txs进行签名验证

	_txs := []*Transaction{}

	for _,tx := range txs  {

		if blockchain.VerifyTransaction(tx,_txs) != true {
			log.Panic("ERROR: Invalid transaction")
		}

		_txs = append(_txs,tx)
	}


	//2. 建立新的区块
	block = jyh_NewBlock(txs, block.Height+1, block.BlockHash)

	//将新区块存储到数据库
	blockchain.BlockDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			b.Put(block.BlockHash, block.jyh_Serialize())

			b.Put([]byte("l"), block.BlockHash)

			blockchain.Tip = block.BlockHash

		}
		return nil
	})

}