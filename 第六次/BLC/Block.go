package BLC

import(
	"time"
	"bytes"
	"encoding/gob"
	"log"
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

/*
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transaction {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
*/
	var transactions [][]byte

	for _, tx := range block.Transaction {
		transactions = append(transactions, tx.jyh_Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// 将区块序列化成字节数组
func (block *Block) jyh_Serialize() []byte {

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
func jyh_DeserializeBlock(blockBytes []byte) *Block {

	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	//fmt.Println("---------in deserialize, %d",block.Nonce)

	if err != nil {
		log.Panic(err)
	}

	return &block
}
/*
func MineNewBlock(Transaction []*Transaction, blockchain *BlockChain){
	// 在建立新区块之前对txs进行签名验证

	for _,tx := range Transaction  {

		if blockchain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

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
*/
// ./main send -from '["1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"]' -to '["1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM"]' -amount '["5"]'
// ./main send -from '["1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"]' -to '["13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"]' -amount '["3"]'
// ./main send -from '["1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM"]' -to '["13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"]' -amount '["1"]'
// ./main send -from '["13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"]' -to '["1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM"]' -amount '["5"]'
// ./main send -from '["13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ","1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"]' -to '["1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG","1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM"]' -amount '["2","3"]'
// ./main createblockchain -address "1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"
// 制作transaction.main createblockchain -address 1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG
/*
func MakeTransaction(from []string, to []string, amount []string, blockchain *BlockChain) {

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
			/*for i, txOut := range tx.TXOut {

				for _, value := range spendableUTXO[string(tx.TxHash)] {
					if (i == value) {
						//	spendableTxout=append(spendableTxout, txOut)
						spendableMap[string(tx.TxHash)] = append(spendableMap[string(tx.TxHash)], txOut)
					}
				}
			}
			if len(spendableUTXO[string(tx.TxHash)])==0{
				continue
			}else{
				for _, value := range spendableUTXO[string(tx.TxHash)] {
						//	spendableTxout=append(spendableTxout, txOut)
						spendableMap[string(tx.TxHash)] = append(spendableMap[string(tx.TxHash)], tx.TXOut[value])
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

	//2、看看与from的人是否一致
	/*	for _, txOut :=range spendableTxout{
			for _ , fromPeople := range from{
				if(strings.Compare(txOut.ScriptPubKey, fromPeople)==0){
				value=value+txOut.Value
				}
			}
		}

	wallets,_ := NewWallets()
	var wallet = new(Wallet)

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

}
*/
func jyh_NewBlock(txs []*Transaction,height int64,prevBlockHash []byte) *Block {

	//创建区块
	block := &Block{height,prevBlockHash,0,txs,time.Now().Unix(),nil}

	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	Pow := ProofOfWork(block)

	// 挖矿验证
	nonce,hash := Pow.jyh_Run()

	block.BlockHash = hash[:]
	block.Nonce = nonce

	return block

}


func jyh_NewGenesisBlock(address string) * Block{
	GenesisTransaction:=jyh_SetCoinbBaseTransaction(address)
	block := &Block{1,IntToHex(0),0,[]*Transaction{GenesisTransaction},time.Now().Unix(),nil}

	Pow:=ProofOfWork(block)
	nonce, hash := Pow.jyh_Run()

	block.BlockHash=hash[:]
	block.Nonce=nonce

	//PrintBlock(block)
	return block
}