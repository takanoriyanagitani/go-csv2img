package csv2img

import (
	"image"
	"image/color"
	"image/gif"
	jpg "image/jpeg"
	"image/png"
	"io"
	"iter"
	"os"
	"strconv"
)

type CsvRow []string

type CsvRowToGray8 func(row CsvRow, buf []uint8) ([]uint8, error)

type IterToImageGray8 func(iter.Seq2[CsvRow, error]) (ImageGray8, error)

func (c IterToImageGray8) IterToImage(
	i iter.Seq2[CsvRow, error],
) (image.Image, error) {
	g, e := c(i)
	switch e {
	case nil:
		return g.AsImage(), nil
	default:
		return nil, e
	}
}

func (c CsvRowToGray8) ToIterToImageGray8(
	rect image.Rectangle,
) IterToImageGray8 {
	return func(rows iter.Seq2[CsvRow, error]) (ImageGray8, error) {
		ret := ImageGray8{
			Rectangle: rect,
			data:      nil,
		}

		var cnt int = ret.TotalPixelCount()
		var buf []uint8 = make([]uint8, 0, ret.Width())
		var dat []uint8 = make([]uint8, 0, cnt)

		for row, e := range rows {
			if nil != e {
				return ret, e
			}

			cols, e := c(row, buf)
			if nil != e {
				return ret, e
			}

			dat = append(dat, cols...)
		}

		if cnt == len(dat) {
			ret.data = dat
		}

		return ret, nil
	}
}

func RectangleFromSize(width, height int) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{
			X: width,
			Y: height,
		},
	}
}

type StrToGray8 func(string) (uint8, error)

var StrToGray8Default StrToGray8 = func(s string) (uint8, error) {
	i, e := strconv.Atoi(s)
	return uint8(i), e
}

func (c StrToGray8) ToCsvRowToGray8() CsvRowToGray8 {
	return func(row CsvRow, buf []uint8) ([]uint8, error) {
		var ret []uint8 = buf[:0]
		for _, col := range row {
			u, e := c(col)
			if nil != e {
				return nil, e
			}
			ret = append(ret, u)
		}
		return ret, nil
	}
}

var CsvRowToGray8Default CsvRowToGray8 = StrToGray8Default.ToCsvRowToGray8()

type PointToIndex func(image.Point) int

type ImageGray8 struct {
	image.Rectangle
	data []uint8
}

func (g ImageGray8) Data(index int) color.Color {
	if len(g.data) <= index {
		return nil
	}
	var y uint8 = g.data[index]
	return color.Gray{Y: y}
}

func (g ImageGray8) Width() int           { return g.Rectangle.Dx() }
func (g ImageGray8) Height() int          { return g.Rectangle.Dy() }
func (g ImageGray8) TotalPixelCount() int { return g.Width() * g.Height() }

func (g ImageGray8) ToPointToIndex() PointToIndex {
	var width int = g.Width()
	return func(p image.Point) int {
		var y int = p.Y
		var x int = p.X
		return y*width + x
	}
}

func (g ImageGray8) At(x, y int) color.Color {
	var ix int = g.ToPointToIndex()(image.Point{X: x, Y: y})
	return g.Data(ix)
}

func (g ImageGray8) ColorModel() color.Model { return color.GrayModel }
func (g ImageGray8) Bounds() image.Rectangle { return g.Rectangle }

func (g ImageGray8) AsImage() image.Image { return g }

func (g ImageGray8) ToPng(w io.Writer) error { return png.Encode(w, g) }
func (g ImageGray8) ToJpg(w io.Writer) error { return jpg.Encode(w, g, nil) }
func (g ImageGray8) ToGif(w io.Writer) error { return gif.Encode(w, g, nil) }

func ImageToPngToStdout(i image.Image) error { return png.Encode(os.Stdout, i) }

func ImageToGifToStdout(i image.Image) error {
	return gif.Encode(os.Stdout, i, nil)
}

func ImageToJpgToStdout(i image.Image) error {
	return jpg.Encode(os.Stdout, i, nil)
}
