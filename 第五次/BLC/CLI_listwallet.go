package BLC

import "fmt"

func (cli *CLI) listWallet(){

	wallets,_ := NewWallets()
	for address,_ := range wallets.WalletsMap {
		fmt.Println(address)
	}

}
