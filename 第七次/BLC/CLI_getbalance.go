package BLC

import (
	"fmt"
)

func (cli *Jyh_CLI) getBalance(address string, nodeID string){

	blockchain := BlockchainObject(nodeID)
	defer blockchain.Jyh_BlockDB.Close()

	utxoSet := &Jyh_UTXOSet{blockchain}
	fmt.Printf("-----\n")
	amount := utxoSet.Jyh_GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n",address,amount)
}