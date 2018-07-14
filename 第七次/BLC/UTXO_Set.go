package BLC

import (
	"github.com/boltdb/bolt"
	"encoding/hex"
	"log"
	"fmt"
	"bytes"
)

const utxoTableName  = "utxoTableName"

type Jyh_UTXOSet struct {
	Blockchain *Jyh_BlockChain
}


// 重置数据库表
func (utxoSet *Jyh_UTXOSet) Jyh_ResetUTXOSet()  {

	err := utxoSet.Blockchain.Jyh_BlockDB.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if err!= nil {
				log.Panic(err)
			}

		}

		b ,_ = tx.CreateBucket([]byte(utxoTableName))
		if b != nil {
			//[string]*TXOutputs
			txOutputsMap := utxoSet.Blockchain.Jyh_FindUTXOMap()
			for keyHash,outs := range txOutputsMap {
				for _,v := range outs.Jyh_UTXOS{
					fmt.Printf("this data is from FindUTXOMap()\n")
					fmt.Printf("tx hash:%x,\nindex:%d,\n",v.Jyh_TxHash,v.Jyh_Index)
					fmt.Printf("it's value:%d, and 160hash: %d",v.Jyh_Output.Jyh_Value,v.Jyh_Output.Jyh_Ripemd160Hash)
				}
				txHash,_ := hex.DecodeString(keyHash)

				b.Put(txHash, outs.Serialize() )

			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

func (utxoSet *Jyh_UTXOSet) Jyh_FindSpendableUTXOS(from string,amount int64,txs []*Jyh_Transaction) (int64,map[string][]int)  {

	unPackageUTXOS := utxoSet.Jyh_FindUnPackageSpendableUTXOS(from,txs)

	spentableUTXO := make(map[string][]int)

	var money int64 = 0

	for _,UTXO := range unPackageUTXOS {

		money += UTXO.Jyh_Output.Jyh_Value;
		txHash := hex.EncodeToString(UTXO.Jyh_TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash],UTXO.Jyh_Index)
		if money >= amount{
			return  money,spentableUTXO
		}
	}


	// 钱还不够
	utxoSet.Blockchain.Jyh_BlockDB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {

			c := b.Cursor()
		UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := DeserializeTXOutputs(v)

				for _,utxo := range txOutputs.Jyh_UTXOS {

					money += utxo.Jyh_Output.Jyh_Value
					txHash := hex.EncodeToString(utxo.Jyh_TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash],utxo.Jyh_Index)

					if money >= amount {
						break UTXOBREAK;
					}
				}
			}

		}

		return nil
	})

	if money < amount{
		log.Panic("余额不足......")
	}


	return  money,spentableUTXO
}
//./main getBalance -address "1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM"
func (utxoSet *Jyh_UTXOSet) Jyh_findUTXOForAddress(address string) []*Jyh_UTXO{

	fmt.Println("11111111111")
	var utxos []*Jyh_UTXO

	err:= utxoSet.Blockchain.Jyh_BlockDB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		// 游标
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOutputs := DeserializeTXOutputs(v)
//			fmt.Println(txOutputs)
//			fmt.Printf("length of UTXOS:%d",len(txOutputs.Jyh_UTXOS))
			for _,utxo := range txOutputs.Jyh_UTXOS  {

				if utxo.Jyh_Output.Jyh_IsLockedPubkeyTxOut(address) {

					utxos = append(utxos,utxo)
//					fmt.Println(utxo.Jyh_Output.Jyh_Ripemd160Hash)
//					fmt.Printf("address:%s\n",address)
				}
			}
		}
//		fmt.Printf("length:%d\n",len(utxos))
		/*for _,v := range utxos{
			fmt.Printf("tx hash:%x,\nindex:%d,\n",v.TxHash,v.Index)
			fmt.Printf("it's value:%d, and 160hash: %d",v.Output.Value,v.Output.Ripemd160Hash)
		}*/
		return nil
	})
	if(err!=nil){
		fmt.Println("err is not nil!!!")
	}


	return utxos

}



