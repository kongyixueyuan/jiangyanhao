package BLC

import "bytes"

type TXInput struct {
	// 1. 交易的Hash
	TxHash      []byte
	// 2. 存储TXOutput在Vout里面的索引
	Vout      int
	Signature []byte
	PubKey    []byte  //实际的public key
}


func (TxIn *TXInput) IsLockedPubkeyTxIn (Rimp160 []byte) bool {
	publicKey := HashPubKey(TxIn.PubKey)

	return bytes.Compare(publicKey,Rimp160) == 0
}