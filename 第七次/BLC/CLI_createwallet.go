package BLC

import "fmt"

func (cli *Jyh_CLI) createwallet(nodeID string){
	wallets,_ := Jyh_NewWallets(nodeID)

	wallets.Jyh_CreateNewWallet(nodeID)

	fmt.Println(len(wallets.Jyh_WalletsMap))

}