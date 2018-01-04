package pixcels

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const (
	PixcelMax = 65535
)

type Pixcel struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

func (p Pixcel) RGBA() (uint32, uint32, uint32, uint32) {
	return p.R, p.G, p.B, p.A
}

func FromArray(v []float32) Pixcel {
	v = append(v, 1.0, 1.0, 1.0, 1.0)

	return Pixcel{
		R: uint32(v[0] * PixcelMax),
		G: uint32(v[1] * PixcelMax),
		B: uint32(v[2] * PixcelMax),
		A: uint32(v[3] * PixcelMax),
	}
}

type PixcelArray struct {
	Width   int
	Height  int
	Pixcels [][]Pixcel
}

func Load(fileName string) (*PixcelArray, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "File open error")
	}
	defer file.Close()

	img, err := loadImage(fileName, file)
	if err != nil {
		return nil, errors.Wrap(err, "IMage load error")
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	pa := NewPixcelArray(width, height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			pa.Pixcels[x][y] = Pixcel{r, g, b, a}
		}
	}

	return pa, nil
}

func NewPixcelArray(width, height int) *PixcelArray {
	pixcels := make([][]Pixcel, width)
	for x := 0; x < width; x++ {
		pixcels[x] = make([]Pixcel, height)
	}
	return &PixcelArray{
		Width:   width,
		Height:  height,
		Pixcels: pixcels,
	}
}

func loadImage(fileName string, file io.Reader) (image.Image, error) {
	switch {
	case strings.HasSuffix(fileName, ".jpeg"):
		fallthrough
	case strings.HasSuffix(fileName, ".jpg"):
		return jpeg.Decode(file)
	case strings.HasSuffix(fileName, ".png"):
		return png.Decode(file)
	}

	return nil, errors.New("Unknown file type")
}

func (pa *PixcelArray) SaveTo(fileName string) error {
	rect := image.Rect(0, 0, pa.Width, pa.Height)
	img := image.NewRGBA(rect)

	for x := 0; x < pa.Width; x++ {
		for y := 0; y < pa.Height; y++ {
			img.Set(x, y, pa.Pixcels[x][y])
		}
	}

	switch {
	case strings.HasSuffix(fileName, ".jpeg"):
		fallthrough
	case strings.HasSuffix(fileName, ".jpg"):
		file, err := os.Create(fileName)
		if err != nil {
			return errors.Wrap(err, "Failed open file")
		}
		jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
		file.Close()
	case strings.HasSuffix(fileName, ".png"):
		file, err := os.Create(fileName)
		if err != nil {
			return errors.Wrap(err, "Failed open file")
		}
		png.Encode(file, img)
		file.Close()
	}

	return nil
}

func (pa *PixcelArray) JoinHorizontal(right *PixcelArray) {
	if pa.Height != right.Height {
		panic("Mismatch PixcelArray Height")
	}

	pa.Width = pa.Width + right.Width

	pixcels := make([][]Pixcel, 0, pa.Width)
	pixcels = append(pixcels, pa.Pixcels...)
	pixcels = append(pixcels, right.Pixcels...)
	pa.Pixcels = pixcels
}

func (pa *PixcelArray) JoinVertical(bottom *PixcelArray) {
	if pa.Width != bottom.Width {
		panic("Mismatch PixcelArray Width")
	}

	pa.Height = pa.Height + bottom.Height

	for x := 0; x < pa.Width; x++ {
		line := make([]Pixcel, 0, pa.Width)
		line = append(line, pa.Pixcels[x]...)
		line = append(line, bottom.Pixcels[x]...)
		pa.Pixcels[x] = line
	}
}

func (pa *PixcelArray) MultColor(p Pixcel) *PixcelArray {
	npa := NewPixcelArray(pa.Width, pa.Height)

	mul32 := func(x, y uint32) uint32 {
		return uint32((uint64(x) * uint64(y)) >> 16)
	}

	for x := 0; x < pa.Width; x++ {
		for y := 0; y < pa.Height; y++ {
			px := pa.Pixcels[x][y]

			npa.Pixcels[x][y] = Pixcel{
				R: mul32(px.R, p.R),
				G: mul32(px.G, p.G),
				B: mul32(px.B, p.B),
				A: mul32(px.A, p.A),
			}
		}
	}

	return npa
}

func (pa *PixcelArray) BackColor(p Pixcel) *PixcelArray {
	npa := NewPixcelArray(pa.Width, pa.Height)

	for x := 0; x < pa.Width; x++ {
		for y := 0; y < pa.Height; y++ {
			px := pa.Pixcels[x][y]
			if px.A == 0 {
				npa.Pixcels[x][y] = p
			} else {
				px.A = PixcelMax
				npa.Pixcels[x][y] = px
			}
		}
	}

	return npa
}

func (pa *PixcelArray) DotImage(dot *PixcelArray) *PixcelArray {
	npa := NewPixcelArray(0, pa.Height*dot.Height)

	for x := 0; x < pa.Width; x++ {
		line := NewPixcelArray(dot.Height, 0)
		for y := 0; y < pa.Height; y++ {
			d := dot.MultColor(pa.Pixcels[x][y])
			line.JoinVertical(d)
		}
		npa.JoinHorizontal(line)
	}

	return npa
}
