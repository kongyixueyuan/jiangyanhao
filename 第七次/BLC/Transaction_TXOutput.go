package BLC

import (
	"bytes"
	"fmt"
)

//TXOutput{100,"zhangbozhi"}
//TXOutput{30,"xietingfeng"}
//TXOutput{40,"zhangbozhi"}

type Jyh_TXOutput struct {
	Jyh_Value int64
	Jyh_Ripemd160Hash []byte  //将public key进行过sha256，rip160的hash

}

func (TxOut *Jyh_TXOutput) Jyh_Lock(address string){
	publicKeyHash := Jyh_Base58Decode([]byte(address))

	TxOut.Jyh_Ripemd160Hash = publicKeyHash[1:len(publicKeyHash) - 4]
}

func (TxOut *Jyh_TXOutput) Jyh_IsLockedPubkeyTxOut (address string) bool {
	publicKeyHash := Jyh_Base58Decode([]byte(address))
	hash160 := publicKeyHash[1:len(publicKeyHash) - 4]

	return bytes.Compare(TxOut.Jyh_Ripemd160Hash,hash160) == 0
}
//./main getBalance -address "

func NewTXOutput(value int64,address string) *Jyh_TXOutput {

	txOutput := &Jyh_TXOutput{value,nil}

	// 设置Ripemd160Hash
	txOutput.Jyh_Lock(address)
	fmt.Printf("rimp::::%x\n",txOutput.Jyh_Ripemd160Hash)
	fmt.Printf("address::::%s\n",address)

//rimp::::67aa8e771cb30d9f3cfd04c770e55cd0af912944
//address::::1AT8vEUJWbdA9J7nNSBqmxJcrYmEqMFcYY

	return txOutput
}