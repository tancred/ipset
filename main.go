package main

import (
	"fmt"
	"net"

	"tancred/testipset/ipset"
)

func main() {
	fmt.Println("Hello, World!")

	set := ipset.NewIPSet()
	defer set.Close()

	a := net.IPv4(10, 255, 0, 0)
	// a := net.IP{0xfc, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	// a := net.ParseIP("::1")
	fmt.Printf("%v %#v\n", a, a)
	fmt.Printf("%v %#v\n", a.To4(), a.To4())
	fmt.Printf("%v %#v\n", a.To16(), a.To16())

	set.Test("bl", net.IPv4(1,2,3,5))
	set.Test("bl", net.IPv4(1,2,3,4))

	set.Save("bl")

	set.Fail()

	fmt.Println("ipset", set)
}