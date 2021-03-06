package main

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/abiosoft/ishell"
	"github.com/fatih/color"

	"golang.org/x/mod/modfile"
)

// FileName is the main file used for parsing
var FileName = "go.mod"

func main() {
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
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	var mods Mods
	urlLength := 0
	curVersionLength := 0
	latestVersionLength := 0

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

		// Check url length for option padding
		if urlLength < len(mod.Path) {
			urlLength = len(mod.Path)
		}
		if curVersionLength < len(mod.CurrentVersion.cleanString()) {
			curVersionLength = len(mod.CurrentVersion.cleanString())
		}
		if latestVersionLength < len(mod.AvailableVersions[0].cleanString()) {
			latestVersionLength = len(mod.AvailableVersions[0].cleanString())
		}

		mods = append(mods, mod)
	}

	// Sort mods by status
	sort.Sort(mods)

	// Create options
	var options []string
	hasIncompatible := false
	for _, m := range mods {
		urlStr := strPadding(m.Path, urlLength) + "   "
		versionStr := strPadding(m.CurrentVersion.cleanString(), curVersionLength) + " -> " + strPadding(m.AvailableVersions[0].cleanString(), latestVersionLength)
		if m.AvailableVersions[0].incompatible {
			versionStr += " " + blue("I")
			hasIncompatible = true
		}
		if m.Status == "major" {
			options = append(options, red(urlStr)+versionStr)
		} else if m.Status == "minor" {
			options = append(options, yellow(urlStr)+versionStr)
		} else if m.Status == "patch" {
			options = append(options, green(urlStr)+versionStr)
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
			text := "Hit space to select modules you want to update. Ctrl + c to cancel\n"
			text += green("Patch") + " " + yellow("Minor") + " " + red("Major") + " "
			if hasIncompatible {
				text += blue("I = Incompatible")
			}
			choices := c.Checklist(options, text, nil)

			if len(choices) > 0 && choices[0] != -1 {
				c.ClearScreen()
				c.Println(green("Modules that were updated!!!"))
				for _, i := range choices {
					err := file.AddRequire(mods[i].Path, mods[i].AvailableVersions[0].original)
					if err != nil {
						c.Err(err)
					}
					c.Println(options[i])
				}
				file.Cleanup()
				dat, err := file.Format()
				if err != nil {
					c.Err(err)
				}

				// Write back to file
				err = ioutil.WriteFile("./"+FileName, dat, 644)
				if err != nil {
					c.Err(err)
				}
			}

			shell.Close()
		},
	})

	shell.Process("checklist")
	shell.Run()
}
