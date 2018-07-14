package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Jyh_TXOutputs struct {
	Jyh_UTXOS []*Jyh_UTXO
}


// 将区块序列化成字节数组
func (txOutputs *Jyh_TXOutputs) Serialize() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func DeserializeTXOutputs(txOutputsBytes []byte) *Jyh_TXOutputs {

	var txOutputs Jyh_TXOutputs

	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	err := decoder.Decode(&txOutputs)
	if err != nil {
		log.Panic(err)
	}

	return &txOutputs
}