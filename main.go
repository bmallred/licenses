package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	var directory string
	var err error

	argc := len(os.Args)
	if argc < 2 {
		// Load the current working directory
		directory, err = os.Getwd()
		if err != nil {
			panic(err)
			return
		}
	} else {
		// Assume the first argument is the target directory
		directory = os.Args[1]
	}

	licenses := []License{}
	metrics := map[string]int{}
	files := walk(directory)
	for file := range files {
		li := parse(file)
		licenses = append(licenses, li)

		for _, t := range li.Types {
			if _, ok := metrics[t]; ok {
				metrics[t]++
			} else {
				metrics[t] = 1
			}
		}
	}

	printTitle()
	printOverview(metrics)
	printLicenses(directory, licenses...)
}

func printTitle() {
	fmt.Println("# Dependencies and licenses\n")
	fmt.Println("This is a collection of dependencies within the project and there associated licenses and copyright information. The data collected was done using [automated tools](https://github.com/bmallred/licenses) and require manual review prior to acceptance.")
	fmt.Println("")
}

func printOverview(metrics map[string]int) {
	// Do some sorting first
	keys := []string{}
	for k, _ := range metrics {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Output a table of the stats
	fmt.Println("## Overview\n")
	for _, k := range keys {
		fmt.Printf("| %s ", k)
	}
	fmt.Printf("|\n")

	for _, _ = range keys {
		fmt.Printf("| -- ")
	}
	fmt.Printf("|\n")

	for _, k := range keys {
		fmt.Printf("| %d ", metrics[k])
	}
	fmt.Printf("|\n\n")
}

func printLicenses(directory string, licenses ...License) {
	fmt.Println("## Packages\n")

	for _, li := range licenses {
		fmt.Printf("### %s\n\n", li.Package)
		fmt.Printf(" - **License**: %s\n", strings.Join(li.Types, ", "))
		if li.Version != "" {
			fmt.Printf(" - **Version**: %s\n", li.Version)
		}
		if li.Year != "" {
			fmt.Printf(" - **Year**: %s\n", li.Year)
		}
		if li.Author != "" {
			fmt.Printf(" - **Author**: %s\n", li.Author)
		}
		relativePath := strings.Replace(li.File, directory, "", 1)
		if strings.HasPrefix(relativePath, "/") {
			relativePath = strings.TrimLeft(relativePath, "/")
		}
		fmt.Printf("\n[View license](%s)\n\n", relativePath)
	}
}

func walk(directory string) chan string {
	c := make(chan string)
	re := regexp.MustCompile(`(license|licence)`)

	go func() {
		filepath.Walk(directory, func(p string, fi os.FileInfo, e error) error {
			if e != nil {
				return e
			}

			if !fi.IsDir() {
				name := strings.ToLower(fi.Name())
				badSuffix := strings.HasSuffix(name, ".after") || strings.HasSuffix(name, ".before")

				if !badSuffix && re.Match([]byte(name)) {
					c <- p
				}
			}

			return nil
		})

		defer close(c)
	}()

	return c
}
