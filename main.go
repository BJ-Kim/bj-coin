package main

import (
	"bj-coin/block"
	"bj-coin/cli"
)

func main() {
	bc := block.GetBlockChain()
	defer bc.DB.Close()

	// for i := 1; i < 10; i++ {
	// 	bc.AddBlock(strconv.Itoa(i))
	// }

	// bc.ShowBlocks()

	c := cli.NewCli(bc)
	c.Run()
}
