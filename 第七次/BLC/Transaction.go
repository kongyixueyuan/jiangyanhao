package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"math/big"
	"crypto/ecdsa"
	"encoding/hex"
	"crypto/rand"
	"crypto/elliptic"
	"time"
)

type Jyh_Transaction struct{
	Jyh_TxHash []byte

	Jyh_TXIn []*Jyh_TXInput
	Jyh_TXOut []*Jyh_TXOutput
}

func (tx *Jyh_Transaction) Jyh_HashTransaction()  {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	resultBytes := bytes.Join([][]byte{IntToHex(time.Now().Unix()),result.Bytes()},[]byte{})

	hash := sha256.Sum256(resultBytes)

	tx.Jyh_TxHash = hash[:]
}

// 判断当前的交易是否是Coinbase交易
func (tx *Jyh_Transaction) Jyh_IsCoinbaseTransaction() bool {
	return len(tx.Jyh_TXIn[0].Jyh_TxHash) == 0 && tx.Jyh_TXIn[0].Jyh_Vout == -1
}

func Jyh_FindSpendableUTXO(blc *Jyh_BlockChain) map[string][]int{


	spendableUTXO := make(map[string][]int)
//	totalUTXO := make(map[string][]int)

	var blockchainIterator *Jyh_BlockchainIterator
	blockchainIterator = blc.Jyh_Iterator()
	for {
		block := blockchainIterator.Jyh_Next()
		//[Txhash]{0,1,2}
		for _, tx := range block.Jyh_Transaction {
			//fmt.Println("it's typing Touts:")
			for i, _ := range tx.Jyh_TXOut {
				spendableUTXO[string(tx.Jyh_TxHash)]=append(spendableUTXO[string(tx.Jyh_TxHash)], i)
				//fmt.Println(i)
			}

		}

		var hashInt big.Int
		hashInt.SetBytes(block.Jyh_PrevBlockHash)


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
	blockchainIterator = blc.Jyh_Iterator()
	//var temp []int
	for {
		block := blockchainIterator.Jyh_Next()
		////input:[Txhash][1]------>[Txhash][0,2]
		for _, tx := range block.Jyh_Transaction {
			//fmt.Println("it's typing in:")
			for _, in := range tx.Jyh_TXIn {
				for _ , value := range spendableUTXO[string(in.Jyh_TxHash)]{
//					temp=append(temp[:0],temp[:0]...)
//					spendableUTXO[string(in.TxHash)]=temp
					if(in.Jyh_Vout==value){
						spendableUTXO[string(in.Jyh_TxHash)]=RemoveInt(spendableUTXO[string(in.Jyh_TxHash)], value)
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
		hashInt.SetBytes(block.Jyh_PrevBlockHash)

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


func Jyh_NewCoinbaseTransaction(address string) *Jyh_Transaction{


	//input: 0000,-1,address
	CoinbaseIn := &Jyh_TXInput{[]byte{},-1,nil , []byte{}}


	//output: 15,address
	CoinbaseOut := NewTXOutput(15,address)


	CoinbaseTransaction := &Jyh_Transaction{[]byte{},[]*Jyh_TXInput{CoinbaseIn},[]*Jyh_TXOutput{CoinbaseOut}}
	CoinbaseTransaction.Jyh_HashTransaction()


	//设置hash值
	//GenesisTransaction.HashTransaction()

	return CoinbaseTransaction
}

func (tx *Jyh_Transaction) Jyh_Hash() []byte {

	txCopy := tx

	txCopy.Jyh_TxHash = []byte{}

	hash := sha256.Sum256(txCopy.Jyh_Serialize())
	return hash[:]
}

func (tx *Jyh_Transaction) Jyh_Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Jyh_Transaction) Jyh_Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Jyh_Transaction) {

	if tx.Jyh_IsCoinbaseTransaction() {
		return
	}


	for _, vin := range tx.Jyh_TXIn {
		if prevTXs[hex.EncodeToString(vin.Jyh_TxHash)].Jyh_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}


	txCopy := tx.Jyh_TrimmedCopy()

	for inID, vin := range txCopy.Jyh_TXIn {
		prevTx := prevTXs[hex.EncodeToString(vin.Jyh_TxHash)]
		txCopy.Jyh_TXIn[inID].Jyh_Signature = nil
		txCopy.Jyh_TXIn[inID].Jyh_PubKey = prevTx.Jyh_TXOut[vin.Jyh_Vout].Jyh_Ripemd160Hash
		txCopy.Jyh_TxHash = txCopy.Jyh_Hash()
		txCopy.Jyh_TXIn[inID].Jyh_PubKey = nil

		// 签名代码
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.Jyh_TxHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Jyh_TXIn[inID].Jyh_Signature = signature
	}
}


// 拷贝一份新的Transaction用于签名                                    T
func (tx *Jyh_Transaction) Jyh_TrimmedCopy() Jyh_Transaction {
	var inputs []*Jyh_TXInput
	var outputs []*Jyh_TXOutput

	for _, vin := range tx.Jyh_TXIn {
		inputs = append(inputs, &Jyh_TXInput{vin.Jyh_TxHash, vin.Jyh_Vout, nil, nil})
	}

	for _, vout := range tx.Jyh_TXOut {
		outputs = append(outputs, &Jyh_TXOutput{vout.Jyh_Value, vout.Jyh_Ripemd160Hash})
	}

	txCopy := Jyh_Transaction{tx.Jyh_TxHash, inputs, outputs}

	return txCopy
}


// 数字签名验证

func (tx *Jyh_Transaction) Jyh_Verify(prevTXs map[string]Jyh_Transaction) bool {
	if tx.Jyh_IsCoinbaseTransaction() {
		return true
	}

	for _, vin := range tx.Jyh_TXIn {
		if prevTXs[hex.EncodeToString(vin.Jyh_TxHash)].Jyh_TxHash == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.Jyh_TrimmedCopy()

	curve := elliptic.P256()

	for inID, vin := range tx.Jyh_TXIn {
		prevTx := prevTXs[hex.EncodeToString(vin.Jyh_TxHash)]
		txCopy.Jyh_TXIn[inID].Jyh_Signature = nil
		txCopy.Jyh_TXIn[inID].Jyh_PubKey = prevTx.Jyh_TXOut[vin.Jyh_Vout].Jyh_Ripemd160Hash
		txCopy.Jyh_TxHash = txCopy.Jyh_Hash()
		txCopy.Jyh_TXIn[inID].Jyh_PubKey = nil


		// 私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Jyh_Signature)
		r.SetBytes(vin.Jyh_Signature[:(sigLen / 2)])
		s.SetBytes(vin.Jyh_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Jyh_PubKey)
		x.SetBytes(vin.Jyh_PubKey[:(keyLen / 2)])
		y.SetBytes(vin.Jyh_PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.Jyh_TxHash, &r, &s) == false {
			return false
		}
	}

	return true
}


func Jyh_NewSimpleTransaction(from string,to string,amount int64,utxoSet *Jyh_UTXOSet,txs []*Jyh_Transaction ,nodeID string) *Jyh_Transaction{

	wallets,_ := Jyh_NewWallets(nodeID)
	wallet := wallets.Jyh_WalletsMap[from]

	// 通过一个函数，返回
	money,spendableUTXODic := utxoSet.Jyh_FindSpendableUTXOS(from,amount,txs)

	var txIntputs []*Jyh_TXInput
	var txOutputs []*Jyh_TXOutput

	for txHash,indexArray := range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {

			txInput := &Jyh_TXInput{txHashBytes,index,nil,wallet.Jyh_PublicKey}
			txIntputs = append(txIntputs,txInput)
		}

	}

	// 转账
	txOutput := NewTXOutput(int64(amount),to)
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = NewTXOutput(int64(money) - int64(amount),from)
	txOutputs = append(txOutputs,txOutput)

	tx := &Jyh_Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.Jyh_HashTransaction()

	//进行签名
	utxoSet.Blockchain.Jyh_SignTransaction(tx, wallet.Jyh_PrivateKey,txs)

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