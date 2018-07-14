package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
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