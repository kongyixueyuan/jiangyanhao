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

type Jyh_BlockChain struct{
	Jyh_Tip []byte //最新的区块的Hash
	Jyh_BlockDB  *bolt.DB
}

//1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG
func (blockchain *Jyh_BlockChain) Jyh_Iterator() *Jyh_BlockchainIterator {
	return &Jyh_BlockchainIterator{blockchain.Jyh_Tip,blockchain.Jyh_BlockDB}
}

// 遍历输出所有区块的信息
func (blc *Jyh_BlockChain) Jyh_Printchain()  {

	blockchainIterator := blc.Jyh_Iterator()

	for {
		block := blockchainIterator.Jyh_Next()

		fmt.Printf("Height：%d\n",block.Jyh_Height)
		fmt.Printf("PrevBlockHash：%x\n",block.Jyh_PrevBlockHash)
		fmt.Printf("Timestamp：%s\n",time.Unix(block.Jyh_Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash：%x\n",block.Jyh_BlockHash)
		fmt.Printf("Nonce：%d\n",block.Jyh_Nonce)

		fmt.Println("Txs:")
		for _, tx := range block.Jyh_Transaction {

			fmt.Printf("%x\n", tx.Jyh_TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Jyh_TXIn {
				fmt.Println("------")
				fmt.Printf("%x\n", in.Jyh_TxHash)
				fmt.Printf("%d\n", in.Jyh_Vout)
				fmt.Println("------")
				//fmt.Printf("%s\n", in.ScriptSig)
			}
			fmt.Println("Vouts:")
			for _, out := range tx.Jyh_TXOut {
				fmt.Println(out.Jyh_Value)
				fmt.Printf("%x\n", out.Jyh_Ripemd160Hash)
			}
		}


		var hashInt big.Int
		hashInt.SetBytes(block.Jyh_PrevBlockHash)

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

func (bclockchain *Jyh_BlockChain) Jyh_SignTransaction(tx *Jyh_Transaction,privKey ecdsa.PrivateKey,txs []*Jyh_Transaction)  {

	if tx.Jyh_IsCoinbaseTransaction() {
		return
	}

	prevTXs := make(map[string]Jyh_Transaction)

	for _, vin := range tx.Jyh_TXIn {
		prevTX, err := bclockchain.Jyh_FindTransaction(vin.Jyh_TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Jyh_TxHash)] = prevTX
	}

	tx.Jyh_Sign(privKey, prevTXs)

}


func (bc *Jyh_BlockChain) Jyh_FindTransaction(ID []byte,txs []*Jyh_Transaction) (Jyh_Transaction, error) {

	for _,tx := range txs  {
		if bytes.Compare(tx.Jyh_TxHash, ID) == 0 {
			return *tx, nil
		}
	}


	bci := bc.Jyh_Iterator()

	for {
		block := bci.Jyh_Next()

		for _, tx := range block.Jyh_Transaction {
			if bytes.Compare(tx.Jyh_TxHash, ID) == 0 {
				return *tx, nil
			}
		}

		var hashInt big.Int
		hashInt.SetBytes(block.Jyh_PrevBlockHash)


		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}
	}

	return Jyh_Transaction{},nil
}


// 验证数字签名
func (bc *Jyh_BlockChain) Jyh_VerifyTransaction(tx *Jyh_Transaction,txs []*Jyh_Transaction) bool {


	prevTXs := make(map[string]Jyh_Transaction)

	for _, vin := range tx.Jyh_TXIn {
		prevTX, err := bc.Jyh_FindTransaction(vin.Jyh_TxHash,txs)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Jyh_TxHash)] = prevTX
	}

	return tx.Jyh_Verify(prevTXs)
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
func (blockchain *Jyh_BlockChain) GetBalance(address string) int64 {

	utxos := blockchain.Jyh_UnUTXOs(address,[]*Jyh_Transaction{})

	var amount int64

	for _, utxo := range utxos {

		amount = amount + utxo.Jyh_Output.Jyh_Value
	}

	return amount
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
func (blockchain *Jyh_BlockChain) Jyh_UnUTXOs(address string,txs []*Jyh_Transaction) []*Jyh_UTXO {

	var unUTXOs []*Jyh_UTXO

	Jyh_spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range txs {

		if tx.Jyh_IsCoinbaseTransaction() == false {
			for _, in := range tx.Jyh_TXIn {
				//是否能够解锁
				publicKeyHash := Jyh_Base58Decode([]byte(address))

				ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]
				if in.Jyh_IsLockedPubkeyTxIn(ripemd160Hash) {

					key := hex.EncodeToString(in.Jyh_TxHash)

					Jyh_spentTXOutputs[key] = append(Jyh_spentTXOutputs[key], in.Jyh_Vout)
				}

			}
		}
	}


	for _,tx := range txs {

	Work1:
		for index,out := range tx.Jyh_TXOut {

			if out.Jyh_IsLockedPubkeyTxOut(address) {
				fmt.Println("看看是否是俊诚...")
				fmt.Println(address)

				fmt.Println(Jyh_spentTXOutputs)

				if len(Jyh_spentTXOutputs) == 0 {
					utxo := &Jyh_UTXO{tx.Jyh_TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range Jyh_spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.Jyh_TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &Jyh_UTXO{tx.Jyh_TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &Jyh_UTXO{tx.Jyh_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}


	blockIterator := blockchain.Jyh_Iterator()

	for {

		block := blockIterator.Jyh_Next()

		fmt.Println(block)
		fmt.Println()

		for i := len(block.Jyh_Transaction) - 1; i >= 0 ; i-- {

			tx := block.Jyh_Transaction[i]
			// txHash
			// Vins
			if tx.Jyh_IsCoinbaseTransaction() == false {
				for _, in := range tx.Jyh_TXIn {
					//是否能够解锁
					publicKeyHash := Jyh_Base58Decode([]byte(address))

					ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]

					if in.Jyh_IsLockedPubkeyTxIn(ripemd160Hash) {

						key := hex.EncodeToString(in.Jyh_TxHash)

						Jyh_spentTXOutputs[key] = append(Jyh_spentTXOutputs[key], in.Jyh_Vout)
					}

				}
			}

			// Vouts

		work:
			for index, out := range tx.Jyh_TXOut {

				if out.Jyh_IsLockedPubkeyTxOut(address) {

					fmt.Println(out)
					fmt.Println(Jyh_spentTXOutputs)

					//&{2 zhangqiang}
					//map[]

					if Jyh_spentTXOutputs != nil {

						//map[cea12d33b2e7083221bf3401764fb661fd6c34fab50f5460e77628c42ca0e92b:[0]]

						if len(Jyh_spentTXOutputs) != 0 {

							var isSpentUTXO bool

							for txHash, indexArray := range Jyh_spentTXOutputs {

								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.Jyh_TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO == false {

								utxo := &Jyh_UTXO{tx.Jyh_TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)

							}
						} else {
							utxo := &Jyh_UTXO{tx.Jyh_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}

		}

		fmt.Println(Jyh_spentTXOutputs)

		var hashInt big.Int
		hashInt.SetBytes(block.Jyh_PrevBlockHash)

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
const dbName  = "blockchain_%s.db"

// 表的名字
const blockTableName  = "blocks"

func BlockchainObject(nodeID string) *Jyh_BlockChain {
	dbName := fmt.Sprintf(dbName,nodeID)

	db,err:= bolt.Open(dbName, 0600,nil)

	if err!=nil{
		log.Fatal(err)
	}

	// 判断数据库是否存在
	if Jyh_DBExists(dbName) == false {
		fmt.Println("数据库不存在....")
		os.Exit(1)
	}

	var Tip []byte


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


	return &Jyh_BlockChain{Tip,db}
}
//1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM
//13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ
//1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG
//./main createBlockchain "13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"
func Jyh_CreateBlockchainWithGenesisBlock(address string,nodeID string) * Jyh_BlockChain{
	// 判断数据库是否存在
	dbName := fmt.Sprintf(dbName,nodeID)

	if Jyh_DBExists(dbName) {
		fmt.Println("创世区块已经存在.......")
		os.Exit(1)
	}

	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var Jyh_genesisHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		//  获取表
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b!=nil{
			txCoinbase:= Jyh_NewCoinbaseTransaction(address)

			genesisBlock:=Jyh_CreateGenesisBlock([]*Jyh_Transaction{txCoinbase})
			err:=b.Put(genesisBlock.Jyh_BlockHash,genesisBlock.Jyh_Serialize())
			if err!=nil{
				log.Panic(err)
			}

			err = b.Put([]byte("l"),genesisBlock.Jyh_BlockHash)
			if err!=nil{
				log.Panic(err)
			}
			Jyh_genesisHash = genesisBlock.Jyh_BlockHash
		}
		return nil
	})
	return &Jyh_BlockChain{Jyh_genesisHash,db}

}

// 判断数据库是否存在
func Jyh_DBExists(dbName string) bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}

	return true
}

//get UTXO from total tables
func (blc *Jyh_BlockChain)Jyh_FindUTXOMap() map[string]*Jyh_TXOutputs{
	blcIterator := blc.Jyh_Iterator()

	// 存储已花费的UTXO的信息
	spentableUTXOsMap := make(map[string][]*Jyh_TXInput)


	utxoMaps := make(map[string]*Jyh_TXOutputs)


	for {
		block := blcIterator.Jyh_Next()

		for i := len(block.Jyh_Transaction) - 1; i >= 0 ;i-- {

			txOutputs := &Jyh_TXOutputs{[]*Jyh_UTXO{}}

			tx := block.Jyh_Transaction[i]

			// coinbase
			if tx.Jyh_IsCoinbaseTransaction() == false {
				for _,txInput := range tx.Jyh_TXIn {

					txHash := hex.EncodeToString(txInput.Jyh_TxHash)
					spentableUTXOsMap[txHash] = append(spentableUTXOsMap[txHash],txInput)

				}
			}

			txHash := hex.EncodeToString(tx.Jyh_TxHash)

			txInputs := spentableUTXOsMap[txHash]

			if len(txInputs) > 0 {


			WorkOutLoop:
				for index,out := range tx.Jyh_TXOut  {

					for _,in := range  txInputs {

						outPublicKey := out.Jyh_Ripemd160Hash
						inPublicKey := in.Jyh_PubKey


						if bytes.Compare(outPublicKey,Jyh_HashPubKey(inPublicKey)) == 0 {
							if index == in.Jyh_Vout {

								continue WorkOutLoop
							} else {

								utxo := &Jyh_UTXO{tx.Jyh_TxHash,index,out}
								txOutputs.Jyh_UTXOS = append(txOutputs.Jyh_UTXOS,utxo)
							}
						}
					}


				}

			} else {

				for index,out := range tx.Jyh_TXOut {
					utxo := &Jyh_UTXO{tx.Jyh_TxHash,index,out}
					txOutputs.Jyh_UTXOS = append(txOutputs.Jyh_UTXOS,utxo)
				}
			}


			// 设置键值对
			utxoMaps[txHash] = txOutputs

		}


		// 找到创世区块时退出
		var hashInt big.Int
		hashInt.SetBytes(block.Jyh_PrevBlockHash)

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
func (blockchain *Jyh_BlockChain) Jyh_FindSpendableUTXOS(from string, amount int,txs []*Jyh_Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO

	utxos := blockchain.Jyh_UnUTXOs(from,txs)

	spendableUTXO := make(map[string][]int)

	//2. 遍历utxos

	var value int64

	for _, utxo := range utxos {

		value = value + utxo.Jyh_Output.Jyh_Value

		hash := hex.EncodeToString(utxo.Jyh_TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Jyh_Index)

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
func (blockchain *Jyh_BlockChain) Jyh_MineNewBlock(from []string, to []string, amount []string,nodeID string) {
	//	$ ./bc send -from '["juncheng"]' -to '["zhangqiang"]' -amount '["2"]'
	//	[juncheng]
	//	[zhangqiang]
	//	[2]

	//1.建立一笔交易


	utxoSet := &Jyh_UTXOSet{blockchain}

	var txs []*Jyh_Transaction

	for index,address := range from {
		value, _ := strconv.Atoi(amount[index])
		tx := Jyh_NewSimpleTransaction(address, to[index], int64(value), utxoSet,txs,nodeID)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	//奖励
	tx := Jyh_NewCoinbaseTransaction(from[0])
	txs = append(txs,tx)


	//1. 通过相关算法建立Transaction数组
	var block *Jyh_Block

	blockchain.Jyh_BlockDB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			block = Jyh_DeserializeBlock(blockBytes)

		}

		return nil
	})
	//	./main send -from '["183UkSrFV72YZ53FS5MAxiSPgNrgWv9uU8"]' -to '["1Ap11ovDK7aorb1Cmkb3L6JxhLZDCC2sbN"]' -amount '["3"]'

	// 在建立新区块之前对txs进行签名验证

	_txs := []*Jyh_Transaction{}

	for _,tx := range txs  {

		if blockchain.Jyh_VerifyTransaction(tx,_txs) != true {
			log.Panic("ERROR: Invalid transaction")
		}

		_txs = append(_txs,tx)
	}


	//2. 建立新的区块
	block = Jyh_NewBlock(txs, block.Jyh_Height+1, block.Jyh_BlockHash)

	//将新区块存储到数据库
	blockchain.Jyh_BlockDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {

			b.Put(block.Jyh_BlockHash, block.Jyh_Serialize())

			b.Put([]byte("l"), block.Jyh_BlockHash)

			blockchain.Jyh_Tip = block.Jyh_BlockHash

		}
		return nil
	})

}

func (bc *Jyh_BlockChain) Jyh_GetBestHeight() int64{
	block := bc.Jyh_Iterator().Jyh_Next()
	return block.Jyh_Height
}

func (bc *Jyh_BlockChain) Jyh_GetBlockHashes() [][]byte{
	blockIterator := bc.Jyh_Iterator()

	var blockHashs [][]byte

	for{
		block := blockIterator.Jyh_Next()

		blockHashs=append(blockHashs,block.Jyh_BlockHash)

		var hashInt big.Int
		hashInt.SetBytes(block.Jyh_PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0))==0{
			break;
		}
	}
	return blockHashs
}

func (bc *Jyh_BlockChain) Jyh_GetBlock(blockHash []byte) ([]byte,error){
	var blockBytes []byte

	err :=bc.Jyh_BlockDB.View(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(blockTableName))

		if b!=nil{
			blockBytes = b.Get(blockHash)
		}
		return nil
	})
	return blockBytes,err
}

func (bc *Jyh_BlockChain) Jyh_AddBlock(block *Jyh_Block){
	err := bc.Jyh_BlockDB.Update(func(tx *bolt.Tx) error{
		b := tx.Bucket([]byte(blockTableName))

		if b!=nil{
			blockExist := b.Get(block.Jyh_BlockHash)

			if blockExist!=nil{
				return nil
			}

			err := b.Put(block.Jyh_BlockHash,block.Jyh_Serialize())

			if err!=nil{
				log.Panic(err)
			}

			blockHash := b.Get([]byte("l"))

			blockBytes := b.Get(blockHash)

			blockInDB := Jyh_DeserializeBlock(blockBytes)

			if blockInDB.Jyh_Height < block.Jyh_Height{
				 b.Put([]byte("l"),block.Jyh_BlockHash)
				 bc.Jyh_Tip=block.Jyh_BlockHash
			}
		}
		return nil
	})
	if err != nil{
		log.Panic(err)
	}
}