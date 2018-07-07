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

type Wallet struct{
	PrivateKey ecdsa.PrivateKey

	PublicKey []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet
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

func (w Wallet) IsValidAddress(address []byte) bool{
	fullpayload := Base58Decode(address)


	payloadversion := fullpayload[:len(fullpayload)-AddressChecksum]
	payloadchecksum := fullpayload[len(fullpayload)-AddressChecksum:]

	if bytes.Compare(payloadchecksum, Checksum(payloadversion))==0 {
		return true
	}
	return false
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{Version}, pubKeyHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)

	return Base58Encode(fullPayload)
}

func HashPubKey(pubKey []byte) []byte{
	PubkeySha256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(PubkeySha256[:])
	if(err!=nil){
		log.Fatal("HashPubkey error!")
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160

}

func Checksum(versionPlayload []byte) []byte{
	firstHash := sha256.Sum256(versionPlayload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:AddressChecksum]
}