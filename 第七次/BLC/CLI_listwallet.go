package BLC

import "fmt"

func (cli *Jyh_CLI) Jyh_listWallet(nodeID string){

	wallets,_ := Jyh_NewWallets(nodeID)
	for address,_ := range wallets.Jyh_WalletsMap {
		fmt.Println(address)
	}

}
