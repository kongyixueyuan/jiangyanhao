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
	if(bytes.Compare(HashPubKey(TxIn.PubKey), Rimp160))==0{
		return true
	}
	return false
}