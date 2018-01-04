package main

import (
	"fmt"
	"os"

	"github.com/IgaguriMK/ledscreen/runes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: no arg")
	}

	r := runes.GetFirstRune(os.Args[1])
	if runes.IsWide(r) {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
