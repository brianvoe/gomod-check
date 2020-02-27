package main

import "errors"

// Mods is an alias to an array of mods
type Mods []*Mod

// Mod is the main struct containing path, current version and possible available versions
type Mod struct {
	Path              string
	CurrentVersion    *version
	AvailableVersions versions
	Status            string
}

func (v Mods) Len() int           { return len(v) }
func (v Mods) Less(i, j int) bool { return v[i].compare(v[j]) > 0 }
func (v Mods) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }

func (m *Mod) compare(m2 *Mod) int {
	if statusInt(m.Status) < statusInt(m2.Status) {
		return -1
	}
	if statusInt(m.Status) > statusInt(m2.Status) {
		return 1
	}

	return 0
}

func statusInt(status string) int {
	if status == "major" {
		return 0
	} else if status == "minor" {
		return 1
	} else if status == "patch" {
		return 2
	}

	return 0
}

// NewMod will parse the current version and get all possible versions for this package
func NewMod(path string, version string) (*Mod, error) {
	// Parse current version
	current, err := parseVersion(version)
	if err != nil {
		return nil, err
	}

	// Get all available versions
	vs := getProxyVersions(path, true)

	// Lets make sure their is a latest version to get
	if len(vs) == 0 {
		return nil, errors.New("No latest version")
	}

	latest := vs[0]

	// Check if current version is up to date
	status := ""
	compare := current.compare(latest)
	if compare >= 0 {
		status = "current"
	} else if current.major < latest.major {
		status = "major"
	} else if current.minor < latest.minor {
		status = "minor"
	} else if current.patch < latest.patch {
		status = "patch"
	}

	return &Mod{
		Path:              path,
		CurrentVersion:    current,
		AvailableVersions: vs,
		Status:            status,
	}, nil
}
