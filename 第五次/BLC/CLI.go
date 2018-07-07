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
	fmt.Println("\tcreateBlockchain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetbalance -- 查看账户余额.")
	fmt.Println("\tcreateWallet -- 创建钱包.")
	fmt.Println("\tlistWallet -- 查看钱包.")
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
	createWalletCmd := flag.NewFlagSet("createWallet",flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance",flag.ExitOnError)

	listWalletCmd := flag.NewFlagSet("listWallet",flag.ExitOnError)

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","创建创世区块的地址")

	flagFrom := sendBlockCmd.String("from","","转账源地址......")
	flagTo := sendBlockCmd.String("to","","转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount","","转账金额......")

//./main send -from jiang -to li -amount 5
//./main send -from ["jiang"] -to ["li"] -amount ["5"]
//./main send -from {"jiang"} -to {"li"} -amount {"5"}
//./main send -from '["jiang"]' -to '["li"]' -amount '["5"]'
	//1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG
	//1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM
	//13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ
	//./main send -from '["1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"]' -to '["1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM"]' -amount '["5"]'
	//./main send -from '["1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"]' -to '["13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"]' -amount '["5"]'
	//./main send -from '["13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"]' -to '["1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"]' -amount '["1"]'
	// ./main send -from '["1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG","1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"]' -to '["1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM","13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"]' -amount '["2","5"]'
	//./main send -from '["1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG","1KaJiPpYWzY2JgCrnr1T7newEK3xkwZmSG"]' -to '["1Gr8d9YMXPGNUsoQkX3N2r3qJGSnwfYzdM","13s5ZSMxyxmRNuaHEWdkpE2yQruahZ2VdJ"]' -amount '["2","5"]'
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
	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listWallet":
		err := listWalletCmd.Parse(os.Args[2:])
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

	if listWalletCmd.Parsed() {
		cli.listWallet()
	}
	if createWalletCmd.Parsed() {
		cli.createwallet()
	}

	if getBalanceCmd.Parsed() {
		fmt.Println("It's balance........")
		cli.getBalance()
	}

}