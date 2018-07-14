package BLC

import "bytes"

type Jyh_TXInput struct {
	// 1. 交易的Hash
	Jyh_TxHash      []byte
	// 2. 存储TXOutput在Vout里面的索引
	Jyh_Vout      int
	Jyh_Signature []byte
	Jyh_PubKey    []byte  //实际的public key
}


func (TxIn *Jyh_TXInput) Jyh_IsLockedPubkeyTxIn (Rimp160 []byte) bool {
	publicKey := Jyh_HashPubKey(TxIn.Jyh_PubKey)

	return bytes.Compare(publicKey,Rimp160) == 0
}