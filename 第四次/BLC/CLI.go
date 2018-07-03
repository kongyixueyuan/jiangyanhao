package BLC

import(
	"os"
	"fmt"
	"log"
	"flag"

)

type CLI struct {}

func printUsage()  {

	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetbalance -address -- 输出区块信息.")

}


func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}


func (cli *CLI) Run (){
	isValidArgs()

	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createBlockchain",flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance",flag.ExitOnError)

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","创建创世区块的地址")

	flagFrom := sendBlockCmd.String("from","","转账源地址......")
	flagTo := sendBlockCmd.String("to","","转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount","","转账金额......")

//./main send -from jiang -to li -amount 5
//./main send -from ["jiang"] -to ["li"] -amount ["5"]
//./main send -from {"jiang"} -to {"li"} -amount {"5"}
//./main send -from '["jiang"]' -to '["li"]' -amount '["5"]'
	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == ""{
			printUsage()
			fmt.Println("00000000111111111")
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount)
	}

	if createBlockchainCmd.Parsed() {

		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateBlockchainWithAddress)

	}

	if printChainCmd.Parsed() {
		cli.printBlockChain()
		fmt.Println("输出所有区块的数据........")
	}

	if getBalanceCmd.Parsed() {
		fmt.Println("It's balance........")
		cli.getBalance()
	}

}