package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/protocol"
	"github.com/wangyoucao577/medialib/protocol/rtmp"
	"github.com/wangyoucao577/medialib/util/dump"
	"github.com/wangyoucao577/medialib/util/exit"
	"github.com/wangyoucao577/medialib/video/avc/nalu"
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

	if contentType == dump.ContentTypeNALUTypes {
		data := nalu.TypesMarshaler{}
		if err := dump.Dump(data, format, flags.outputFilePath); err != nil {
			glog.Error(err)
			exit.Fail()
		}
		return
	}

	if strings.HasPrefix(flags.inputFilePath, protocol.SchemaRTMP) { // rtmp
		h, err := rtmp.NewHandler(flags.inputFilePath)
		if err != nil {
			glog.Error(err)
			exit.Fail()
		}
		if err := h.Connect(time.Second * 3); err != nil {
			glog.Error(err)
			exit.Fail()
		}
		defer h.Close()
		return
	}

	if err := parseFLV(flags.inputFilePath, format, contentType, flags.outputFilePath); err != nil {
		glog.Error(err)
		exit.Fail()
	}
}
