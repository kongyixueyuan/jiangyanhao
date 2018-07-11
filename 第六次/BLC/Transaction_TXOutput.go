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
	publicKeyHash := Base58Decode([]byte(address))
	hash160 := publicKeyHash[1:len(publicKeyHash) - 4]

	return bytes.Compare(TxOut.Ripemd160Hash,hash160) == 1
}
//./main getBalance -address "

func NewTXOutput(value int64,address string) *TXOutput {

	txOutput := &TXOutput{value,nil}

	// 设置Ripemd160Hash
	txOutput.Lock(address)

	return txOutput
}