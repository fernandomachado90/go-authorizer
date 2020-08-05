package main

import (
	"fmt"
	"os"
)

func initHandler() Handler {
	db := NewMemoryDB()
	return Handler{
		db:             db,
		accountHandler: NewAccountManager(db),
	}
}

func main() {
	h := initHandler()
	for {
		fmt.Println(h.Encode(h.Dispatch(h.Decode(os.Stdin))))
	}
}
