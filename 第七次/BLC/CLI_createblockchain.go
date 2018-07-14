package BLC

func (cli *Jyh_CLI) createBlockChain(address string,nodeID string){

	Blockchain := Jyh_CreateBlockchainWithGenesisBlock(address,nodeID)


	utxoSet := &Jyh_UTXOSet{Blockchain}
	Blockchain.Jyh_Printchain()
	utxoSet.Jyh_ResetUTXOSet()

	defer Blockchain.Jyh_BlockDB.Close()

}