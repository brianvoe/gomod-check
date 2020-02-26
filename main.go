package main

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/abiosoft/ishell"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"

	"golang.org/x/mod/modfile"
)

// FileName is the main file used for parsing
var FileName = "go.mod"

func main() {
	gofakeit.BS()

	// Read in go mod file
	dat, err := ioutil.ReadFile("./" + FileName)
	if err != nil {
		panic(err)
	}

	// Pares file
	file, err := modfile.Parse(FileName, dat, nil)
	if err != nil {
		panic(err)
	}

	// Colors
	red := color.New(color.FgHiRed).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	green := color.New(color.FgHiGreen).SprintFunc()

	var mods Mods

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

	}

	// Sort mods by status
	sort.Sort(mods)

	// Create options
	var options []string
	for _, m := range mods {
		optionStr := m.Path + " " + m.CurrentVersion.original + " " + m.AvailableVersions[0].original
		if m.Status == "major" {
			options = append(options, red(optionStr))
		} else if m.Status == "minor" {
			options = append(options, yellow(optionStr))
		} else if m.Status == "patch" {
			options = append(options, green(optionStr))
		}
	}

	if len(options) == 0 {
		fmt.Println(green("You are all up to date!!!"))
		return
	}

	shell := ishell.New()

	shell.AddCmd(&ishell.Cmd{
		Name: "checklist",
		Help: "checklist prompt",
		Func: func(c *ishell.Context) {
			choices := c.Checklist(options, "Hit space to select packages you want to update. Ctrl + c to cancel\n"+green("Patch")+" "+yellow("Minor")+" "+red("Major"), nil)

			if len(choices) > 0 && choices[0] != -1 {
				for _, i := range choices {
					c.Println(mods[i].Path, mods[i].CurrentVersion.original, mods[i].AvailableVersions[0].original, mods[i].Status)
				}
			}

			shell.Close()
		},
	})

	shell.Process("checklist")
	shell.Run()
}
