package BLC

import (
	"github.com/boltdb/bolt"
	"encoding/hex"
	"log"
	"fmt"
	"bytes"
)

const utxoTableName  = "utxoTableName"

type jyh_UTXOSet struct {
	Blockchain *BlockChain
}


// 重置数据库表
func (utxoSet *jyh_UTXOSet) jyh_ResetUTXOSet()  {

	err := utxoSet.Blockchain.BlockDB.Update(func(tx *bolt.Tx) error {

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
			txOutputsMap := utxoSet.Blockchain.jyh_FindUTXOMap()
			for keyHash,outs := range txOutputsMap {
				for _,v := range outs.UTXOS{
					fmt.Printf("this data is from FindUTXOMap()\n")
					fmt.Printf("tx hash:%x,\nindex:%d,\n",v.TxHash,v.Index)
					fmt.Printf("it's value:%d, and 160hash: %d",v.Output.Value,v.Output.Ripemd160Hash)
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

func (utxoSet *jyh_UTXOSet) jyh_FindSpendableUTXOS(from string,amount int64,txs []*Transaction) (int64,map[string][]int)  {

	unPackageUTXOS := utxoSet.jyh_FindUnPackageSpendableUTXOS(from,txs)

	spentableUTXO := make(map[string][]int)

	var money int64 = 0

	for _,UTXO := range unPackageUTXOS {

		money += UTXO.Output.Value;
		txHash := hex.EncodeToString(UTXO.TxHash)
		spentableUTXO[txHash] = append(spentableUTXO[txHash],UTXO.Index)
		if money >= amount{
			return  money,spentableUTXO
		}
	}


	// 钱还不够
	utxoSet.Blockchain.BlockDB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {

			c := b.Cursor()
		UTXOBREAK:
			for k, v := c.First(); k != nil; k, v = c.Next() {

				txOutputs := DeserializeTXOutputs(v)

				for _,utxo := range txOutputs.UTXOS {

					money += utxo.Output.Value
					txHash := hex.EncodeToString(utxo.TxHash)
					spentableUTXO[txHash] = append(spentableUTXO[txHash],utxo.Index)

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
func (utxoSet *jyh_UTXOSet) jyh_findUTXOForAddress(address string) []*UTXO{


	var utxos []*UTXO

	utxoSet.Blockchain.BlockDB.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(utxoTableName))

		// 游标
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			txOutputs := DeserializeTXOutputs(v)

			for _,utxo := range txOutputs.UTXOS  {

				if utxo.Output.IsLockedPubkeyTxOut(address) {
					utxos = append(utxos,utxo)
				}
			}
		}
		fmt.Printf("length:%d\n",len(utxos))
		for _,v := range utxos{
			fmt.Printf("tx hash:%x,\nindex:%d,\n",v.TxHash,v.Index)
			fmt.Printf("it's value:%d, and 160hash: %d",v.Output.Value,v.Output.Ripemd160Hash)
		}
		return nil
	})


	return utxos

}



func (utxoSet *jyh_UTXOSet) jyh_GetBalance(address string) int64 {

	UTXOS := utxoSet.jyh_findUTXOForAddress(address)

	var amount int64

	for _,utxo := range UTXOS  {
		amount += utxo.Output.Value
	}

	return amount
}

// 返回要凑多少钱，对应TXOutput的TX的Hash和index
func (utxoSet *jyh_UTXOSet) jyh_FindUnPackageSpendableUTXOS(from string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	//{hash:[0]}

	for _,tx := range txs {

		if tx.jyh_IsCoinbaseTransaction() == false {
			for _, in := range tx.TXIn {
				//是否能够解锁
				publicKeyHash := Base58Decode([]byte(from))

				ripemd160Hash := publicKeyHash[1:len(publicKeyHash) - 4]
				if in.IsLockedPubkeyTxIn(ripemd160Hash) {

					key := hex.EncodeToString(in.TxHash)

					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}

			}
		}
	}


	for _,tx := range txs {

	Work1:
		for index,out := range tx.TXOut {

			if out.IsLockedPubkeyTxOut(from) {
				fmt.Println("看看是否是俊诚...")
				fmt.Println(from)

				fmt.Println(spentTXOutputs)

				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash,indexArray := range spentTXOutputs {

						txHashStr := hex.EncodeToString(tx.TxHash)

						if hash == txHashStr {

							var isUnSpentUTXO bool

							for _,outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO == false {
									utxo := &UTXO{tx.TxHash, index, out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
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
func (utxoSet *jyh_UTXOSet) jyh_Update()  {

	// blocks
	//


	// 最新的Block
	block := utxoSet.Blockchain.Iterator().jyh_Next()


	// utxoTable
	//

	ins := []*TXInput{}

	outsMap := make(map[string]*TXOutputs)

	// 找到所有我要删除的数据
	for _,tx := range block.Transaction {

		for _,in := range tx.TXIn {
			ins = append(ins,in)
		}
	}

	for _,tx := range block.Transaction  {


		utxos := []*UTXO{}

		for index,out := range tx.TXOut  {

			isSpent := false

			for _,in := range ins  {

				if in.Vout == index && bytes.Compare(tx.TxHash ,in.TxHash) == 0 && bytes.Compare(out.Ripemd160Hash,HashPubKey(in.PubKey)) == 0 {

					isSpent = true
					continue
				}
			}

			if isSpent == false {
				utxo := &UTXO{tx.TxHash,index,out}
				utxos = append(utxos,utxo)
			}

		}

		if len(utxos) > 0 {
			txHash := hex.EncodeToString(tx.TxHash)
			outsMap[txHash] = &TXOutputs{utxos}
		}

	}



	err := utxoSet.Blockchain.BlockDB.Update(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(utxoTableName))

		if b != nil {


			// 删除
			for _,in := range ins {

				txOutputsBytes := b.Get(in.TxHash)

				if len(txOutputsBytes) == 0 {
					continue
				}

				fmt.Println("DeserializeTXOutputs")
				fmt.Println(txOutputsBytes)

				txOutputs := DeserializeTXOutputs(txOutputsBytes)

				fmt.Println(txOutputs)

				UTXOS := []*UTXO{}

				// 判断是否需要
				isNeedDelete := false

				for _,utxo := range txOutputs.UTXOS  {

					if in.Vout == utxo.Index && bytes.Compare(utxo.Output.Ripemd160Hash,HashPubKey(in.PubKey)) == 0 {

						isNeedDelete = true
					} else {
						UTXOS = append(UTXOS,utxo)
					}
				}



				if isNeedDelete {
					b.Delete(in.TxHash)
					if len(UTXOS) > 0 {

						preTXOutputs := outsMap[hex.EncodeToString(in.TxHash)]

						preTXOutputs.UTXOS = append(preTXOutputs.UTXOS,UTXOS...)

						outsMap[hex.EncodeToString(in.TxHash)] = preTXOutputs

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




