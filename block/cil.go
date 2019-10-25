package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	bc *Blockchain
}

const usage = `
Usage:
  addblock -data BLOCK_DATA    add a block to the blockchain
  printchain                   print all the blocks of the blockchain
`

func (cli *CLI) printUsage() {
	fmt.Println(usage)
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}