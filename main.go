package main

import (
	"os"
	"fmt"
	"cli"
)

func main() {
	usage := `Chaingo - an blockchain write by golang

Usage:
	chaingo [command] options

Commands:
	mine	start mining
	account	account manage
	version	show chaingo version
	help	shows a list of commands or help for one command.

Options:
	-h 	show command help

Author:
   	cheuk <zhuohong@live.com>`

	if len(os.Args) == 1 {
		fmt.Println(usage)
		return
	}
	if os.Args[1] == "mine" {

	}
	if os.Args[1] == "account" {
		accountUsage := `Usage:
  chaingo account [new|list|balance]`
		if len(os.Args) == 2 || os.Args[1] == "-h" {
			fmt.Println(accountUsage)
			return
		}
		if os.Args[2] == "list" {
			cli.ListAccount()
			return
		}
		if os.Args[2] == "new" {
			cli.NewAccount()
			return
		}
		if os.Args[2] == "balance" {
			if len(os.Args) == 3 {
				fmt.Println("ERROR: please enter an address")
				return
			}
			cli.Balance(os.Args[3])
			return
		}
	}
}
