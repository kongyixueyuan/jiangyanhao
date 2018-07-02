package main

import
(
	"PublicChain1801/姜延浩/第二次/BLC"
)

func main() {

	genesisblock := BLC.NewGenesisBlock([]byte("it's genesis block"))
	blockchain := BLC.NewBlockChain(genesisblock)


	newblock := BLC.NewBlock(blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].BlockHash , []byte("2nd block"))
	BLC.AddToChain(newblock, blockchain)

	newblock = BLC.NewBlock(blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].BlockHash , []byte("3rd block"))
	BLC.AddToChain(newblock, blockchain)
}