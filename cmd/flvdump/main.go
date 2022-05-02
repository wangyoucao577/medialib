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

	format, err := dump.GetFormat(flags.format)
	if err != nil {
		glog.Error(err)
		exit.Fail()
	}
	contentType, err := getConentType()
	if err != nil {
		glog.Error(err)
		exit.Fail()
	}

	if len(flags.inputFilePath) == 0 {
		glog.Error("Input file is required.")
		exit.Fail()
	}

	if flags.outputFilePath == dump.OutputStdout {
		defer fmt.Println() // new line to avoid `%` displayed at the end in Mac shell
	}

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
