package main

import (
	"fmt"
	"os"
)

func main() {
	for {
		fmt.Println(Parse(os.Stdin))
	}
}
