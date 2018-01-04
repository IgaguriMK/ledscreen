package main

import (
	"fmt"
	"os"

	"golang.org/x/text/width"

	"github.com/IgaguriMK/ledscreen/runes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: no arg")
	}

	r := runes.GetFirstRune(os.Args[1])
	k := width.LookupRune(r).Kind()

	switch k {
	case width.Neutral:
		fallthrough
	case width.EastAsianNarrow:
		fallthrough
	case width.EastAsianHalfwidth:
		os.Exit(0)
	default:
		os.Exit(1)
	}
}
