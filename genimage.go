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

var Black pixcels.Pixcel = pixcels.Pixcel{0, 0, 0, pixcels.PixcelMax}

func main() {
	var dotFile string
	flag.StringVar(&dotFile, "d", "image/dot.png", "Dot image file")
	var outFile string
	flag.StringVar(&outFile, "o", "output", "Output name")
	var colorFile string
	flag.StringVar(&colorFile, "c", "colors.yml", "Colormap file")
	var width int
	flag.IntVar(&width, "w", 256, "Image maximum width in dot")

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

	printer := NewPrinter(dot, colors["-"])

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, ":"):
			cn := strings.TrimPrefix(arg, ":")
			col, ok := colors[cn]
			if !ok {
				fmt.Fprintf(os.Stderr, "Cannot find color '%s'\n", cn)
				os.Exit(1)
			}
			printer.Fill = col
		case arg == "_-":
			printer.Background = Black
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

	printer.Image.SaveTo(outFile + ".png")
}

type Printer struct {
	Background pixcels.Pixcel
	Fill       pixcels.Pixcel
	Dot        *pixcels.PixcelArray
	Image      *pixcels.PixcelArray
}

func NewPrinter(dot *pixcels.PixcelArray, def pixcels.Pixcel) *Printer {
	return &Printer{
		Background: Black,
		Fill:       def,
		Dot:        dot,
		Image:      pixcels.NewPixcelArray(0, 16*dot.YSize),
	}
}

func (pr *Printer) Print(str string) {
	for _, r := range str {
		imageName := runes.ImageName(r)
		runeImg, err := pixcels.Load(imageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Rune '%v' not found", r)
			os.Exit(1)
		}

		runeImg = runeImg.MultColor(pr.Fill)
		runeImg = runeImg.BackColor(pr.Background)

		charImg := runeImg.DotImage(pr.Dot)
		pr.Image.JoinHorizontal(charImg)

		fmt.Printf("'%v' %s\n", string(r), runes.RuneCode(r))
	}
}
