package main

import (
	"fmt"
	"io/ioutil"

	"github.com/brianvoe/gofakeit/v4"
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
	}

}
