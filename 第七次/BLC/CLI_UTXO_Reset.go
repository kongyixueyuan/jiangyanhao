package BLC

import "fmt"

func (cli *Jyh_CLI) resetUTXOSet(nodeID string)  {


	fmt.Println("resetUTXOSet")

	blockchain := BlockchainObject(nodeID)

	defer blockchain.Jyh_BlockDB.Close()

	utxoSet := &Jyh_UTXOSet{blockchain}

	utxoSet.Jyh_ResetUTXOSet()
}
