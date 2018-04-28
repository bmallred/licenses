package main

import (
	"io/ioutil"
	"path"
	"regexp"
	"sort"
)

// License store the metadata which may be found within different types of legalese.
type License struct {
	File    string
	Package string
	Types   []string
	Version string
	Year    string
	Author  string
}

// Parse a license file extracting as much information as we can.
func parse(filepath string) License {
	license := License{
		File:  filepath,
		Types: []string{"Unknown"},
	}

	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return license
	}

	// Extract the package name... completely not correct since
	// we are just taking the base directory name.
	license.Package = path.Base(path.Dir(filepath))

	// Search for license version information
	if rxVersion.Match(contents) {
		submatches := rxVersion.FindSubmatch(contents)
		license.Version = string(submatches[1])
	}

	// Search for copyright information
	if rxCopyright.Match(contents) {
		submatches := rxCopyright.FindSubmatch(contents)
		license.Year = string(submatches[1])
		license.Author = string(submatches[2])
	}

	// Search for the type of license
	hits := []string{}
	for name, expr := range expressions {
		for _, re := range expr {
			if re.Match(contents) {
				add := true
				for _, h := range hits {
					if h == name {
						add = false
					}
				}
				if add {
					hits = append(hits, name)
				}
			}
		}
	}

	if len(hits) > 0 {
		sort.Strings(hits)
		license.Types = hits
	}

	return license
}

// Simple regular expression to try and pull some version information.
var rxVersion = regexp.MustCompile(`Version (\d+)`)

// Simple regular expression to try and pull some copyright information.
var rxCopyright = regexp.MustCompile(`Copyright.*(\d{4}),?\s([\w -!]*\w)`)

