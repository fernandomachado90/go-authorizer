package main

import (
	"fmt"
	"os"
)

func main() {
	db := NewMemoryDB()
	p := Parser{
		accountManager: NewAccountManager(db),
	}

	for {
		fmt.Println(p.Parse(os.Stdin))
	}
}
