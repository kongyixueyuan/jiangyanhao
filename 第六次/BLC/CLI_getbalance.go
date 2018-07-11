package BLC

import (
	"fmt"
)

func (cli *CLI) getBalance(address string){

	blockchain := BlockchainObject()
	defer blockchain.BlockDB.Close()

	utxoSet := &jyh_UTXOSet{blockchain}
	fmt.Printf("-----\n")
	amount := utxoSet.jyh_GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n",address,amount)
}