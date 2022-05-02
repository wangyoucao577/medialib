package main

import (
	"flag"
	"fmt"

	"github.com/wangyoucao577/medialib/util/dump"
)

var flags struct {
	inputFilePath  string
	outputFilePath string
	content        string // content to output
	format         string
}

var supportedContentTypes = []dump.ContentType{
	dump.ContentTypeNALUTypes,
	dump.ContentTypeES, // NOTE: put default at the end to align with `-h` shown
}

func supportedConentTypesHelper() string {
	var maxLen int
	for _, n := range supportedContentTypes {
		if maxLen < len(n) {
			maxLen = len(n)
		}
	}

	var s string
	for _, n := range supportedContentTypes {
		s += "\n"
		s += n.FixedLenString(maxLen)
		s += ": "
		s += n.Description()
	}
	return s
}

func getConentType() (dump.ContentType, error) {
	for _, c := range supportedContentTypes {
		if c == dump.ContentType(flags.content) {
			return c, nil
		}
	}
	return "", fmt.Errorf("invalid content type %s", flags.content)
}

func init() {

	flag.StringVar(&flags.inputFilePath, "i", "", `Input Elementary Stream file path, such as 'x.h264' or 'x.h265'.
Be aware that the Elementary Stream file is mandatory stored by AnnexB byte stream format.`)
	flag.StringVar(&flags.content, "content", dump.ContentTypeES, fmt.Sprintf("Contents to parse and output, available values: %s", supportedConentTypesHelper()))
	flag.StringVar(&flags.format, "format", dump.FormatJSON, fmt.Sprintf("Output format, available values:%s", dump.FormatsHelper()))
	flag.StringVar(&flags.outputFilePath, "o", "stdout", "Output file path.")
}
