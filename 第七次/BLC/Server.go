package BLC

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
)

func Jyh_startServer(nodeID string,minerAdd string)  {

	// 当前节点的IP地址
	nodeAddress = fmt.Sprintf("localhost:%s",nodeID)

	ln,err := net.Listen(PROTOCOL,nodeAddress)

	if err != nil {
		log.Panic(err)
	}

	defer ln.Close()


	bc := BlockchainObject(nodeID)

	defer bc.Jyh_BlockDB.Close()

	// 第一个终端：端口为3000,启动的就是主节点
	// 第二个终端：端口为3001，钱包节点
	// 第三个终端：端口号为3002，矿工节点
	if nodeAddress != knowNodes[0]{
		// 此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据

		Jyh_sendVersion(knowNodes[0],bc)
	}

	for {
		// 收到的数据的格式是固定的，12字节+结构体字节数组

		// 接收客户端发送过来的数据
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}

		go Jyh_handleConnection(conn,bc)

	}

}



func Jyh_handleConnection(conn net.Conn,bc *Jyh_BlockChain) {

	// 读取客户端发送过来的所有的数据
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Receive a Message:%s\n",request[:COMMANDLENGTH])

	command := Jyh_bytesToCommand(request[:COMMANDLENGTH])


	// 12字节 + 某个结构体序列化以后的字节数组

	switch command {
	case COMMAND_VERSION:
		Jyh_handleVersion(request, bc)
	case COMMAND_ADDR:
		Jyh_handleAddr(request, bc)
	case COMMAND_BLOCK:
		Jyh_handleBlock(request, bc)
	case COMMAND_GETBLOCKS:
		Jyh_handleGetblocks(request, bc)
	case COMMAND_GETDATA:
		Jyh_handleGetData(request, bc)
	case COMMAND_INV:
		Jyh_handleInv(request, bc)
	case COMMAND_TX:
		Jyh_handleTx(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

