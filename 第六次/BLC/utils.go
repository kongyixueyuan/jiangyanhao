package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"fmt"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
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
	fullpayload := Base58Decode([]byte(address))
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

