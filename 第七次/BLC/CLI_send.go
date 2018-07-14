package BLC

import "fmt"

// 转账
func (cli *Jyh_CLI) send(from []string,to []string,amount []string, nodeID string, mineNow bool)  {
	//1AT8vEUJWbdA9J7nNSBqmxJcrYmEqMFcYY =15
	//199mwdskttaV3v7K7ezWyRe6mU6XHNeoMn
	//1DZhVyKBJGjJteCX49V2jWnfgKCh81gueG

	blockchain := BlockchainObject(nodeID)

	//签名，简历transaction
	//MakeTransaction(from,to,amount,blockchain)

	//矿工对签名进行验证

	defer blockchain.Jyh_BlockDB.Close()
	if mineNow{
		blockchain.Jyh_MineNewBlock(from,to,amount,nodeID)

		utxoSet := &Jyh_UTXOSet{blockchain}

		//转账成功以后，需要更新一下
		utxoSet.Jyh_Update()
	}else{
		// 把交易发送到矿工节点去进行验证
		fmt.Println("由矿工节点处理......")
	}


}
