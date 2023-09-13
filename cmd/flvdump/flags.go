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
	content        string // content to output
	format         string
}

var supportedContentTypes = []dump.ContentType{
	dump.ContentTypeNALUTypes,

	dump.ContentTypeES,
	dump.ContentTypeTags, // NOTE: put default at the end to align with `-h` shown
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
	flag.StringVar(&flags.inputFilePath, "i", "", fmt.Sprintf("Input flv file url, '%s' if stdin", util.InputStdin))
	flag.StringVar(&flags.format, "format", dump.FormatJSON, fmt.Sprintf("Output format, available values:%s", dump.FormatsHelper()))
	flag.StringVar(&flags.content, "content", dump.ContentTypeTags, fmt.Sprintf("Contents to parse and output, available values: %s", supportedConentTypesHelper()))
	flag.StringVar(&flags.outputFilePath, "o", "stdout", "Output file path.")
}

func validateFlags() error {
	_, err := dump.GetFormat(flags.format)
	if err != nil {
		return err
	}

	contentType, err := getConentType()
	if err != nil {
		return err
	}

	if len(flags.outputFilePath) == 0 {
		return fmt.Errorf("output should not be empty")
	}

	if contentType == dump.ContentTypeES ||
		contentType == dump.ContentTypeTags {
		if len(flags.inputFilePath) == 0 {
			return fmt.Errorf("input file is required")
		}
	}

	return nil
}
