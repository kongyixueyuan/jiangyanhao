package BLC

import(
	"os"
	"fmt"
	"log"
	"flag"

)

type Jyh_CLI struct {}

func printUsage()  {

	fmt.Println("Usage:")

	fmt.Println("\tcreateBlockchain -address -- 交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细.")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetBalance -- 查看账户余额.")
	fmt.Println("\tcreateWallet -- 创建钱包.")
	fmt.Println("\tlistWallet -- 查看钱包.")
	fmt.Println("\tresetUTXO -- 重置.")
	fmt.Println("\tstartnode -miner ADDRESS -- 启动节点服务器，并且指定挖矿奖励的地址.")
}


func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}


func (cli *Jyh_CLI) Run (){
	isValidArgs()

	// 设置ID
	// export NODE_ID=8888
	// 读取
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!\n")
		os.Exit(1)
	}

	fmt.Printf("NODE_ID:%s\n",nodeID)

	sendBlockCmd := flag.NewFlagSet("send",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createBlockchain",flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet",flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getBalance",flag.ExitOnError)
	resetUTXOCMD := flag.NewFlagSet("resetUTXO",flag.ExitOnError)
	listWalletCmd := flag.NewFlagSet("listWallet",flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode",flag.ExitOnError)

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address","","创建创世区块的地址")
	getBalanceWithAdress := getBalanceCmd.String("address","","要查询某一个账号的余额.......")


	flagFrom := sendBlockCmd.String("from","","转账源地址......")
	flagTo := sendBlockCmd.String("to","","转账目的地地址......")
	flagAmount := sendBlockCmd.String("amount","","转账金额......")


	flagMine := sendBlockCmd.Bool("mine",false,"是否在当前节点中立即验证....")


	flagMiner := startNodeCmd.String("miner","","定义挖矿奖励的地址......")

//./main createBlockchain -address 1AT8vEUJWbdA9J7nNSBqmxJcrYmEqMFcYY
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
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {

			log.Panic(err)
		}
	case "resetUTXO":
		err := resetUTXOCMD.Parse(os.Args[2:])
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
		fmt.Println("your CMD is wrong!!!")
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

		for index,fromAdress := range from {
			if IsValidForAdress([]byte(fromAdress)) == false || IsValidForAdress([]byte(to[index])) == false {
				fmt.Printf("地址无效......")
				printUsage()
				os.Exit(1)
			}
		}

		amount := JSONToArray(*flagAmount)
		cli.send(from,to,amount,nodeID,*flagMine)
	}

	if createBlockchainCmd.Parsed() {

		if IsValidForAdress([]byte(*flagCreateBlockchainWithAddress)) == false {
			fmt.Println("地址不能为空....")
			printUsage()
			os.Exit(1)
		}
		cli.createBlockChain(*flagCreateBlockchainWithAddress,nodeID)

	}
	if resetUTXOCMD.Parsed() {

		fmt.Println("测试....")
		cli.resetUTXOSet(nodeID)
	}
	if printChainCmd.Parsed() {
		cli.Jyh_printBlockChain(nodeID)
		fmt.Println("输出所有区块的数据........")
	}

	if listWalletCmd.Parsed() {
		cli.Jyh_listWallet(nodeID)
	}
	if createWalletCmd.Parsed() {

		cli.createwallet(nodeID)
	}

	if getBalanceCmd.Parsed() {

		if IsValidForAdress([]byte(*getBalanceWithAdress)) == false {
			fmt.Println("地址无效....")
			printUsage()
			os.Exit(1)
		}

		cli.getBalance(*getBalanceWithAdress,nodeID)
	}
	if startNodeCmd.Parsed() {


		cli.startNode(nodeID,*flagMiner)
	}

}

