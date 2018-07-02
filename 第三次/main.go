package main

import
(
	"./BLC"
	"log"
	"os"
	"fmt"
	"flag"
)

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


func main()  {

	isValidArgs()

	addBlockCmd := flag.NewFlagSet("addBlock",flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain",flag.ExitOnError)

	flagAddBlockData := addBlockCmd.String("data","http://liyuechun.org","交易数据......")

	genesisblock := BLC.NewGenesisBlock([]byte("it's genesis block"))
	Blockchain := BLC.NewBlockChain(genesisblock)
	BLC.NewBlock([]byte("2nd"), Blockchain)

	switch os.Args[1] {
	case "addBlock":
		BLC.NewBlock([]byte("3rd"), Blockchain)
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		Blockchain.Printchain()
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
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

	defer Blockchain.BlockDB.Close()

	//	fmt.Println(genesisblock)
}


