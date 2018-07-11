package BLC

import (
	"fmt"
	"os"
)

// 转账
func (cli *CLI) send(from []string,to []string,amount []string)  {



	if DBExists() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObject()

	//签名，简历transaction
	//MakeTransaction(from,to,amount,blockchain)

	//矿工对签名进行验证

	defer blockchain.BlockDB.Close()

	blockchain.jyh_MineNewBlock(from,to,amount)

	utxoSet := &jyh_UTXOSet{blockchain}

	//转账成功以后，需要更新一下
	utxoSet.jyh_Update()

}
