package main

import (
	"encoding/csv"
	"github.com/alecthomas/kingpin"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	in    = kingpin.Arg("input", "The input file to use, if not specified, stdin is used").File()
	out   = kingpin.Arg("out", "The output file to use, if not specified, stdout is used").File()
	h = kingpin.Flag("header", "Whether or not to skip the first row or not").Bool()
	latIdx   = kingpin.Flag("lat", "The column index of the latitude column").Required().Int()
	lngIdx   = kingpin.Flag("lng", "The column index of the longitude column").Required().Int()
	delim = kingpin.Flag("d", "The delimiter to use for CSV files").Default(",").String()
)

func main() {
	kingpin.Parse()
	var input io.Reader
	if in != nil && *in != nil {
		input = *in
	} else {
		input = os.Stdin
	}
	var output io.WriteCloser
	if out != nil && *out != nil {
		output = *out
	} else {
		output = os.Stdout
	}
	reader := csv.NewReader(input)
	reader.Comma = rune((*delim)[0])
	writer := csv.NewWriter(output)
	writer.Comma = rune((*delim)[0])

	var row int
	if *h {
		_, err := reader.Read()
		if err != nil {
			log.Fatalf("failed to read CSV: %s\n", err)
		}
	}
	for {
		rec, err := reader.Read()
		row++
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read CSV: %s\n", err)
		}
		lat, lng := parseFloat(rec[*latIdx], row), parseFloat(rec[*lngIdx], row)
		p := orb.Point{lng, lat}
		dst := make([]string, len(rec) - 1)
		var i int
		for j, s := range rec {
			if j != *latIdx && j != *lngIdx {
				dst[i] = s
				i++
			}
		}
		dst[len(dst) - 1] = wkt.MarshalString(p)
		err = writer.Write(dst)
		if err == io.ErrClosedPipe {
			break
		}
		if err != nil {
			log.Fatalf("unable to write to csv: %v\n", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatalf("failed to close csv: %v\n", err)
	}
}

func logParseFloatErr(val string, row int, err error) {
	log.Fatalf("failed to convert value %s to a float, row %d: %v\n", val, row, err)
}

func min(is ...int) int {
	m := is[0]
	for _, i := range is {
		if i < m {
			m = i
		}
	}
	return m
}

func max(is ...int) int {
	m := is[0]
	for _, i := range is {
		if i > m {
			m = i
		}
	}
	return m
}


func parseFloat(val string, row int) float64 {
	if len(val) == 0 {
		return 0
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		logParseFloatErr(val, row, err)
	}
	return f
}
