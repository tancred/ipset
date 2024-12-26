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
	set := New()
	defer set.Close()

	set.Destroy(namedSetV4)
	set.Destroy(noSuchSet)
	set.Destroy(namedSetV6)
	set.Destroy(noSuchSet)

	set.Create(namedSetV4)
	set.Add(namedSetV4, net.IPv4(1,2,3,4))

	set.Create(namedSetV6, CreateOptionFamily("inet6"))
	set.Add(namedSetV6, net.ParseIP("::1").To16())

	return func(t *testing.T) {
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

func TestTestV6(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	name := namedSetV6

	addr := net.ParseIP("::1")
	r, err := set.Test(name, addr)
	if err != nil {
		t.Errorf("address %s: unexpected error %v", addr.String(), err)
	}
	if !r {
		t.Errorf("address %s expected in the set %s", addr.String(), name)
	}

	addr = net.ParseIP("::2")
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

func TestCreateWithFamilyInet6(t *testing.T) {
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

func TestCreateDuplicate(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	err := set.Create(namedSetV4)

	if err == nil {
		t.Fatalf("expected error on missing set %s, got nothing", namedSetV4)
	}

	if !strings.Contains(err.Error(), "set with the same name already exists") {
		t.Errorf("Expected error on existing set but got '%v'", err)
	}
}

func testAddIPv4(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	ok, err := set.Add(namedSetV4, net.IPv4(1,2,3,5))

	if err != nil {
		t.Errorf("expected no error on add, got '%v'", err)
	}

	if !ok {
		t.Errorf("expected ok")
	}
}

func testAddIPv6(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	ok, err := set.Add(namedSetV6, net.ParseIP("::1").To16())

	if err != nil {
		t.Errorf("expected no error on add, got '%v'", err)
	}

	if !ok {
		t.Errorf("expected ok")
	}
}

func TestAddNoSet(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	_, err := set.Add(noSuchSet, net.IPv4(1,2,3,4))

	if err == nil {
		t.Fatalf("expected error on missing set, got nothing")
	}

	if !strings.Contains(err.Error(), "The set with the given name does not exist") {
		t.Errorf("Expected error on missing set but got '%v'", err)
	}
}

func TestAddDuplicateV4(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	_, err := set.Add(namedSetV4, net.IPv4(1,2,3,4))

	if err == nil {
		t.Fatalf("expected error on missing set, got nothing")
	}

	if !strings.Contains(err.Error(), "Element cannot be added to the set: it's already added") {
		t.Errorf("Expected error on address present but got '%v'", err)
	}
}

func TestAddDuplicateV6(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	_, err := set.Add(namedSetV6, net.ParseIP("::1").To16())

	if err == nil {
		t.Fatalf("expected error on missing set, got nothing")
	}

	if !strings.Contains(err.Error(), "Element cannot be added to the set: it's already added") {
		t.Errorf("Expected error on address present but got '%v'", err)
	}
}

func TestAddIPv6OnIPv4(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	_, err := set.Add(namedSetV4, net.ParseIP("::2").To16())

	if err == nil {
		t.Fatalf("expected error on address family mismatch, got nothing")
	}

	if !strings.Contains(err.Error(), "Syntax error: cannot parse ::2: resolving to IPv4 address failed") {
		t.Errorf("Expected parse error on IPv6 address but got '%v'", err)
	}
}

func TestAddIPv4OnIPv6(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	_, err := set.Add(namedSetV6, net.IPv4(1,2,3,4))

	if err == nil {
		t.Fatalf("expected error on address family mismatch, got nothing")
	}

	if !strings.Contains(err.Error(), "Syntax error: cannot parse 1.2.3.4: resolving to IPv6 address failed") {
		t.Errorf("Expected parse error on IPv4 address but got '%v'", err)
	}
}

func TestAdd6(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	ok, err := set.Add6(namedSetV6, net.IPv4(1,2,3,4))

	if err != nil {
		t.Fatalf("expected no error on missing set, got nothing")
	}

	if !ok {
		t.Errorf("expected ok")
	}

	ok, err = set.Add6(namedSetV6, net.ParseIP("fe80::842f:57ff:fea2:3864"))

	if err != nil {
		t.Fatalf("expected no error on missing set, got nothing")
	}

	if !ok {
		t.Errorf("expected ok")
	}
}

func TestTest6(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	set := New()
	defer set.Close()

	_, err := set.Add6(namedSetV6, net.IPv4(1,2,3,4))
	if err != nil {
		t.Fatalf("unexpected error adding v6: %v", err)
	}

	addr := net.IPv4(1,2,3,4)
	r, err := set.Test6(namedSetV6, addr)
	if err != nil {
		t.Errorf("address %s: unexpected error %v", addr.String(), err)
	}
	if !r {
		t.Errorf("address %s expected in the set %s", addr.String(), namedSetV6)
	}

	addr = net.ParseIP("::1")
	r, err = set.Test6(namedSetV6, addr)
	if err != nil {
		t.Errorf("address %s: unexpected error %v", addr.String(), err)
	}
	if !r {
		t.Errorf("address %s expected in the set %s", addr.String(), namedSetV6)
	}

	addr = net.ParseIP("::2")
	r, err = set.Test6(namedSetV6, addr)
	if err != nil {
		t.Errorf("address %s: unexpected error %v", addr.String(), err)
	}
	if r {
		t.Errorf("address %s not expected on set %s but was", addr.String(), namedSetV6)
	}
}
