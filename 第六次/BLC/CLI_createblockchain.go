package BLC

func (cli *CLI) createBlockChain(address string){

	genesisblock := jyh_NewGenesisBlock(address)
	Blockchain := NewBlockChain(genesisblock)


	utxoSet := &jyh_UTXOSet{Blockchain}
	Blockchain.Printchain()
	utxoSet.jyh_ResetUTXOSet()

	defer Blockchain.BlockDB.Close()

}