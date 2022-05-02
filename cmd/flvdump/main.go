package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/flv"
	"github.com/wangyoucao577/medialib/util/dump"
	"github.com/wangyoucao577/medialib/util/exit"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	// validate and get flags
	if err := validateFlags(); err != nil {
		glog.Error(err)
		exit.Fail()
	}
	format, _ := dump.GetFormat(flags.format) // ignore error since they've been validated already
	contentType, _ := getConentType()

	if flags.outputFilePath == dump.OutputStdout {
		defer fmt.Println() // new line to avoid `%` displayed at the end in Mac shell
	}

	// parse file
	h := flv.NewHandler(flags.inputFilePath)
	if err := h.Parse(); err != nil {
		if err != io.EOF {
			glog.Errorf("Parse ES failed, err %v", err)
			exit.Fail()
		}
	}

	if contentType == dump.ContentTypeTags {
		if err := dump.Dump(h.FLV, format, flags.outputFilePath); err != nil {
			glog.Error(err)
			exit.Fail()
		}
	}
}
