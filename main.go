package main

import (
	"encoding/csv"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	in     = kingpin.Arg("input", "The input file to use, if not specified, stdin is used").File()
	out    = kingpin.Arg("out", "The output file to use, if not specified, stdout is used").String()
	h      = kingpin.Flag("header", "Whether or not to skip the first row or not").Bool()
	latIdx = kingpin.Flag("lat", "The column index of the latitude column").Required().Int()
	lngIdx = kingpin.Flag("lng", "The column index of the longitude column").Required().Int()
	delim  = kingpin.Flag("d", "The delimiter to use for CSV files").Default(",").String()
	name   = kingpin.Flag("name", "Column name for WKT").Default("WKT").String()
)

func main() {
	kingpin.Parse()
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var input io.ReadCloser
	if in != nil && *in != nil {
		input = *in
	} else {
		input = os.Stdin
	}
	go closeOnEnd(input)
	var output io.WriteCloser
	if out != nil && len(*out) > 0 {
		f, err := os.OpenFile(*out, os.O_CREATE, os.ModeAppend)
		if err != nil {
			return err
		}
		output = f
	} else {
		output = os.Stdout
	}
	reader := csv.NewReader(input)
	reader.Comma = rune((*delim)[0])
	writer := csv.NewWriter(output)
	writer.Comma = rune((*delim)[0])

	var row int
	if *h {
		header, err := reader.Read()
		if err != nil {
			return fmt.Errorf("failed to read CSV: %s", err)
		}
		headers := removeElements(header, *latIdx, *lngIdx)
		headers[len(headers)-1] = *name
		err = writer.Write(headers)
		if err == io.ErrClosedPipe {
			return nil
		}
		if err != nil {
			return err
		}
	}
	for {
		rec, err := reader.Read()
		row++
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV: %s", err)
		}
		lat, lng := parseFloat(rec[*latIdx], row), parseFloat(rec[*lngIdx], row)
		p := orb.Point{lng, lat}
		dst := make([]string, len(rec)-1)
		var i int
		// Exclude lat/lng rows from CSV
		for j, s := range rec {
			if j != *latIdx && j != *lngIdx {
				dst[i] = s
				i++
			}
		}
		// Append wkt POINT geometry to end of csv
		dst[len(dst)-1] = wkt.MarshalString(p)
		err = writer.Write(dst)
		if err == io.ErrClosedPipe {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to write to csv: %v", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to close csv: %v", err)
	}
	return nil
}

func parseFloat(val string, row int) float64 {
	if len(val) == 0 {
		return 0
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Fatalf("failed to convert value %s to a float, row %d: %v\n", val, row, err)
	}
	return f
}

func removeElements(str []string, one, two int) []string {
	f := make([]string, len(str)-1)
	var i int
	for j := 0; j < len(str); j++ {
		if j != one && j != two {
			f[i] = str[j]
			i++
		}
	}
	return append(f, "wkt")
}

func closeOnEnd(r io.ReadCloser) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGKILL, syscall.SIGINT)
	<-sig
	err := r.Close()
	if err != nil {
		log.Fatal(err)
	}
}
