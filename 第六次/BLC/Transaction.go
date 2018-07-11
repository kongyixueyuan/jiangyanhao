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
	"time"
)

type Transaction struct{
	TxHash []byte

	TXIn []*TXInput
	TXOut []*TXOutput
}

func (tx *Transaction) jyh_HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	resultBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()),result.Bytes()},[]byte{})

	hash := sha256.Sum256(resultBytes)

	tx.TxHash = hash[:]
}

// 判断当前的交易是否是Coinbase交易
func (tx *Transaction) jyh_IsCoinbaseTransaction() bool {
	return len(tx.TXIn[0].TxHash) == 0 && tx.TXIn[0].Vout == -1
}

func jyh_FindSpendableUTXO(blc *BlockChain) map[string][]int{


	spendableUTXO := make(map[string][]int)
//	totalUTXO := make(map[string][]int)

	var blockchainIterator *BlockchainIterator
	blockchainIterator = blc.Iterator()
	for {
		block := blockchainIterator.jyh_Next()
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
		block := blockchainIterator.jyh_Next()
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


func jyh_SetCoinbBaseTransaction(address string) *Transaction{
	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	//input: 0000,-1,address
	CoinbaseIn := &TXInput{[]byte{},-1,nil , []byte{}}
	txIntputs = append(txIntputs,CoinbaseIn)


	//output: 15,address
	CoinbaseOut := &TXOutput{15,ConvertAddtoRip(address)}
	txOutputs = append(txOutputs,CoinbaseOut)
	fmt.Println("tx value:%d, address:%s",CoinbaseOut.Value,address)


	CoinbaseTransaction := &Transaction{[]byte{},txIntputs,txOutputs}
	CoinbaseTransaction.jyh_HashTransaction()


	//设置hash值
	//GenesisTransaction.HashTransaction()

	return CoinbaseTransaction
}

func (tx *Transaction) jyh_Hash() []byte {

	txCopy := tx

	txCopy.TxHash = []byte{}

	hash := sha256.Sum256(txCopy.jyh_Serialize())
	return hash[:]
}

func (tx *Transaction) jyh_Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) jyh_Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {

	if tx.jyh_IsCoinbaseTransaction() {
		return
	}


	for _, vin := range tx.TXIn {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}


	txCopy := tx.jyh_TrimmedCopy()

	for inID, vin := range txCopy.TXIn {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.TXIn[inID].Signature = nil
		txCopy.TXIn[inID].PubKey = prevTx.TXOut[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.jyh_Hash()
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
func (tx *Transaction) jyh_TrimmedCopy() Transaction {
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

func (tx *Transaction) jyh_Verify(prevTXs map[string]Transaction) bool {
	if tx.jyh_IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.TXIn {
		if prevTXs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.jyh_TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.TXIn {
		prevTx := prevTXs[hex.EncodeToString(vin.TxHash)]
		txCopy.TXIn[inID].Signature = nil
		txCopy.TXIn[inID].PubKey = prevTx.TXOut[vin.Vout].Ripemd160Hash
		txCopy.TxHash = txCopy.jyh_Hash()
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


func jyh_NewSimpleTransaction(from string,to string,amount int64,utxoSet *jyh_UTXOSet,txs []*Transaction) *Transaction{

	wallets,_ := NewWallets()
	wallet := wallets.WalletsMap[from]

	// 通过一个函数，返回
	money,spendableUTXODic := utxoSet.jyh_FindSpendableUTXOS(from,amount,txs)

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {

			txInput := &TXInput{txHashBytes,index,nil,wallet.PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := NewTXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = NewTXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.jyh_HashTransaction()

	//进行签名
	utxoSet.Blockchain.SignTransaction(tx, wallet.PrivateKey,txs)

	return tx

}

/*

func NewSimpleTransaction(from string,to string,amount int64,utxoSet *UTXOSet,txs []*Transaction) *Transaction{

	wallets,_ := NewWallets()
	wallet := wallets.WalletsMap[from]

	// 通过一个函数，返回
	money,spendableUTXODic := utxoSet.FindSpendableUTXOS(from,amount,txs)

	blockchainIterator := blockChain.Iterator()
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
	//		if len(spendableUTXO[string(tx.TxHash)])==0{
				continue
			}else{
				for _, value := range spendableUTXO[string(tx.TxHash)] {
					//	spendableTxout=append(spendableTxout, txOut)
					spendableMap[string(tx.TxHash)] = append(spendableMap[string(tx.TxHash)], tx.TXOut[value])
				}
			}
		}
		var hashInt bigInt
		//hashInt.SetBytes(block.PrevBlockHash)

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



	for i, fromPeople := range from {
		var dont []*TXInput
		TempTXIn=dont
		amount, err := strconv.Atoi(amount[i])
		if (err != nil) {
			log.Fatal("no amount! in %d", i)
		}

		//翻spendableMap
		for txhash, TxOutputs := range spendableMap {
			for j, txOut := range TxOutputs {
				//	if (strings.Compare(txOut.PubkeyHash, fromPeople) == 0) {
				if wallet.IsValidAddress([]byte(fromPeople))==false{
					log.Fatal("Send is not bitcoin address!", fromPeople)
					os.Exit(0)
				}
				if (txOut.IsLockedPubkeyTxOut(fromPeople)) {
					wallet = wallets.WalletsMap[fromPeople]
					if wallet.IsValidAddress([]byte(to[i]))==false{
						log.Fatal("To is not bitcoin address!", fromPeople)
						os.Exit(0)
					}
					value = value + txOut.Value
					//txIn先搭着，余额不够就清空, j
					TempTXIn = append(TempTXIn, &TXInput{[]byte(txhash), spendableUTXO[string(txhash)][j], nil, wallet.PublicKey})
					//如果余额够的话，搭建txin与txout

					if (value >= int64(amount)) {
						TempTXOut = append(TempTXOut, &TXOutput{value-int64(amount), ConvertAddtoRip(from[i])})
						TempTXOut = append(TempTXOut, &TXOutput{int64(amount), ConvertAddtoRip(to[i])})
						break
					}
				}
			}
			if (value >= int64(amount)) {
				break
			}
		}
		//如果余额不够，结束
		if (value < int64(amount)) {
			log.Fatal("not enough money! go to work! %s", fromPeople)

		}
		value = 0
		newTransaction = Transaction{[]byte{},TempTXIn,TempTXOut}
		newTransaction.HashTransaction()

		//进行签名
		blockchain.SignTransaction(&newTransaction, wallet.PrivateKey)
	}


	MineNewBlock([]*Transaction{&newTransaction},blockchain)

}*/