package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type Transaction struct{
	TxHash []byte

	TXIn []*TXInput
	TXOut []*TXOutput
}

func (tx *Transaction) HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())

	tx.TxHash = hash[:]
}

func FindSpendableUTXO(blc *BlockChain) []*TXOutput{

	var txOutputs []*TXOutput

	blockchainIterator := blc.Iterator()
	totalUTXO := make(map[string][]int)
	spendableUTXO := make(map[string][]int)


	for {
		block := blockchainIterator.Next()

		for _, tx := range block.Transaction {
			fmt.Println("it's typing Touts:")
			for i, _ := range tx.TXOut {
				totalUTXO[string(tx.TxHash)]=append(totalUTXO[string(tx.TxHash)], i)
				fmt.Println(i)
			}
			//fmt.Println(tx.TxHash)
			//fmt.Println(spendableUTXO[string(tx.TxHash)]).flag≥encodingflag
		}

		for _, tx := range block.Transaction {
			for _, in := range tx.TXIn {
				for _ , value := range totalUTXO[string(tx.TxHash)]{
					if(in.Vout!=value){
						spendableUTXO[string(tx.TxHash)]= append(spendableUTXO[string(tx.TxHash)], value)
					}
				}
			}
			//fmt.Println(tx.TxHash)
			//fmt.Println(spendableUTXO[string(tx.TxHash)])
		}


		fmt.Println("------------------------------")

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		// Cmp compares x and y and returns:
		//
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y

		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}

	}
	//fmt.Println(Blockchain)



	return txOutputs
}


func SetCoinbBaseTransaction(address string) []*Transaction{
	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	//input: 0000,-1,address
	CoinbaseIn := &TXInput{[]byte{},-1,address}
	txIntputs = append(txIntputs,CoinbaseIn)
	fmt.Println("tx hash:%x, value:%d, address:%s",CoinbaseIn.TxHash,CoinbaseIn.Vout,CoinbaseIn.ScriptSig)

	//output: 15,address
	CoinbaseOut := &TXOutput{15,address}
	txOutputs = append(txOutputs,CoinbaseOut)
	fmt.Println("tx value:%d, address:%s",CoinbaseOut.Value,CoinbaseOut.ScriptPubKey)

	GenesisTransaction := &Transaction{[]byte{},txIntputs,txOutputs}
	GenesisTransaction.HashTransaction()


	//设置hash值
	//GenesisTransaction.HashTransaction()

	return []*Transaction{GenesisTransaction}
}