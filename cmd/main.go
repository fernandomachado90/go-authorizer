package main

import (
	"fmt"
	"os"
)

func main() {
	// given
	db := NewMemoryDB()
	p := Parser{
		accountManager: NewAccountManager(db),
	}

	for {
		fmt.Println(p.Parse(os.Stdin))
	}
}
