package main

import (
	"flag"
	"fmt"

	"github.com/wangyoucao577/medialib/util/dump"
)

var flags struct {
	inputFilePath string
	format        string
}

func init() {
	flag.StringVar(&flags.inputFilePath, "i", "", `Input FLV file path.`)
	flag.StringVar(&flags.format, "format", "json", fmt.Sprintf("Output format, available values:%s", dump.FormatsHelper()))
}
