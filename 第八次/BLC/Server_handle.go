package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
	"encoding/hex"
	"os"
	"github.com/boltdb/bolt"
)

func Jyh_handleVersion(request []byte, bc *Jyh_BlockChain){
	var buff bytes.Buffer
	var payload Jyh_Version

	dataBytes := request[COMMANDLENGTH:]

	buff.Write(dataBytes)
	dec:=gob.NewDecoder(&buff)
	err:=dec.Decode(&payload)

	if err!=nil{
		log.Panic(err)
	}

	bestHeight := bc.Jyh_GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if bestHeight>foreignerBestHeight{
		Jyh_sendVersion(payload.AddrFreom, bc)
	}else if bestHeight<foreignerBestHeight{
		Jyh_sendGetBlocks(payload.AddrFreom)
	}
}

func Jyh_handleAddr (request []byte, bc *Jyh_BlockChain){

}

func Jyh_handleGetData(request []byte,bc *Jyh_BlockChain)  {

	var buff bytes.Buffer
	var payload Jyh_GetData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == BLOCK_TYPE {

		block, err := bc.Jyh_GetBlock([]byte(payload.Hash))
		if err != nil {
			return
		}

		Jyh_sendBlock(payload.AddrFrom, block)
	}

	if payload.Type == "tx" {

	}
}

func Jyh_handleGetblocks (request []byte, bc *Jyh_BlockChain){
	var buff bytes.Buffer
	var payload Jyh_GetBlocks

	dataBytes := request[COMMANDLENGTH:]

	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err!=nil{
		log.Panic(err)
	}

	blocks :=bc.Jyh_GetBlockHashes()

	//
	Jyh_sendInv(payload.AddrFrom, BLOCK_TYPE, blocks)
}



func Jyh_handleBlock(request []byte,bc *Jyh_BlockChain)  {
	var buff bytes.Buffer
	var payload Jyh_BlockData

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockBytes := payload.Block

	block := Jyh_DeserializeBlock(blockBytes)

	fmt.Println("Recevied a new block!")
	bc.Jyh_AddBlock(block)

	fmt.Printf("Added block %x\n", block.Jyh_BlockHash)

	if len(transactionArray) > 0 {
		blockHash := transactionArray[0]
		Jyh_sendGetData(payload.AddrFrom, "block", blockHash)

		transactionArray = transactionArray[1:]
	} else {

		fmt.Println("数据库重置......")
		UTXOSet := &Jyh_UTXOSet{bc}
		UTXOSet.Jyh_ResetUTXOSet()

	}

}

func Jyh_handleTx(request []byte,bc *Jyh_BlockChain)  {

	var buff bytes.Buffer
	var payload Jyh_Tx

	dataBytes := request[COMMANDLENGTH:]

	// 反序列化
	buff.Write(dataBytes)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic("发序列化错误:",err)
	}

	//-----

	tx := payload.Tx
	memoryTxPool[hex.EncodeToString(tx.Jyh_TxHash)] = tx

	// 说明主节点自己
	if nodeAddress == knowNodes[0] {
		// 给矿工节点发送交易hash
		for _,nodeAddr := range knowNodes {

			if nodeAddr != nodeAddress && nodeAddr != payload.AddrFrom {
				Jyh_sendInv(nodeAddr,TX_TYPE,[][]byte{tx.Jyh_TxHash})
			}

		}
	}

	// 矿工进行挖矿验证
	// "" | 1DVFvyCK8qTQkLBTZ5fkh5eDSbcZVoHAsj
	if len(minerAddress) > 0 {

		bc.Jyh_BlockDB.Close()

		blockchain := BlockchainObject(os.Getenv("NODE_ID"))
		defer blockchain.Jyh_BlockDB.Close()

		//1.建立一笔交易
		//
		utxoSet := &Jyh_UTXOSet{blockchain}

		var txs []*Jyh_Transaction

		txs = append(txs, tx)

		//奖励
		coinTX := Jyh_NewCoinbaseTransaction(minerAddress)
		txs = append(txs, coinTX)

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

		// 在建立新区块之前对txs进行签名验证

		_txs := []*Jyh_Transaction{}

		for _, tx := range txs {

			if blockchain.Jyh_VerifyTransaction(tx, _txs) != true {
				log.Panic("ERROR: Invalid transaction")
			}

			_txs = append(_txs, tx)
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
		//转账成功以后，需要更新一下
		utxoSet.Jyh_Update()
		Jyh_sendBlock(knowNodes[0], block.Jyh_Serialize())






		//
		//txs := []*Transaction{tx}
		//
		//coinbaseTx := NewCoinbaseTransaction(minerAddress)
		//txs = append(txs,coinbaseTx)
		//
		//_txs := []*Transaction{}
		//
		//fmt.Println("开始进行数字签名验证.....")
		//
		//for index,tx := range txs  {
		//
		//	fmt.Printf("开始第%d次验证...\n",index)
		//
		//	if bc.VerifyTransaction(tx,_txs) != true {
		//		log.Panic("ERROR: Invalid transaction")
		//	}
		//
		//	fmt.Printf("第%d次验证成功\n",index)
		//	_txs = append(_txs,tx)
		//}
		//
		//fmt.Println("数字签名验证成功.....")
		//
		////1. 通过相关算法建立Transaction数组
		//var block *Block
		//
		//bc.DB.View(func(tx *bolt.Tx) error {
		//
		//	b := tx.Bucket([]byte(blockTableName))
		//	if b != nil {
		//
		//		hash := b.Get([]byte("l"))
		//
		//		blockBytes := b.Get(hash)
		//
		//		block = DeserializeBlock(blockBytes)
		//
		//	}
		//
		//	return nil
		//})
		//
		////2. 建立新的区块
		//block = NewBlock(txs, block.Height+1, block.Hash)
		//
		////将新区块存储到数据库
		//bc.DB.Update(func(tx *bolt.Tx) error {
		//	b := tx.Bucket([]byte(blockTableName))
		//	if b != nil {
		//
		//		b.Put(block.Hash, block.Serialize())
		//
		//		b.Put([]byte("l"), block.Hash)
		//
		//		bc.Tip = block.Hash
		//
		//	}
		//	return nil
		//})
		//
		//sendBlock(knowNodes[0],block.Serialize())
	}


}



func Jyh_handleInv(request []byte, bc *Jyh_BlockChain){
	var buff bytes.Buffer
	var payload Jyh_Inv

	dataBytes:=request[COMMANDLENGTH:]

	buff.Write(dataBytes)
	dec:=gob.NewDecoder(&buff)
	err:=dec.Decode(&payload)
	if err!=nil{
		log.Panic(err)
	}

	if payload.Type == BLOCK_TYPE{
		blockHash := payload.Items[0]
		Jyh_sendGetData(payload.AddrFrom, BLOCK_TYPE, blockHash)
	}
	if len(payload.Items) >= 1 {
		transactionArray = payload.Items[1:]
	}

	if payload.Type == TX_TYPE {

	}



}