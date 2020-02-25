package main

// Mod is the main struct containing path, current version and possible available versions
type Mod struct {
	Path              string
	CurrentVersion    *version
	AvailableVersions versions
	Status            string
}

// NewMod will parse the current version and get all possible versions for this package
func NewMod(path string, version string) (*Mod, error) {
	// Parse current version
	current, err := parseVersion(version)
	if err != nil {
		return nil, err
	}

	// Get all available versions
	vs := getProxyVersions(path)
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
