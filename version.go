package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var allowed string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"
var versionRegex = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

type versions []*version

func (v versions) Len() int           { return len(v) }
func (v versions) Less(i, j int) bool { return v[i].compare(v[j]) < 0 }
func (v versions) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

// version represents a single semantic version.
type version struct {
	major, minor, patch uint64
	prerelease          string
	metadata            string
	original            string
}

// parseVersion parses a given version string
func parseVersion(v string) (*version, error) {
	m := versionRegex.FindStringSubmatch(v)
	if m == nil {
		return nil, errors.New("Invalid Semantic Version")
	}

	var err error

	// Regex breakdown
	major := m[1]
	minor := m[2]
	patch := m[3]
	prerelease := m[5]
	metadata := m[8]

	sv := &version{
		prerelease: prerelease,
		metadata:   metadata,
		original:   v,
	}

	// Major check
	sv.major, err = strconv.ParseUint(major, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing version segment: %s", err)
	}

	// Minor check
	if minor != "" {
		sv.minor, err = strconv.ParseUint(strings.TrimPrefix(minor, "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version segment: %s", err)
		}
	} else {
		sv.minor = 0
	}

	// Patch check
	if patch != "" {
		sv.patch, err = strconv.ParseUint(strings.TrimPrefix(patch, "."), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version segment: %s", err)
		}
	} else {
		sv.patch = 0
	}

	// Prerelease check
	if sv.prerelease != "" {
		if err = validatePrerelease(sv.prerelease); err != nil {
			return nil, err
		}
	}

	// Metadata check
	if sv.metadata != "" {
		if err = validateMetadata(sv.metadata); err != nil {
			return nil, err
		}
	}

	return sv, nil
}

func (v *version) compare(v2 *version) int {
	if d := compareSegment(v.major, v2.major); d != 0 {
		return d
	}
	if d := compareSegment(v.minor, v2.minor); d != 0 {
		return d
	}
	if d := compareSegment(v.patch, v2.patch); d != 0 {
		return d
	}

	// Major, minor, and patch are the same lets check prerelease
	if v.prerelease == "" && v2.prerelease == "" {
		return 0
	}
	if v.prerelease == "" {
		return 1
	}
	if v2.prerelease == "" {
		return -1
	}

	return 0
}

func compareSegment(v, o uint64) int {
	if v < o {
		return -1
	}
	if v > o {
		return 1
	}

	return 0
}

// Like strings.ContainsAny but does an only instead of any.
func containsOnly(s string, comp string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(comp, r)
	}) == -1
}

// validatePrerelease loops through values to check for valid characters
func validatePrerelease(p string) error {
	eparts := strings.Split(p, ".")
	for _, p := range eparts {
		if containsOnly(p, "0123456789") {
			if len(p) > 1 && p[0] == '0' {
				return errors.New("Version segment starts with 0")
			}
		} else if !containsOnly(p, allowed) {
			return errors.New("Invalid Prerelease string")
		}
	}
	return nil
}

// validateMetadata loops through values to check for valid characters
func validateMetadata(m string) error {
	eparts := strings.Split(m, ".")
	for _, p := range eparts {
		if !containsOnly(p, allowed) {
			return errors.New("Invalid Metadata string")
		}
	}
	return nil
}
