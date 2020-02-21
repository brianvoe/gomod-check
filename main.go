package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"

	"golang.org/x/mod/modfile"
)

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

		// fmt.Println(r.Mod.Path, " ", r.Mod.Version)
		// vs := getVersions(r.Mod.Path)
		// for _, v := range vs {
		// 	fmt.Println(v.original)
		// }
	}

	shell := ishell.New()

	// Colors
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	boldRed := color.New(color.FgRed, color.Bold).SprintFunc()

	shell.AddCmd(&ishell.Cmd{
		Name: "checklist",
		Help: "checklist prompt",
		Func: func(c *ishell.Context) {
			c.Print(cyan("cyan\n"))
			c.Println(yellow("yellow"))
			c.Printf("%s\n", boldRed("bold red"))

			languages := []string{cyan("Python"), "Go", "Haskell", "Rust"}
			choices := c.Checklist(languages, "What are your favourite programming languages ?", nil)
			out := func() []string {
				var checked []string
				for _, i := range choices {
					checked = append(checked, languages[i])
				}
				return checked
			}
			c.Println("Your choices are", strings.Join(out(), ", "))
		},
	})

	shell.Process("checklist")

	// run shell
	shell.Run()

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
