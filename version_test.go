package main

import (
	"sort"
	"testing"
)

func TestParseVersion(t *testing.T) {
	v, err := parseVersion("v3.2.1")
	if err != nil {
		t.Fatal(err)
	}

	if v.major != 3 {
		t.Fatal("major doesnt equal 3")
	}

	if v.minor != 2 {
		t.Fatal("minor doesnt equal 2")
	}

	if v.patch != 1 {
		t.Fatal("patch doesnt equal 1")
	}
}

func TestSortVersion(t *testing.T) {
	vs := []string{"v4.2.1", "v0.5.7", "v2.15.99"}

	var versionsArray versions
	for _, v := range vs {
		ver, err := parseVersion(v)
		if err != nil {
			t.Fatal(err)
		}
		versionsArray = append(versionsArray, ver)
	}

	sort.Sort(versionsArray)

	if versionsArray[0].original != "v0.5.7" {
		t.Fatal("0 Not in order got:", versionsArray[0].original)
	}
	if versionsArray[1].original != "v2.15.99" {
		t.Fatal("1 Not in order got:", versionsArray[1].original)
	}
	if versionsArray[2].original != "v4.2.1" {
		t.Fatal("2 Not in order got:", versionsArray[2].original)
	}
}
