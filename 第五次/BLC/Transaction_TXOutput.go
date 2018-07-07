package BLC

import "bytes"

//TXOutput{100,"zhangbozhi"}
//TXOutput{30,"xietingfeng"}
//TXOutput{40,"zhangbozhi"}


type TXOutput struct {
	Value int64
	Ripemd160Hash []byte  //将public key进行过sha256，rip160的hash
}

func (TxOut *TXOutput) Lock(address string){
	fullpayload := Base58Decode([]byte(address))
	TxOut.Ripemd160Hash = fullpayload[Version:len(fullpayload)-AddressChecksum]
}

func (TxOut *TXOutput) IsLockedPubkeyTxOut (address string) bool {
	fullpayload := Base58Decode([]byte(address))
	if(bytes.Compare(fullpayload[Version:len(fullpayload)-AddressChecksum], TxOut.Ripemd160Hash))==0{
		return true
	}
	return false
}