package ipset

import (
	"net"
	"testing"
)

// rename package iplist; we only support type hash:ip and family inet or inet6, and optionally a timeout.
// type connState int
//
//const (
//	startState connState = iota
//)

// IPList.Create(name string, )
// create and destroy a {v4,v6} set (needs "show"!)
// create a {v4,v6} set with timeout
// create a duplicate {v4,v6} set

// Create()
// show {v4,v6}, with/without timeout -> {"hash:ip", inet/inet6, -/timeout}

func TestTestV4(t *testing.T) {
	// for now we assume the set 'bl' exists
	// and that 1.2.3.4 is in the set
	set := New()
	defer set.Close()

	name := "bl"

	addr := net.IPv4(1,2,3,4)
	r, err := set.Test(name, addr)
	if err != nil {
		t.Errorf("address %s: unexpected error %v", addr.String(), err)
	}
	if !r {
		t.Errorf("address %s expected in the set %s", addr.String(), name)
	}

	addr = net.IPv4(1,2,3,5)
	r, err = set.Test(name, addr)
	if err != nil {
		t.Errorf("address %s: unexpected error %v", addr.String(), err)
	}
	if r {
		t.Errorf("address %s not expected on set %s but was", addr.String(), name)
	}
}

func TestInfo(t *testing.T) {
	set := New()
	defer set.Close()

	info, err := set.Info("bl")

	if err != nil {
		t.Fatalf("SetInfo returned error unexpectedly: %v", err)
	}

	if info.Name != "bl" {
		t.Errorf("expected name 'bl', was '%s'", info.Name)
	}
	if info.Type != "hash:ip" {
		t.Errorf("expected type 'hash:ip', was '%s'", info.Type)
	}
	if info.Family != "inet" {
		t.Errorf("expected family 'inet', was '%s'", info.Family)
	}
	if info.Timeout == nil {
		t.Errorf("expected timeout 600, was nil")
	} else if *info.Timeout != 600 {
		t.Errorf("expected timeout 600, was '%v'", *info.Timeout)
	}
}
