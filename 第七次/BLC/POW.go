package BLC

import (
	"math/big"
	"strconv"
	"bytes"
	"crypto/sha256"
)

type POW struct{
	block *Jyh_Block
	target *big.Int
}

// 256位Hash里面前面至少要有16个零
const targetBit  = 8

func (Pow POW)jyh_Run() (int64 , []byte){
	nonce:=0

	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte

	for{

		heightBytes := IntToHex(Pow.block.Jyh_Height)

		timeString := strconv.FormatInt(Pow.block.Jyh_Timestamp,2)

		timeBytes := []byte(timeString)

		blockBytes := bytes.Join([][]byte{
			heightBytes,
			Pow.block.Jyh_PrevBlockHash,
			Pow.block.HashTransactions(),
			timeBytes,
			IntToHex(int64(nonce)),
		},[]byte{})

		hash = sha256.Sum256(blockBytes)
		hashInt.SetBytes(hash[:])

		if Pow.target.Cmp(&hashInt) == 1 {
			break
		}
		nonce=nonce+1
	}
	return int64(nonce), hash[:]
}

func ProofOfWork(block *Jyh_Block) *POW{
	target := big.NewInt(1)

	target=target.Lsh(target,256-targetBit)

	return &POW{block,target}
}


