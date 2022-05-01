package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4/box"
	"github.com/wangyoucao577/medialib/util/exit"
	"github.com/wangyoucao577/medialib/util/marshaler"
	"github.com/wangyoucao577/medialib/video/avc/nalu"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	format, err := marshaler.GetFormat(flags.format)
	if err != nil {
		glog.Error(err)
		exit.Fail()
	}

	var data []byte

	contentFlag := getContentFlag()
	if contentFlag == flagContentBoxTypes {
		data, err = marshaler.Marshal(box.TypesMarshaler{}, format)
	} else if contentFlag == flagContentNALUTypes {
		data, err = marshaler.Marshal(nalu.TypesMarshaler{}, format)
	} else { // need to parse

		if len(flags.inputFilePath) == 0 {
			glog.Error("Input file is required.")
			exit.Fail()
		}

		data, err = parseMP4(flags.inputFilePath, format, getContentFlag())
	}

	if err != nil {
		glog.Error(err)
	} else {
		fmt.Println(string(data))
	}

}
