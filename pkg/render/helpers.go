package render

import (
	//"github.com/CharlyF/cluster-monitoring/pkg/aggregator"
	"github.com/dustin/go-humanize"
	"golang.org/x/text/unicode/norm"
	"html/template"
	"strings"
)

// Fmap return a fresh copy of a map of utility functions for templating
func Fmap() template.FuncMap {
	return template.FuncMap{
		"printDashes":        printDashes,
		"humanize": 			mkHuman,
	}
}

func printDashes(s string, dash string) string {
	return strings.Repeat(dash, stringLength(s))
}

func stringLength(s string) int {
	/*
		len(string) is wrong if the string has unicode characters in it,
		for example, something like 'Agent (v6.0.0+Χελωνη)' has len(s) == 27.
		This is a better way of counting a string length
		(credit goes to https://stackoverflow.com/a/12668840)
	*/
	var ia norm.Iter
	ia.InitString(norm.NFKD, s)
	nc := 0
	for !ia.Done() {
		nc = nc + 1
		ia.Next()
	}
	return nc
}

// mkHuman makes large numbers more readable
func mkHuman(f float64) string {
	var str string
	if f > 1000000.0 {
		str = humanize.SIWithDigits(f, 1, "")
	} else {
		str = humanize.Commaf(f)
	}

	return str
}
/*
----------------
src      | dest
----------------
tag1      | tag2
tag2      | tag4
----------------
 */
