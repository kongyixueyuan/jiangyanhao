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

func FindSpendableUTXO(blc *BlockChain) map[string][]int{


	spendableUTXO := make(map[string][]int)

	var blockchainIterator *BlockchainIterator
	blockchainIterator = blc.Iterator()
	for {
		block := blockchainIterator.Next()
		//[Txhash]{0,1,2}
		for _, tx := range block.Transaction {
			//fmt.Println("it's typing Touts:")
			for i, _ := range tx.TXOut {
				spendableUTXO[string(tx.TxHash)]=append(spendableUTXO[string(tx.TxHash)], i)
				//fmt.Println(i)
			}
			//fmt.Println(tx.TxHash)
			//fmt.Println(spendableUTXO[string(tx.TxHash)]).flag≥encodingflag
		}

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
	/*
	for tx, nums := range spendableUTXO {
		fmt.Printf(":txout:%x--------",tx)
		for i := range nums{
			fmt.Printf("%d---",nums[i])
		}
		fmt.Printf("\n")
	}*/
	//fmt.Println(Blockchain)
	blockchainIterator = blc.Iterator()
	var ff []int
	for {

		block := blockchainIterator.Next()
		////input:[Txhash][1]------>[Txhash][0,2]
		for _, tx := range block.Transaction {
			//fmt.Println("it's typing in:")
			for _, in := range tx.TXIn {
				//	fmt.Printf("for in------------,%x",in.TxHash)
				//	fmt.Printf("-------in Vout:%d-----\n",in.Vout)

				for _ , value := range spendableUTXO[string(in.TxHash)]{
					//		fmt.Printf("it's in.Txhash-----%x\n", in.TxHash)
					//fmt.Printf("it's tx.Txhash-----%x", tx.TxHash)
					//		fmt.Printf("------vout: %d comparing:%d\n",in.Vout,value)

					//			fmt.Println("it's in tx ======== in ")
					ff=append(ff[:0],ff[:0]...)
					spendableUTXO[string(in.TxHash)]=ff

					if(in.Vout!=value){
						spendableUTXO[string(in.TxHash)]= append(spendableUTXO[string(in.TxHash)], value)
						//fmt.Printf("it's going to add in spendableUTXO-----%x",tx.TxHash)
						//fmt.Println("in Vout:",in.Vout,"it's value:",value)
					}

				}
			}
			//fmt.Println(tx.TxHash)
			//fmt.Println(spendableUTXO[string(tx.TxHash)])
		}
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
	/*fmt.Println("p--------------")
	for tx, nums := range spendableUTXO {
		fmt.Printf(":txout:%x----",tx)
		for i := range nums{
			fmt.Printf("%d---",nums[i])
		}
		fmt.Printf("\n")
	}*/

	fmt.Println("------------------------------")

	return spendableUTXO
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