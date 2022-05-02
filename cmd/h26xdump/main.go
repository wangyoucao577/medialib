package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util/dump"
	"github.com/wangyoucao577/medialib/util/exit"
	"github.com/wangyoucao577/medialib/video/avc/annexbes"
	"github.com/wangyoucao577/medialib/video/avc/nalu"
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

	if flags.outputFilePath == dump.OutputStdout {
		defer fmt.Println() // new line to avoid `%` displayed at the end in Mac shell
	}

	var data dump.Marshaler
	if contentType == dump.ContentTypeNALUTypes {
		data = nalu.TypesMarshaler{}
	} else { // need to parse

		if len(flags.inputFilePath) == 0 {
			glog.Error("Input file is required.")
			exit.Fail()
		}

		h := annexbes.NewHandler(flags.inputFilePath)
		if err := h.Parse(); err != nil {
			if err != io.EOF {
				glog.Errorf("Parse ES failed, err %v", err)
				exit.Fail()
			}
		}
		data = &h.ElementaryStream
	}

	if err := dump.Dump(data, format, flags.outputFilePath); err != nil {
		glog.Error(err)
		exit.Fail()
	}
}
