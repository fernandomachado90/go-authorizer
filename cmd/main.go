package main

import (
	"fmt"
	"os"
)

var CurrentAccount *Account

func main() {
	for {
		fmt.Println(Parse(os.Stdin))
	}
}
