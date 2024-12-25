package main

import (
	"fmt"
	"net"

	"tancred/testipset/ipset"
)

func main() {
	fmt.Println("Hello, World!")

	set := ipset.New()
	defer set.Close()

	printIP(net.IPv4(10, 255, 0, 0))
	printIP(net.IP{0xfc, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	testIPv4(set, "bl", net.IPv4(1,2,3,4))
	testIPv4(set, "bl", net.IPv4(1,2,3,5))

	info, err := set.Info("bl")
	if err != nil {
		fmt.Printf("info failed: %v\n", err)
	} else {
		fmt.Printf("info: %v\n", info)
	}

	set.Fail()

	set.Test("bl6", net.ParseIP("::1").To16())
	set.Test("bl6", net.ParseIP("::2").To16())

	set.Save("bl6")

	fmt.Println("ipset", set)
}

func printIP(a net.IP) {
	fmt.Println("printIP")
	fmt.Printf("  %v %#v\n", a, a)
	fmt.Printf("  %v %#v\n", a.To4(), a.To4())
	fmt.Printf("  %v %#v\n", a.To16(), a.To16())
	fmt.Printf("  '%s' '%s'\n", a.To4().String(), a.To16().String())
}

func testIPv4(set *ipset.IPSet, name string, addr net.IP) {
	fmt.Println("testing", addr)
	ok, err := set.Test(name, addr)
	if err != nil {
		fmt.Println("  ", addr, "error", err)
	} else if ok {
		fmt.Println("  ", addr, "is ON")
	} else {
		fmt.Println("  ", addr, "is off")
	}
}
