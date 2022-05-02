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

	var data []byte

	contentFlag := getContentFlag()
	if contentFlag == flagContentNALUTypes {
		data, err = dump.Marshal(nalu.TypesMarshaler{}, format)
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
		data, err = dump.Marshal(&h.ElementaryStream, format)
	}

	if err != nil {
		glog.Error(err)
	} else {
		fmt.Println(string(data))
	}
}
