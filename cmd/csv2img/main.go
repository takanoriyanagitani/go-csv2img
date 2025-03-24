package main

import (
	"context"
	"fmt"
	"image"
	"iter"
	"log"
	"os"
	"strconv"

	ci "github.com/takanoriyanagitani/go-csv2img"
	rs "github.com/takanoriyanagitani/go-csv2img/csv/reader/std"
	. "github.com/takanoriyanagitani/go-csv2img/util"
)

var envValByKey func(string) IO[string] = Lift(
	func(key string) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	},
)

var width IO[int] = Bind(envValByKey("ENV_WIDTH"), Lift(strconv.Atoi))

var height IO[int] = Bind(envValByKey("ENV_HEIGHT"), Lift(strconv.Atoi))

var rect IO[image.Rectangle] = Bind(
	All(width, height),
	Lift(func(i []int) (image.Rectangle, error) {
		return ci.RectangleFromSize(i[0], i[1]), nil
	}),
)

var row2gray8 ci.CsvRowToGray8 = ci.CsvRowToGray8Default

var iter2grayImg8 IO[ci.IterToImageGray8] = Bind(
	rect,
	Lift(func(r image.Rectangle) (ci.IterToImageGray8, error) {
		return row2gray8.ToIterToImageGray8(r), nil
	}),
)

var rows IO[iter.Seq2[ci.CsvRow, error]] = rs.CsvRowsStdin

var img IO[image.Image] = Bind(
	iter2grayImg8,
	func(i2g ci.IterToImageGray8) IO[image.Image] {
		return Bind(
			rows,
			Lift(i2g.IterToImage),
		)
	},
)

var img2png2stdout IO[Void] = Bind(
	img,
	Lift(func(i image.Image) (Void, error) {
		return Empty, ci.ImageToPngToStdout(i)
	}),
)

func main() {
	_, e := img2png2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
