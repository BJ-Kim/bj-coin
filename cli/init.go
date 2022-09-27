package cli

import (
	"bj-coin/block"
	"flag"
	"fmt"
	"log"
	"os"
)

type Cli struct {
	bc *block.BlockChain
}

func NewCli(bc *block.BlockChain) *Cli {
	return &Cli{bc}
}

func (cli *Cli) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *Cli) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *Cli) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success add block")
}

func (cli *Cli) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			fmt.Println("?D!!@?")
			fmt.Println(*addBlockData)
			fmt.Println("?D!!@?")
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.bc.ShowBlocks()
	}
}
