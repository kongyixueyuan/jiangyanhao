package BLC

func (cli *CLI) createBlockChain(address string) *BlockChain{

	genesisblock := NewGenesisBlock("jiang")
	Blockchain := NewBlockChain(genesisblock)

	return Blockchain
}