package main

import "testing"

func TestReplace(t *testing.T) {
	vals := []string{"0", "1", "2", "3", "4"}
	one := 1
	two := 3
	dst := make([]string, len(vals)-1)
	var i int
	for j, s := range vals {
		if j != one && j != two {
			dst[i] = s
			i++
		}
	}
	for i, s := range []string{"0", "2", "4"} {
		if dst[i] != s {
			t.Errorf("%q should match %q", dst[i], s)
		}
	}
}