func (utxoSet *Jyh_UTXOSet) Jyh_GetBalance(address string) int64 {

	UTXOS := utxoSet.Jyh_findUTXOForAddress(address)

	var amount int64

	for _,utxo := range UTXOS  {
		amount += utxo.Jyh_Output.Jyh_Value
	}

	return amount
}

// 返回要凑多少钱，对应TXOutput的TX的Hash和index
func (utxoSet *Jyh_UTXOSet) Jyh_FindUnPackageSpendableUTXOS(from string, txs []*Jyh_Transaction) []*Jyh_UTXO {

	var unUTXOs []*Jyh_UTXO

	spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range txs {

		if tx.Jyh_IsCoinbaseTransaction() == false {
			for _, in := range tx.Jyh_TXIn {
				//是否能够解锁
				publicKeyHash := Jyh_Base58Decode([]byte(from))

				ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]
				if in.Jyh_IsLockedPubkeyTxIn(ripemd160Hash) {

					key := hex.EncodeToString(in.Jyh_TxHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Jyh_Vout)
				}

			}
		}
	}


	for _,tx := range txs {

	Work1:
		for index,out := range tx.Jyh_TXOut {

			if out.Jyh_IsLockedPubkeyTxOut(from) {
				fmt.Println("看看是否是俊诚...")
				fmt.Println(from)

				fmt.Println(spentTXOutputs)

				if len(spentTXOutputs) == 0 {
					utxo := &Jyh_UTXO{tx.Jyh_TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.Jyh_TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &Jyh_UTXO{tx.Jyh_TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &Jyh_UTXO{tx.Jyh_TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}



	return unUTXOs

}

// 更新
func (utxoSet *Jyh_UTXOSet) Jyh_Update()  {

	// blocks
	//


	// 最新的Block
	block := utxoSet.Blockchain.Jyh_Iterator().Jyh_Next()


	// utxoTable
	//

	ins := []*Jyh_TXInput{}

	outsMap := make(map[string]*Jyh_TXOutputs)

	// 找到所有我要删除的数据
	for _,tx := range block.Jyh_Transaction {

		for _,in := range tx.Jyh_TXIn {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.Jyh_Transaction  {


		utxos := []*Jyh_UTXO{}

		for index,out := range tx.Jyh_TXOut  {

			isSpent := false

			for _,in := range ins  {

				if in.Jyh_Vout == index && bytes.Compare(tx.Jyh_TxHash ,in.Jyh_TxHash) == 0 && bytes.Compare(out.Jyh_Ripemd160Hash,Jyh_HashPubKey(in.Jyh_PubKey)) == 0 {

					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &Jyh_UTXO{tx.Jyh_TxHash,index,out}
				utxos = append(utxos,utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.Jyh_TxHash)
			outsMap[txHash] = &Jyh_TXOutputs{utxos}
		}

	}



	err := utxoSet.Blockchain.Jyh_BlockDB.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {


			// 删除
			for _,in := range ins {

				txOutputsBytes := b.Get(in.Jyh_TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				fmt.Println("DeserializeTXOutputs")
				fmt.Println(txOutputsBytes)

				txOutputs := DeserializeTXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*Jyh_UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _,utxo := range txOutputs.Jyh_UTXOS  {

					if in.Jyh_Vout == utxo.Jyh_Index && bytes.Compare(utxo.Jyh_Output.Jyh_Ripemd160Hash,Jyh_HashPubKey(in.Jyh_PubKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}



				if isNeedDelete {
					b.Delete(in.Jyh_TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.Jyh_TxHash)]

						preTXOutputs.Jyh_UTXOS = append(preTXOutputs.Jyh_UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.Jyh_TxHash)] = preTXOutputs

					}
				}

			}

			// 新增

			for keyHash,outPuts := range outsMap  {
				keyHashBytes,_ := hex.DecodeString(keyHash)
				b.Put(keyHashBytes,outPuts.Serialize())
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}




