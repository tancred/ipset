package main

import (
	"fmt"

	"tancred/ipset/ipset"
)

func main() {
	fmt.Println("Hello, World!")

	set := ipset.NewIPSet()
	defer set.Close()

	fmt.Println("ipset", set)
}