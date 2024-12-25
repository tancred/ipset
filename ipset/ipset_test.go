package ipset

import (
	"net"
	"strings"
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

// add, no set
// add, duplicate
// add, ipv6 on ipv4
// add, ipv4 on ipv6

const (
	namedSetV4 = "bl4"
	namedSetV6 = "bl6"
	noSuchSet = "bl2"
)

func setup(t *testing.T) func(t *testing.T) {
	t.Log("setup")
	set := New()
	defer set.Close()

	set.Destroy(namedSetV4)
	set.Destroy(noSuchSet)
	set.Destroy(namedSetV6)
	set.Destroy(noSuchSet)

	set.Create(namedSetV4)
	set.Add(namedSetV4, net.IPv4(1,2,3,4))

	return func(t *testing.T) {
		t.Log("teardown")
		set := New()
		defer set.Close()

		set.Destroy(namedSetV4)
		set.Destroy(noSuchSet)
		set.Destroy(namedSetV6)
		set.Destroy(noSuchSet)
	}
}

func TestTestV4(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	// for now we assume the set 'bl' exists
	// and that 1.2.3.4 is in the set
	set := New()
	defer set.Close()

	name := namedSetV4

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

func TestInfoOnNonexistent(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	info, err := set.Info(noSuchSet)

	if err == nil {
		t.Fatalf("expected error on missing set %s, instead got info %v", noSuchSet, info)
	}

	if !strings.Contains(err.Error(), "The set with the given name does not exist") {
		t.Errorf("Expected error on missing set but got '%v'", err)
	}
}

func TestCreateDefault(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	err := set.Create(noSuchSet)

	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	info, err := set.Info(noSuchSet)
	if err != nil {
		t.Fatalf("expected set '%s', got error: %v", noSuchSet, err)
	}

	if info.Name != noSuchSet {
		t.Errorf("expected name '%s', was '%s'", noSuchSet, info.Name)
	}
	if info.Type != "hash:ip" {
		t.Errorf("expected type 'hash:ip', was '%s'", info.Type)
	}
	if info.Family != "inet" {
		t.Errorf("expected family 'inet', was '%s'", info.Family)
	}
	if info.Timeout != nil {
		t.Errorf("expected no timeout, was '%v'", *info.Timeout)
	}
}

func TestCreateWithTimeout(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	err := set.Create(noSuchSet, CreateOptionTimeout(601))

	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	info, err := set.Info(noSuchSet)
	if err != nil {
		t.Fatalf("expected set '%s', got error: %v", noSuchSet, err)
	}

	if info.Name != noSuchSet {
		t.Errorf("expected name '%s', was '%s'", noSuchSet, info.Name)
	}
	if info.Type != "hash:ip" {
		t.Errorf("expected type 'hash:ip', was '%s'", info.Type)
	}
	if info.Family != "inet" {
		t.Errorf("expected family 'inet', was '%s'", info.Family)
	}
	expectedTimeout := 601
	if info.Timeout == nil {
		t.Errorf("expected timeout %d, was nil", expectedTimeout)
	} else if *info.Timeout != expectedTimeout {
		t.Errorf("expected timeout %d, was '%v'", expectedTimeout, *info.Timeout)
	}
}

func TestCreateV6(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	err := set.Create(noSuchSet, CreateOptionFamily("inet6"))

	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	info, err := set.Info(noSuchSet)
	if err != nil {
		t.Fatalf("expected set '%s', got error: %v", noSuchSet, err)
	}

	if info.Name != noSuchSet {
		t.Errorf("expected name '%s', was '%s'", noSuchSet, info.Name)
	}
	if info.Type != "hash:ip" {
		t.Errorf("expected type 'hash:ip', was '%s'", info.Type)
	}
	if info.Family != "inet6" {
		t.Errorf("expected family 'inet6', was '%s'", info.Family)
	}
	if info.Timeout != nil {
		t.Errorf("expected no timeout, was '%v'", *info.Timeout)
	}
}
