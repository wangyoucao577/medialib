package main

import (
	"flag"
	"fmt"

	"github.com/wangyoucao577/medialib/util"
	"github.com/wangyoucao577/medialib/util/dump"
)

var flags struct {
	inputFilePath  string
	outputFilePath string
	outputFormat   string

	parseES bool // parse and dump es layer rather than container layer

	dumpBoxTypes  bool
	dumpNALUTypes bool
}

func init() {
	flag.StringVar(&flags.inputFilePath, "i", "", fmt.Sprintf("input file url, '%s' if stdin", util.InputStdin))
	flag.StringVar(&flags.outputFilePath, "o", "", "output to file instead of stdout")
	flag.StringVar(&flags.outputFormat, "of", dump.FormatJSON, fmt.Sprintf("output format, available values:%s", dump.FormatsHelper()))

	flag.BoolVar(&flags.parseES, "parse_es", false, "parse and dump Elementry Stream layer rather than container layer")

	flag.BoolVar(&flags.dumpBoxTypes, "box_types", false, "dump supported mp4 box types")
	flag.BoolVar(&flags.dumpNALUTypes, "nalu_types", false, "dump supported NALU types")
}

func validateFlags() error {
	_, err := dump.GetFormat(flags.outputFormat)
	if err != nil {
		return err
	}

	if !flags.dumpBoxTypes && !flags.dumpNALUTypes && len(flags.inputFilePath) == 0 {
		return fmt.Errorf("input file is mandantory")
	}

	return nil
}
