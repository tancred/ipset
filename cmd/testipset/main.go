package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/tancred/ipset"
)

func main() {
	set := ipset.New()
	defer set.Close()

	timeout := 604800

	createSetIfNecessary(set, ipset.Info{Name: "bl", Type: "hash:ip", Family: "inet", Timeout: &timeout})
	set.Add("bl", net.IPv4(1, 2, 3, 5))

	testIPv4(set, "bl", net.IPv4(1, 2, 3, 4))
	testIPv4(set, "bl", net.IPv4(1, 2, 3, 5))

	info, err := set.Info("bl")
	if err != nil {
		fmt.Printf("info failed: %v\n", err)
	} else {
		fmt.Printf("info: %v\n", info)
	}

	createSetIfNecessary(set, ipset.Info{Name: "bl6", Type: "hash:ip", Family: "inet6", Timeout: &timeout})
	set.Add6("bl6", net.ParseIP("fe80::842f:57ff:fea2:3864"))
	testIPv6(set, "bl6", net.ParseIP("::1"))
	testIPv6(set, "bl6", net.ParseIP("fe80::842f:57ff:fea2:3864"))
}

func createSetIfNecessary(set *ipset.IPSet, expInfo ipset.Info) {
	actInfo, err := set.Info(expInfo.Name)

	if err != nil {
		if errors.Is(err, ipset.ErrSetNotFound) {
			log.Printf("creating set %s", expInfo.Name)
			createSet(set, expInfo)
			return
		}
		log.Fatalf("error: %v", err)
	} else if !checkInfo(expInfo, actInfo) {
		log.Printf("destroying set %s", expInfo.Name)

		err = set.Destroy(expInfo.Name)
		if err != nil {
			log.Fatalf("failed to remove set %s", expInfo.Name)
		}

		log.Printf("recreating set %s", expInfo.Name)
		createSet(set, expInfo)
	}
}

func createSet(set *ipset.IPSet, expInfo ipset.Info) {
	var opts []ipset.CreateOption

	opts = append(opts, ipset.CreateOptionFamily(expInfo.Family))

	if expInfo.Timeout != nil {
		opts = append(opts, ipset.CreateOptionTimeout(*expInfo.Timeout))
	}

	err := set.Create(expInfo.Name, opts...)

	if err != nil {
		log.Fatalf("can't create ipset '%s': %v", expInfo.Name, err)
	}

	return
}

func checkInfo(expInfo ipset.Info, actInfo ipset.Info) bool {
	res := true

	if actInfo.Type != expInfo.Type {
		log.Printf("Set %s has wrong type, expected %s but was %s", expInfo.Name, expInfo.Type, actInfo.Type)
		res = false
	}

	if actInfo.Family != expInfo.Family {
		log.Printf("Set %s has wrong family, expected %s but was %s", expInfo.Name, expInfo.Family, actInfo.Family)
		res = false
	}

	eqTimeout := func(a ipset.Info, b ipset.Info) bool {
		return (actInfo.Timeout == nil && expInfo.Timeout == nil) || (actInfo.Timeout != nil && expInfo.Timeout != nil && *actInfo.Timeout == *expInfo.Timeout)
	}

	prntTimeout := func(t *int) string {
		if t == nil {
			return "<nil>"
		}
		return fmt.Sprintf("%d", *t)
	}

	if !eqTimeout(actInfo, expInfo) {
		log.Printf("Set %s has wrong timeout, expected %v but was %v", actInfo.Name, prntTimeout(expInfo.Timeout), prntTimeout(actInfo.Timeout))
		res = false
	}

	if res {
		log.Printf("set %s is present", expInfo.Name)
	}

	return res
}

func printIP(a net.IP) {
	fmt.Println("printIP")
	fmt.Printf("a        %v %#v\n", a, a)
	fmt.Printf("a.To4()  %v %#v\n", a.To4(), a.To4())
	fmt.Printf("a.To16() %v %#v\n", a.To16(), a.To16())
	fmt.Printf("a str   '%s' '%s'\n", a.To4().String(), a.To16().String())
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

func testIPv6(set *ipset.IPSet, name string, addr net.IP) {
	fmt.Println("testing", addr.To16())
	ok, err := set.Test(name, addr.To16())
	if err != nil {
		fmt.Println("  ", addr.To16(), "error", err)
	} else if ok {
		fmt.Println("  ", addr.To16(), "is ON")
	} else {
		fmt.Println("  ", addr.To16(), "is off")
	}
}
