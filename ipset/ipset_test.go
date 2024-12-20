package ipset

import (
	"testing"
)

func Test(t *testing.T) {
	set := NewIPSet()
	defer func () {
		set.Destroy("thelist")
		set.Close()
	}()

	t.Fatalf("Oh, no!")
}
