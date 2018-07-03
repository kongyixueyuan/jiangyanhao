package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) getBalance(){
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	Blockchain:=BlockchainObject()
	Blockchain.getBalance()

	defer Blockchain.BlockDB.Close()
}