// List of expressions used to find types of licenses.
// Most of these are straight forward in the real world people dork with the language
// of the licenses and remove or alter pieces.
//
// Reference: https://choosealicense.com/appendix
var expressions = map[string][]*regexp.Regexp{
	"MIT": []*regexp.Regexp{
		regexp.MustCompile(`MIT License`),
		regexp.MustCompile(`MIT\/X11 License`),
		regexp.MustCompile(`Permission\s+is\s+hereby\s+granted,\s+free\s+of\s+charge,\s+to\s+any\s+person\s+obtaining\s+a\s+copy\s+of\s+this\s+software\s+and\s+associated\s+documentation\s+files\s+\(the\s+"Software"\),\s+to\s+deal\s+in\s+the\s+Software\s+without\s+restriction,\s+including\s+without\s+limitation\s+the\s+rights\s+to\s+use,\s+copy,\s+modify,\s+merge,\s+publish,\s+distribute,\s+sublicense,\s+and\/or\s+sell\s+copies\s+of\s+the\s+Software,\s+and\s+to\s+permit\s+persons\s+to\s+whom\s+the\s+Software\s+is\s+furnished\s+to\s+do\s+so,\s+subject\s+to\s+the\s+following\s+conditions`),
	},
	"AGPL": []*regexp.Regexp{
		regexp.MustCompile(`GNU AFFERO GENERAL PUBLIC LICENSE`),
	},
	"LGPL": []*regexp.Regexp{
		regexp.MustCompile(`GNU LESSER GENERAL PUBLIC LICENSE`),
	},
	"GPL": []*regexp.Regexp{
		regexp.MustCompile(`GNU GENERAL PUBLIC LICENSE`),
	},
	"MPL": []*regexp.Regexp{
		regexp.MustCompile(`Mozilla Public License`),
	},
	"Apache": []*regexp.Regexp{
		regexp.MustCompile(`Apache License`),
		// regexp.MustCompile(`https?:\/\/www\.apache\.org\/licenses`),
	},
	"Unlicense": []*regexp.Regexp{
		regexp.MustCompile(`free and unencumbered software`),
		regexp.MustCompile(`https?:\/\/unlicense\.org`),
	},
	"AFL": []*regexp.Regexp{
		regexp.MustCompile(`Academic Free License`),
	},
	"Artistic": []*regexp.Regexp{
		regexp.MustCompile(`Artistic License`),
	},
	"BSD": []*regexp.Regexp{
		regexp.MustCompile(`BSD 2\-Clause License`),
		regexp.MustCompile(`The Clear BSD License`),
		regexp.MustCompile(`BSD 3-Clause License`),
		regexp.MustCompile(`BSD`),
		regexp.MustCompile(`Redistribution\s+and\s+use\s+in\s+source\s+and\s+binary\s+forms,\s+with\s+or\s+without\s+modification,\s+are\s+permitted\s+provided\s+that\s+the\s+following\s+conditions\s+are\s+met`),
		regexp.MustCompile(`Redistribution\s+and\s+use\s+of\s+this\s+software\s+in\s+source\s+and\s+binary\s+forms,\s+with\s+or\s+without\s+modification,\s+are\s+permitted\s+provided\s+that\s+the\s+following\s+conditions\s+are\s+met`),
	},
	"Boost": []*regexp.Regexp{
		regexp.MustCompile(`Boost Software License`),
	},
	"CC Attribution": []*regexp.Regexp{
		regexp.MustCompile(`Attribution \d+\.\d+ International`),
	},
	"CC Attribution-ShareAlike": []*regexp.Regexp{
		regexp.MustCompile(`Attribution\-ShareAlike \d+\.\d+ International`),
	},
	"CC0": []*regexp.Regexp{
		regexp.MustCompile(`CC0 \d+\.\d+ Universal`),
	},
	"Educational": []*regexp.Regexp{
		regexp.MustCompile(`Educational Community License`),
	},
	"Eclipse": []*regexp.Regexp{
		regexp.MustCompile(`Eclipse Public License`),
	},
	"European Union": []*regexp.Regexp{
		regexp.MustCompile(`European Union Public Licen[sc]e`),
	},
	"ISC": []*regexp.Regexp{
		regexp.MustCompile(`ISC License`),
		regexp.MustCompile(`Permission\s+to\s+use,\s+copy,\s+modify,\s+and\/or\s+distribute\s+this\s+software\s+for\s+any\s+purpose\s+with\s+or\s+without\s+fee\s+is\s+hereby\s+granted,\s+provided\s+that\s+the\s+above\s+copyright\s+notice\s+and\s+this\s+permission\s+notice\s+appear\s+in\s+all\s+copies.`),
	},
	"LaTeX": []*regexp.Regexp{
		regexp.MustCompile(`The LaTeX Project Public License`),
	},
	"Ms-PL": []*regexp.Regexp{
		regexp.MustCompile(`Microsoft Public License`),
	},
	"Ms-RL": []*regexp.Regexp{
		regexp.MustCompile(`Microsoft Reciprocal License`),
	},
	"OSL": []*regexp.Regexp{
		regexp.MustCompile(`Open Software License`),
	},
	"PostgreSQL": []*regexp.Regexp{
		regexp.MustCompile(`PosgreSQL Licen[cs]e`),
	},
	"SIL OFL": []*regexp.Regexp{
		regexp.MustCompile(`SIL Open Font License`),
	},
	"UIUC/NCSA": []*regexp.Regexp{
		regexp.MustCompile(`University of Illinois\/NCSA Open Sourse License`),
		regexp.MustCompile(`University of Illinois Open Sourse License`),
		regexp.MustCompile(`NCSA Open Sourse License`),
	},
	"UPL": []*regexp.Regexp{
		regexp.MustCompile(`Universal Permissive License`),
	},
	"WTF": []*regexp.Regexp{
		regexp.MustCompile(`DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE`),
		regexp.MustCompile(`DO WTF YOU WANT TO PUBLIC LICENSE`),
	},
	"zlib": []*regexp.Regexp{
		regexp.MustCompile(`zlib License`),
	},
}
