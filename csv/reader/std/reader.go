package reader

import (
	"encoding/csv"
	"io"
	"iter"
	"os"

	ci "github.com/takanoriyanagitani/go-csv2img"
	. "github.com/takanoriyanagitani/go-csv2img/util"
)

type Reader struct{ *csv.Reader }

func (r Reader) ToIter() iter.Seq2[ci.CsvRow, error] {
	return func(yield func(ci.CsvRow, error) bool) {
		for {
			row, e := r.Reader.Read()
			if io.EOF == e {
				return
			}

			if !yield(row, e) {
				return
			}
		}
	}
}

func ReaderToIter(rdr io.Reader) iter.Seq2[ci.CsvRow, error] {
	var cr *csv.Reader = csv.NewReader(rdr)
	cr.ReuseRecord = true
	return Reader{cr}.ToIter()
}

var CsvRowsStdin IO[iter.Seq2[ci.CsvRow, error]] = Of(ReaderToIter(os.Stdin))
