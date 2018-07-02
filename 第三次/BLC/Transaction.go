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

func FindUTXOOutput(blc *BlockChain) []*TXOutput{

	var txOutputs []*TXOutput

	blockchainIterator := blc.Iterator()

	spendableUTXO := make(map[string][]int)

	d := []int{};

	for {
		block := blockchainIterator.Next()

		for _, tx := range block.Transaction {
			fmt.Println("it's typing Touts:")
			for i, _ := range tx.TXOut {
				d=append(d, i)
				fmt.Println(i)
			}
			spendableUTXO[string(tx.TxHash)]=d
			d = []int{}
			//fmt.Println(tx.TxHash)
			//fmt.Println(spendableUTXO[string(tx.TxHash)]).flag≥encoding≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥≥/*/*/*/*/*/*/*/*/*/*flag*/*/*/*/*/*/*/*/*/*/
		}

		for _, tx := range block.Transaction {
			for i, in := range tx.TXIn {

				spendableUTXO[string(tx.TxHash)]=append(spendableUTXO[string(tx.TxHash)][:i], spendableUTXO[string(tx.TxHash)][i+1:]...)
				//重组

				fmt.Printf("%x\n", in.TxHash)
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}
			spendableUTXO[string(tx.TxHash)]=d
			d = []int{}
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