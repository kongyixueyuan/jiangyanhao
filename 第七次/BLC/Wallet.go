package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"

	"bytes"
)

const Version=byte(0x00)
const AddressChecksum=4

type Jyh_Wallet struct{
	Jyh_PrivateKey ecdsa.PrivateKey

	Jyh_PublicKey []byte
}

func Jyh_NewWallet() *Jyh_Wallet {
	privateKey,publicKey := newKeyPair()

	return &Jyh_Wallet{privateKey,publicKey}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if(err!=nil){
		log.Fatal("error when create wallet")
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

func IsValidForAdress(address []byte) bool{
	fullpayload := Jyh_Base58Decode(address)


	payloadversion := fullpayload[:len(fullpayload)-AddressChecksum]
	payloadchecksum := fullpayload[len(fullpayload)-AddressChecksum:]

	if bytes.Compare(payloadchecksum, Checksum(payloadversion))==0 {
		return true
	}
	return false
}

func (w Jyh_Wallet) GetAddress() []byte {
	//1. hash160
	// 20字节
	ripemd160Hash := Jyh_HashPubKey(w.Jyh_PublicKey)

	// 21字节
	version_ripemd160Hash := append([]byte{Version},ripemd160Hash...)

	// 两次的256 hash
	checkSumBytes := Checksum(version_ripemd160Hash)

	//25
	bytes := append(version_ripemd160Hash,checkSumBytes...)

	return Jyh_Base58Encode(bytes)
}

func Jyh_HashPubKey(pubKey []byte) []byte{
	//1. 256

	hash256 := sha256.New()
	hash256.Write(pubKey)
	hash := hash256.Sum(nil)

	//2. 160

	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)

}

func Checksum(versionPlayload []byte) []byte{
	firstHash := sha256.Sum256(versionPlayload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:AddressChecksum]
}