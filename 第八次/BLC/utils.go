package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"fmt"
	"encoding/gob"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

//字节数组转version
func Jyh_bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

//version 转字节数组
func Jyh_commandToBytes(command string) []byte {
	var bytes [COMMANDLENGTH]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}


// 字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
func RemoveInt(slice []int, value int) []int{

	for i, in := range slice{
		if in==value{
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func ConvertAddtoRip(address string) []byte{
	fullpayload := Jyh_Base58Decode([]byte(address))
	return fullpayload[Version:len(fullpayload)-AddressChecksum]
}

// 标准的JSON字符串转数组
func JSONToArray(jsonString string) []string {

	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		fmt.Println("1222222223333333")
		log.Panic(err)
	}
	return sArr
}

// 将结构体序列化成字节数组
func Jyh_gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
