package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) printBlockChain(){
	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}
	Blockchain:=BlockchainObject()
	Blockchain.Printchain()

	defer Blockchain.BlockDB.Close()
}