package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

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

		fmt.Println(r.Mod.Path, " ", r.Mod.Version)
		getVersions(r.Mod.Path)
	}

}

func getVersions(url string) {
	resp, err := http.Get(fmt.Sprintf("https://proxy.golang.org/%s/@v/list", url))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
