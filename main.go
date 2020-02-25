package main

import (
	"fmt"
	"io/ioutil"

	"github.com/abiosoft/ishell"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"

	"golang.org/x/mod/modfile"
)

var fileName = "go.mod"

func main() {
	gofakeit.BS()

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

	// Colors
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	var mods []*Mod
	var options []string

	// Loop through required mods and exclude indirect ones
	for _, r := range file.Require {
		if r.Indirect {
			continue
		}

		// Check if mod is on the most current version
		mod, err := NewMod(r.Mod.Path, r.Mod.Version)
		if err != nil {
			continue
		}

		// Do nothing if mod is current
		if mod.Status == "current" {
			continue
		}

		mods = append(mods, mod)
		optionStr := mod.Path + " " + mod.CurrentVersion.original + " " + mod.AvailableVersions[0].original
		if mod.Status == "major" {
			options = append(options, red(optionStr))
		} else if mod.Status == "minor" {
			options = append(options, yellow(optionStr))
		} else if mod.Status == "patch" {
			options = append(options, blue(optionStr))
		}
	}

	if len(options) == 0 {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Println(green("You are all up to date!!!"))
		return
	}

	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "checklist",
		Help: "checklist prompt",
		Func: func(c *ishell.Context) {
			choices := c.Checklist(options, "What are your favourite programming languages ?", nil)

			for _, i := range choices {
				c.Println(mods[i].Path, mods[i].CurrentVersion.original, mods[i].AvailableVersions[0].original, mods[i].Status)
			}
		},
	})

	shell.Process("checklist")

	// run shell
	shell.Run()

}
