package BLC

func (cli *CLI) createBlockChain(address string) *BlockChain{

	genesisblock := NewGenesisBlock(address)
	Blockchain := NewBlockChain(genesisblock)

	return Blockchain
}