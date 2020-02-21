package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/brianvoe/gofakeit"
	"golang.org/x/mod/modfile"
)

// https://proxy.golang.org/github.com/brianvoe/gofakeit/@v/list

var fileName = "go.mod"

func main() {
	gofakeit.UUID()

	// Read in go mod file
	dat, err := ioutil.ReadFile("./" + fileName)
	if err != nil {
		panic(err)
	}

	// Pares file
	file, err := modfile.Parse(fileName, dat, nil)
	if err != nil {
		panic(err)
	}

	// Loop through required mods and exclude indirect ones
	for _, r := range file.Require {
		if r.Indirect {
			continue
		}

		// Check if mod is on the most current version

		fmt.Println(r.Mod.Path, " ", r.Mod.Version)
		vs := getVersions(r.Mod.Path)
		for _, v := range vs {
			fmt.Println(v.original)
		}
	}

}

func getVersions(url string) versions {
	resp, err := http.Get(fmt.Sprintf("https://proxy.golang.org/%s/@v/list", url))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var vs versions
	versionsString := strings.Split(string(body), "\n")
	for _, vss := range versionsString {
		v, err := parseVersion(vss)
		if err != nil {
			// If has error parsing skip it
			continue
		}
		vs = append(vs, v)
	}
	sort.Sort(vs)

	return vs
}
