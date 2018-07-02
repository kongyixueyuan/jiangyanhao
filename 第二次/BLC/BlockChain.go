package BLC

type BlockChain struct{
	Blocks []*Block
}

func NewBlockChain(block * Block) * BlockChain{

	blockchain := &BlockChain{[]*Block{block}}

	return blockchain
}

func AddToChain(block *Block, blockchain *BlockChain) *BlockChain{
	blockchain.Blocks = append(blockchain.Blocks, block)
	return blockchain
}