package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/IgaguriMK/ledscreen/colors"
	"github.com/IgaguriMK/ledscreen/pixcels"
	"github.com/IgaguriMK/ledscreen/runes"
)

func main() {
	var dotFile string
	flag.StringVar(&dotFile, "d", "image/dot.png", "Dot image file")
	var outFile string
	flag.StringVar(&outFile, "o", "output", "Output name")
	var colorFile string
	flag.StringVar(&colorFile, "c", "colors.yml", "Colormap file")
	var width int
	flag.IntVar(&width, "w", 256, "Image maximum width in dot. 0 is no scroll")

	flag.Parse()
	args := flag.Args()

	dot, err := pixcels.Load(dotFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	colors, err := colors.LoadTable(colorFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	printer := NewPrinter()
	printer.Fill = colors[":"]
	printer.Background = colors["_"]

	for _, arg := range args {
		switch {
		case arg == ":":
			printer.Background = colors[":"]
		case strings.HasPrefix(arg, ":"):
			cn := strings.TrimPrefix(arg, ":")
			col, ok := colors[cn]
			if !ok {
				fmt.Fprintf(os.Stderr, "Cannot find color '%s'\n", cn)
				os.Exit(1)
			}
			printer.Fill = col
		case arg == "_":
			printer.Background = colors["_"]
		case strings.HasPrefix(arg, "_"):
			cn := strings.TrimPrefix(arg, "_")
			col, ok := colors[cn]
			if !ok {
				fmt.Fprintf(os.Stderr, "Cannot find color '%s'\n", cn)
				os.Exit(1)
			}
			printer.Background = col
		default:
			printer.Print(arg)
		}
	}

	img := printer.Image
	switch {
	case width == 0:
		result := img.DotImage(dot)
		result.SaveTo(outFile + ".png")
	case img.Width <= width:
		left := (width - img.Width) / 2
		result := img.Cut(-left, width, colors["_"])
		result = result.DotImage(dot)
		result.SaveTo(outFile + ".png")
	default:
		for i := 0; i < img.Width+width; i++ {
			result := img.Cut(i-width, width, colors["_"])
			result = result.DotImage(dot)
			cnt := fmt.Sprintf("%06d", i)
			result.SaveTo(outFile + "/" + cnt + ".png")
		}
	}
}

type Printer struct {
	Background pixcels.Pixcel
	Fill       pixcels.Pixcel
	Image      *pixcels.PixcelArray
}

func NewPrinter() *Printer {
	return &Printer{
		Image: pixcels.NewPixcelArray(0, 16),
	}
}

func (pr *Printer) Print(str string) {
	for _, r := range str {
		imageName := runes.ImageName(r)
		runeImg, err := pixcels.Load(imageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Rune '%s' not found\n", string(r))
			os.Exit(1)
		}

		runeImg = runeImg.MultColor(pr.Fill)
		runeImg = runeImg.BackColor(pr.Background)
		pr.Image.JoinHorizontal(runeImg)

		fmt.Printf("'%v' %s\n", string(r), runes.RuneCode(r))
	}
}
