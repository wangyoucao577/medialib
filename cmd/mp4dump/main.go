package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util/dump"
	"github.com/wangyoucao577/medialib/util/exit"
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
	var data []byte

	if contentType == dump.ContentTypeBoxTypes {
		data, err = dump.Marshal(box.TypesMarshaler{}, format)
	} else if contentType == dump.ContentTypeNALUTypes {
		data, err = dump.Marshal(nalu.TypesMarshaler{}, format)
	} else { // need to parse

		if len(flags.inputFilePath) == 0 {
			glog.Error("Input file is required.")
			exit.Fail()
		}

		data, err = parseMP4(flags.inputFilePath, format, contentType)
	}

	if err != nil {
		glog.Error(err)
	} else {
		fmt.Println(string(data))
	}

}
