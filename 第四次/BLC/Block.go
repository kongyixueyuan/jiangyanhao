package BLC

import(
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"github.com/boltdb/bolt"
	"crypto/sha256"
	"strings"
	"strconv"
	"math/big"
)
type Block struct{
	Height int64

	PrevBlockHash []byte

	Nonce int64

	Transaction []*Transaction

	Timestamp int64

	BlockHash []byte
}
/*
func PrintBlock(block *Block){
	fmt.Println("===============")
	fmt.Println(block.Height)
	fmt.Println(hex.EncodeToString(block.PrevBlockHash))
	fmt.Println(block.nonce)
	fmt.Println(block.Timestamp)
	fmt.Println(hex.EncodeToString(block.BlockHash))
	fmt.Println("===============")
}*/

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
	//fmt.Println("in serialize, %d",block.Nonce)
	if err != nil {
		log.Panic(err)
	}

	//fmt.Println("/////////in serialize, %d",DeserializeBlock(result.Bytes()))
	return result.Bytes()
}

// 反序列化
func DeserializeBlock(blockBytes []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	//fmt.Println("---------in deserialize, %d",block.Nonce)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewBlock(Transaction []*Transaction, blockchain *BlockChain){


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
			newBlock :=&Block{block.Height+1, block.BlockHash,0,Transaction,time.Now().Unix(),nil}

			Pow:=ProofOfWork(newBlock)
			nonce, hash := Pow.run()
			newBlock.Nonce=nonce
			newBlock.BlockHash=hash[:]
			//fmt.Println("it's hash+++++++++++,%s,%s",hex.EncodeToString(hash),hex.EncodeToString(newBlock.BlockHash))

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

// 制作transaction
func MineNewBlock(from []string, to []string, amount []string, blockchain *BlockChain) {

	blockchainIterator := blockchain.Iterator()


	//把可花费的utxo算出来，都是output
	//[Txhash][0,2]
	spendableUTXO := FindSpendableUTXO(blockchain)
	//var spendableTxout []*TXOutput
	spendableMap := make(map[string][]*TXOutput)
	value := int64(0)

	var newTransaction Transaction
	var TempTXIn []*TXInput
	var TempTXOut []*TXOutput
	//看看from需要用哪几个output，算余额够不够
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
/*
	for txhash, Txs := range spendableMap {
		for j, txOut := range Txs {
			fmt.Println(txhash)
			fmt.Println(txOut)
			fmt.Println("j:%d",j)
			fmt.Println("-----")
		}
	}*/

	//2、看看与from的人是否一致
	/*	for _, txOut :=range spendableTxout{
			for _ , fromPeople := range from{
				if(strings.Compare(txOut.ScriptPubKey, fromPeople)==0){
				value=value+txOut.Value
				}
			}
		}*/
	for i, fromPeople := range from {
		var dont []*TXInput
		TempTXIn=dont
		amount, err := strconv.Atoi(amount[i])
		if (err != nil) {
			log.Fatal("no amount! in %d", i)
		}

		//翻spendableMap
		for txhash, Txs := range spendableMap {
			for j, txOut := range Txs {
				if (strings.Compare(txOut.ScriptPubKey, fromPeople) == 0) {
					value = value + txOut.Value
					//txIn先搭着，余额不够就清空
					TempTXIn = append(TempTXIn, &TXInput{[]byte(txhash), j, fromPeople})
					//如果余额够的话，搭建txin与txout

					if (value >= int64(amount)) {
						TempTXOut = append(TempTXOut, &TXOutput{value-int64(amount), from[i]})
						TempTXOut = append(TempTXOut, &TXOutput{int64(amount), to[i]})
					}
				}
			}
		}
		//如果余额不够，结束
		if (value < int64(amount)) {
			log.Fatal("not enough money! go to work! %s", fromPeople)

		}
		newTransaction = Transaction{[]byte{},TempTXIn,TempTXOut}
		newTransaction.HashTransaction()
	}


	NewBlock([]*Transaction{&newTransaction},blockchain)

}
func NewGenesisBlock(address string) * Block{
	GenesisTransaction:=SetCoinbBaseTransaction(address)
	block := &Block{1,IntToHex(0),0,GenesisTransaction,time.Now().Unix(),nil}

	Pow:=ProofOfWork(block)
	nonce, hash := Pow.run()

	block.BlockHash=hash[:]
	block.Nonce=nonce

	//PrintBlock(block)
	return block
}