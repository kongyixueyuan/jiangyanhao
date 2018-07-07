package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"math/big"
	"crypto/ecdsa"
	"encoding/hex"
	"crypto/rand"
	"crypto/elliptic"
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

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {

	return len(tx.TXIn[0].TxHash) == 0 && tx.TXIn[0].Vout == -1
}

func FindSpendableUTXO(blc *BlockChain) map[string][]int{


	spendableUTXO := make(map[string][]int)
//	totalUTXO := make(map[string][]int)

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

		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)


		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break;
		}

	}
	/*fmt.Println(":--------------------:%x----")
	for tx, nums := range spendableUTXO {
		fmt.Printf(":txout:%x----", tx)
		for i := range nums {
			fmt.Printf("%d---", nums[i])
		}
		fmt.Printf("\n")
	}
	fmt.Println(":--------------------:%x----")*/
	blockchainIterator = blc.Iterator()
	//var temp []int
	for {
		block := blockchainIterator.Next()
		////input:[Txhash][1]------>[Txhash][0,2]
		for _, tx := range block.Transaction {
			//fmt.Println("it's typing in:")
			for _, in := range tx.TXIn {
				for _ , value := range spendableUTXO[string(in.TxHash)]{
//					temp=append(temp[:0],temp[:0]...)
//					spendableUTXO[string(in.TxHash)]=temp
					if(in.Vout==value){
						spendableUTXO[string(in.TxHash)]=RemoveInt(spendableUTXO[string(in.TxHash)], value)
					}

					/*if(in.Vout!=value){
						spendableUTXO[string(in.TxHash)]= append(spendableUTXO[string(in.TxHash)], value)
						fmt.Printf("it's going to add in spendableUTXO-----%x",tx.TxHash)
						fmt.Println("in Vout:",in.Vout,"it's value:",value)
					}*/

				}
				/*if(len(spendableUTXO[string(in.TxHash)])!=0){
					totalUTXO[string(in.TxHash)]=spendableUTXO[string(in.TxHash)]
				}*/
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
	/*fmt.Println(":-/////////////////")
	for tx, nums := range totalUTXO {
		fmt.Printf("totalxxxxxxxxxxx:%x----", tx)
		for i := range nums {
			fmt.Printf("%d---", nums[i])
		}
		fmt.Println("nums lenlen", len(nums))

	}
	for tx, nums := range spendableUTXO {
		fmt.Printf(":txout:%x----", tx)
		for i := range nums {
			fmt.Printf("%d---", nums[i])
		}
		fmt.Println("nums lenlen", len(nums))

	}
	fmt.Println("spendable lenlen", len(spendableUTXO))
	*/
	return spendableUTXO
}


func SetCoinbBaseTransaction(address string) []*Transaction{
	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	//input: 0000,-1,address
	CoinbaseIn := &TXInput{[]byte{},-1,nil , nil}
	txIntputs = append(txIntputs,CoinbaseIn)


	//output: 15,address
	CoinbaseOut := &TXOutput{15,ConvertAddtoRip(address)}
	txOutputs = append(txOutputs,CoinbaseOut)
	fmt.Println("tx value:%d, address:%s",CoinbaseOut.Value,address)

	GenesisTransaction := &Transaction{[]byte{},txIntputs,txOutputs}
	GenesisTransaction.HashTransaction()


	//设置hash值
	//GenesisTransaction.HashTransaction()

	return []*Transaction{GenesisTransaction}
}

func (tx *Transaction) Hash() []byte {

	txCopy := tx

	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.IsCoinbaseTransaction() {
		return
	}


	for _, vin := range tx.TXIn {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}


	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.TXIn {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.TXIn[inID].Signature = nil
		txCopy.TXIn[inID].PubKey = prevTx.TXOut[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.TXIn[inID].PubKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.TXIn[inID].Signature = signature
	}
}


// 拷贝一份新的Transaction用于签名                                    T
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []*TXInput
	var outputs []*TXOutput

	for _, vin := range tx.TXIn {
		inputs = append(inputs, &TXInput{vin.TxHash, vin.Vout, nil, nil})
	}

	for _, vout := range tx.TXOut {
		outputs = append(outputs, &TXOutput{vout.Value, vout.Ripemd160Hash})
	}

	txCopy := Transaction{tx.TxHash, inputs, outputs}

	return txCopy
}


// 数字签名验证

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.TXIn {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.TXIn {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.TXIn[inID].Signature = nil
		txCopy.TXIn[inID].PubKey = prevTx.TXOut[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.Hash()
		txCopy.TXIn[inID].PubKey = nil


		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}