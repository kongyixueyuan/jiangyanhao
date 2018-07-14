package BLC

func (cli *Jyh_CLI) Jyh_printBlockChain(nodeID string){
	Blockchain:=BlockchainObject(nodeID)

	Blockchain.Jyh_Printchain()

	defer Blockchain.Jyh_BlockDB.Close()
}