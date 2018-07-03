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
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")

}


func isValidArgs()  {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}


func (cli *CLI) Run (){
	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addBlock",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockchain",flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data","http://liyuechun.org","交易数据......")



	switch os.Args[1] {
	case "addBlock":
		Blockchain:=BlockchainObject()
		NewBlock([]byte("3rd"), Blockchain)
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		defer Blockchain.BlockDB.Close()
	case "createBlockchain":
		Blockchain:=cli.createBlockChain("put a genesis address")
		err := createBlockChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		defer Blockchain.BlockDB.Close()
	case "printchain":
		Blockchain:=BlockchainObject()
		Blockchain.Printchain()
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		defer Blockchain.BlockDB.Close()
	default:
		printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}

		fmt.Println(*flagAddBlockData)
	}

	if printChainCmd.Parsed() {

		fmt.Println("输出所有区块的数据........")

	}


	//	fmt.Println(genesisblock)
}