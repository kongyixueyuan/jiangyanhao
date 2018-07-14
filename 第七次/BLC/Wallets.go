package BLC

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"log"
	"io/ioutil"
	"os"
)

type Jyh_Wallets struct {
	Jyh_WalletsMap map[string]*Jyh_Wallet
}

const walletFile  = "Wallets_%s.dat"

// 创建钱包集合
func Jyh_NewWallets(nodeID string) (*Jyh_Wallets,error){
	walletFile := fmt.Sprintf(walletFile,nodeID)

	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Jyh_Wallets{}
		wallets.Jyh_WalletsMap = make(map[string]*Jyh_Wallet)
		return wallets,err
	}


	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Jyh_Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets,nil
}





// 创建一个新钱包
func (w *Jyh_Wallets) Jyh_CreateNewWallet(nodeID string)  {

	wallet := Jyh_NewWallet()
	fmt.Printf("Address：%s\n",wallet.GetAddress())
	w.Jyh_WalletsMap[string(wallet.GetAddress())] = wallet
	w.Jyh_SaveWallets(nodeID)
}

// 将钱包信息写入到文件
func (w *Jyh_Wallets) Jyh_SaveWallets(nodeID string)  {
	walletFile := fmt.Sprintf(walletFile,nodeID)

	var content bytes.Buffer

	// 注册的目的，是为了，可以序列化任何类型
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}

	// 将序列化以后的数据写入到文件，原来文件的数据会被覆盖
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}


}