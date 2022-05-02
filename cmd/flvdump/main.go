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

	h := flv.NewHandler(flags.inputFilePath)
	if err := h.Parse(); err != nil {
		if err != io.EOF {
			glog.Errorf("Parse ES failed, err %v", err)
			exit.Fail()
		}
	}

	if contentType == dump.ContentTypeTags {
		if data, err := dump.Marshal(h.FLV, format); err != nil {
			glog.Error(err)
		} else {
			fmt.Println(string(data))
		}
	}

}
