package main

import (
	"fmt"
	"os"

	"github.com/IgaguriMK/ledscreen/pixcels"
	"github.com/IgaguriMK/ledscreen/runes"
)

func main() {
	dot, err := pixcels.Load("image/dot.png")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "No args")
		os.Exit(1)
	}
	r := runes.GetFirstRune(os.Args[1])
	imageName := runes.ImageName(r)
	runeImg, err := pixcels.Load(imageName)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	result := runeImg.DotImage(dot)

	result.SaveTo("view.png")
}
