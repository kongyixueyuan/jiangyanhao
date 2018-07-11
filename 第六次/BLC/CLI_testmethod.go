package BLC

import "fmt"

func (cli *CLI) TestMethod()  {


	fmt.Println("TestMethod")

	blockchain := BlockchainObject()

	defer blockchain.BlockDB.Close()

	utxoSet := &jyh_UTXOSet{blockchain}

	utxoSet.jyh_ResetUTXOSet()

	//fmt.Println(blockchain.FindUTXOMap())
}
