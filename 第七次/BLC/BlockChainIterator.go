package BLC


import (
	"github.com/boltdb/bolt"
	"log"
)

type Jyh_BlockchainIterator struct {
	Jyh_CurrentHash []byte
	Jyh_DB  *bolt.DB
}

func (blockchainIterator *Jyh_BlockchainIterator) Jyh_Next() *Jyh_Block {

	var block *Jyh_Block

	err := blockchainIterator.Jyh_DB.View(func(tx *bolt.Tx) error{

		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			currentBloclBytes := b.Get(blockchainIterator.Jyh_CurrentHash)
			//  获取到当前迭代器里面的currentHash所对应的区块
			block = Jyh_DeserializeBlock(currentBloclBytes)

			// 更新迭代器里面CurrentHash
			blockchainIterator.Jyh_CurrentHash = block.Jyh_PrevBlockHash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}


	return block

}