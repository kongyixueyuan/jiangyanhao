package BLC

import (
	"net"
	"io"
	"bytes"
	"log"
)

//COMMAND_VERSION
func Jyh_sendVersion(toAddress string,bc *Jyh_BlockChain)  {


	bestHeight := bc.Jyh_GetBestHeight()

	payload := Jyh_gobEncode(Jyh_Version{NODE_VERSION, bestHeight, nodeAddress})

	//version
	request := append(Jyh_commandToBytes(COMMAND_VERSION), payload...)

	Jyh_sendData(toAddress,request)


}

func Jyh_sendBlock(toAddress string, block []byte)  {


	payload := Jyh_gobEncode(Jyh_BlockData{nodeAddress,block})

	request := append(Jyh_commandToBytes(COMMAND_BLOCK), payload...)

	Jyh_sendData(toAddress,request)

}

func Jyh_sendGetData(toAddress string, kind string ,blockHash []byte) {

	payload := Jyh_gobEncode(Jyh_GetData{nodeAddress,kind,blockHash})

	request := append(Jyh_commandToBytes(COMMAND_GETDATA), payload...)

	Jyh_sendData(toAddress,request)
}


func Jyh_sendData(to string,data []byte)  {

	conn, err := net.Dial("tcp", to)
	if err != nil {
		panic("error")
	}
	defer conn.Close()

	// 附带要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

//
func Jyh_sendInv(toAddress string, kind string, hashes [][]byte) {

	payload := Jyh_gobEncode(Jyh_Inv{nodeAddress,kind,hashes})

	request := append(Jyh_commandToBytes(COMMAND_INV), payload...)

	Jyh_sendData(toAddress,request)

}

//COMMAND_GETBLOCKS
func Jyh_sendGetBlocks(toAddress string)  {

	payload := Jyh_gobEncode(Jyh_GetBlocks{nodeAddress})

	request := append(Jyh_commandToBytes(COMMAND_GETBLOCKS), payload...)

	Jyh_sendData(toAddress,request)

}