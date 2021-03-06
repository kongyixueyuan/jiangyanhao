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