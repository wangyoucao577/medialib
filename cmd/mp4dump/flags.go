package main

import (
	"flag"
	"fmt"

	"github.com/wangyoucao577/medialib/util/dump"
)

var flags struct {
	inputFilePath string
	content       string // content to output
	format        string
}

var supportedContentTypes = []dump.ContentType{
	dump.ContentTypeBoxTypes,
	dump.ContentTypeNALUTypes,

	dump.ContentTypeES,
	dump.ContentTypeRawES,
	dump.ContentTypeRawAnnexBES,
	dump.ContentTypeBoxes, // NOTE: put default at the end to align with `-h` shown
}

func supportedConentTypesHelper() string {
	var s string
	for _, n := range supportedContentTypes {
		s += "\n"
		s += n.FixedLenString()
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

	flag.StringVar(&flags.inputFilePath, "i", "", `Input mp4/fmp4 file url.`)
	flag.StringVar(&flags.content, "content", dump.ContentTypeES, fmt.Sprintf("Contents to parse and output, available values: %s", supportedConentTypesHelper()))
	flag.StringVar(&flags.format, "format", dump.FormatJSON, fmt.Sprintf("Output format, available values:%s", dump.FormatsHelper()))
}
