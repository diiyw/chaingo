package main

import (
	"os"
	"cli"
	"log"
	"strconv"
	"fmt"
)

const VERSION = "0.0.1"

func main() {
	usage := `Chaingo - an blockchain write by golang

Usage:
	chaingo [command] options

Commands:
	mine	start mining
	account	account manage
	chain	print blockchain
	version	show chaingo version
	help	shows this help.

Options:
	-h 	show command help

Author:
   	cheuk <zhuohong@live.com>`

	if len(os.Args) == 1 {
		fmt.Println(usage)
		return
	}
	if os.Args[1] == "mine" {
		var (
			mineUsage = `Usage:
  chaingo mine [-ip 0.0.0.0 -port 2048 -address]`
			ip   = "0.0.0.0"
			port = 2048
		)
		if len(os.Args) < 8 || os.Args[1] == "-h" {
			fmt.Println(mineUsage)
			return
		}
		if os.Args[2] == "-ip" {
			if len(os.Args) == 3 {
				log.Println("ERROR: unspecified ip")
				return
			}
			ip = os.Args[3]
		}
		if os.Args[4] == "-port" {
			if len(os.Args) == 5 {
				log.Println("ERROR: unspecified ip")
				return
			}
			p, err := strconv.ParseInt(os.Args[5], 10, 32)
			if err != nil {
				log.Println("ERROR: error port")
			}
			port = int(p)
		}
		if os.Args[6] == "-address" {
			if len(os.Args) == 7 {
				log.Println("ERROR: unspecified address")
				return
			}
			cli.StartMine(ip, port, os.Args[7])
			return
		}
		fmt.Println(mineUsage)
	}
	if os.Args[1] == "account" {
		accountUsage := `Usage:
  chaingo account [-new|-list|-fund|-send]`
		if len(os.Args) == 2 || os.Args[1] == "-h" {
			fmt.Println(accountUsage)
			return
		}
		if os.Args[2] == "-list" {
			cli.ListAccount()
			return
		}
		if os.Args[2] == "-new" {
			cli.NewAccount()
			return
		}
		if os.Args[2] == "-fund" {
			if len(os.Args) == 3 {
				log.Println("ERROR: please enter an address")
				return
			}
			cli.Fund(os.Args[3])
			return
		}
		if os.Args[2] == "-send" {
			if len(os.Args) <= 5 {
				fmt.Println(`Usage:
  chaingo account -send (amount) (from) (to)`)
				return
			}
			p, err := strconv.ParseInt(os.Args[3], 10, 32)
			if err != nil {
				log.Println("ERROR: error amount")
			}
			cli.Send(int(p), os.Args[4], os.Args[5])
			return
		}
		fmt.Println(accountUsage)
	}
	if os.Args[1] == "version" {
		fmt.Println("Version:", VERSION)
		return
	}
	if os.Args[1] == "chain" {
		cli.PrintChain()
		return
	}
	if os.Args[1] == "help" {
		fmt.Println(usage)
		return
	}
	fmt.Println(usage)
}
