package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/wangyoucao577/medialib/util/marshaler"
)

var flags struct {
	inputFilePath string
	content       string // content to output
	format        string
}

func init() {
	flag.StringVar(&flags.inputFilePath, "i", "", `Input Elementary Stream file path, such as 'x.h264' or 'x.h265'.
Be aware that the Elementary Stream file is mandatory stored by AnnexB byte stream format.`)
	flag.StringVar(&flags.content, "content", "es", `Contents to parse and output, available values: 
  nalu_types: NALU types(no parse)  
  es: AVC/HEVC elementary stream parsing data`)
	flag.StringVar(&flags.format, "format", "json", fmt.Sprintf("Output format, available values:%s", marshaler.FormatsHelper()))
}

const (
	flagContentBoxes = iota // mp4 boxes
	flagContentES           // AVC/HEVC Elementary Stream parsing data

	// no parse needed
	flagContentNALUTypes
)

func getContentFlag() int {
	switch strings.ToLower(flags.content) {
	case "es":
		return flagContentES

	case "nalu_types":
		return flagContentNALUTypes
	}
	return flagContentBoxes
}
