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

type PixcelArray struct {
	XSize   int
	YSize   int
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
	xSize := bounds.Dx()
	ySize := bounds.Dy()
	pa := NewPixcelArray(xSize, ySize)

	for y := 0; y < ySize; y++ {
		for x := 0; x < xSize; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pa.Pixcels[y][x] = Pixcel{r, g, b, a}
		}
	}

	return pa, nil
}

func NewPixcelArray(xSize, ySize int) *PixcelArray {
	pixcels := make([][]Pixcel, ySize)
	for y := 0; y < ySize; y++ {
		pixcels[y] = make([]Pixcel, xSize)
	}
	return &PixcelArray{
		XSize:   xSize,
		YSize:   ySize,
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
	rect := image.Rect(0, 0, pa.XSize, pa.YSize)
	img := image.NewRGBA(rect)

	for y := 0; y < pa.YSize; y++ {
		for x := 0; x < pa.XSize; x++ {
			img.Set(x, y, pa.Pixcels[y][x])
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
	if pa.YSize != right.YSize {
		panic("Mismatch PixcelArray YSize")
	}

	pa.XSize = pa.XSize + right.XSize

	for y := 0; y < pa.YSize; y++ {
		line := make([]Pixcel, 0, pa.XSize)
		line = append(line, pa.Pixcels[y]...)
		line = append(line, right.Pixcels[y]...)
		pa.Pixcels[y] = line
	}
}

func (pa *PixcelArray) JoinVertical(bottom *PixcelArray) {
	if pa.XSize != bottom.XSize {
		panic("Mismatch PixcelArray XSize")
	}

	pa.YSize = pa.YSize + bottom.YSize

	pixcels := make([][]Pixcel, 0, pa.YSize)
	pixcels = append(pixcels, pa.Pixcels...)
	pixcels = append(pixcels, bottom.Pixcels...)
	pa.Pixcels = pixcels
}

func (pa *PixcelArray) MultColor(p Pixcel) *PixcelArray {
	npa := NewPixcelArray(pa.XSize, pa.YSize)

	mul32 := func(x, y uint32) uint32 {
		return uint32((uint64(x) * uint64(y)) >> 16)
	}

	for y := 0; y < pa.YSize; y++ {
		for x := 0; x < pa.XSize; x++ {
			px := pa.Pixcels[y][x]

			npa.Pixcels[y][x] = Pixcel{
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
	npa := NewPixcelArray(pa.XSize, pa.YSize)

	for y := 0; y < pa.YSize; y++ {
		for x := 0; x < pa.XSize; x++ {
			px := pa.Pixcels[y][x]
			if px.A == 0 {
				npa.Pixcels[y][x] = p
			} else {
				px.A = PixcelMax
				npa.Pixcels[y][x] = px
			}
		}
	}

	return npa
}

func (pa *PixcelArray) DotImage(dot *PixcelArray) *PixcelArray {
	npa := NewPixcelArray(pa.XSize*dot.XSize, 0)

	for y := 0; y < pa.YSize; y++ {
		line := NewPixcelArray(0, dot.YSize)
		for x := 0; x < pa.XSize; x++ {
			d := dot.MultColor(pa.Pixcels[y][x])
			line.JoinHorizontal(d)
		}
		npa.JoinVertical(line)
	}

	return npa
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
