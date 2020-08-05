package main

import (
	"fmt"
	"os"
)

func initHandler() Handler {
	db := NewMemoryDB()
	return Handler{
		accountManager: NewAccountManager(db),
	}
}

func main() {
	h := initHandler()
	for {
		fmt.Println(h.Encode(h.Process(h.Decode(os.Stdin))))
	}
}
