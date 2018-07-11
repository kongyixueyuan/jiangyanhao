package BLC

func (cli *CLI) createwallet(){
	wallets,_ := NewWallets()

	wallets.CreateNewWallet()

}