package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/flv"
	"github.com/wangyoucao577/medialib/container/mp4"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util/dump"
	"github.com/wangyoucao577/medialib/util/exit"
	"github.com/wangyoucao577/medialib/util/mediaformat"
	"github.com/wangyoucao577/medialib/video/avc/annexbes"
	avcnalu "github.com/wangyoucao577/medialib/video/avc/nalu"
	hevcnalu "github.com/wangyoucao577/medialib/video/hevc/nalu"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	// validate and get flags
	if err := validateFlags(); err != nil {
		glog.Error(err)
		exit.Fail()
	}
	outputFormat, _ := dump.GetFormat(flags.outputFormat) // ignore error since they've been validated already

	if flags.outputFilePath == dump.OutputStdout {
		defer fmt.Println() // new line to avoid `%` displayed at the end in Mac shell
	}

	var data dump.Marshaler
	if flags.dumpBoxTypes {
		data = box.TypesMarshaler{}
	} else if flags.dumpAVCNALUTypes {
		data = avcnalu.TypesMarshaler{}
	} else if flags.dumpHEVCNALUTypes {
		data = hevcnalu.TypesMarshaler{}
	} else {
		if m, err := parseInput(flags.inputFilePath, flags.parseES, flags.printDurations); err != nil {
			glog.Error(err)
			exit.Fail()
		} else {
			data = m
		}
	}
	if data != nil {
		if err := dump.Dump(data, outputFormat, flags.outputFilePath); err != nil {
			glog.Error(err)
			exit.Fail()
		}
		return
	}

}

func parseInput(inputFilePath string, parseES bool, printDuration bool) (dump.Marshaler, error) {

	if strings.HasSuffix(inputFilePath, mediaformat.AsExtension(mediaformat.MP4)) ||
		strings.HasSuffix(inputFilePath, mediaformat.AsExtension(mediaformat.FMP4)) ||
		strings.HasSuffix(inputFilePath, mediaformat.AsExtension(mediaformat.M4S)) {
		m := mp4.New(inputFilePath)
		if err := m.Parse(); err != nil {
			if err != io.EOF {
				glog.Warningf("Parse mp4 failed but ignore to leverage the data has been parsed already, err %v", err)
				// exit.Fail()	// ignore the error so that able to leverage the data has been parsed already
			}
		}

		if printDuration {
			m.DumpDurations()
		}

		if !parseES {
			return m, nil
		}

		es, err := m.Boxes.ExtractES(0)
		if err != nil {
			return nil, fmt.Errorf("extract es failed, err %v", err)
		}
		return es, nil

	} else if strings.HasSuffix(inputFilePath, mediaformat.AsExtension(mediaformat.FLV)) {
		h := flv.New(inputFilePath)
		if err := h.Parse(); err != nil {
			if err != io.EOF {
				glog.Warningf("Parse FLV failed but ignore to leverage the data has been parsed already, err %v", err)
				// exit.Fail()	// ignore the error so that able to leverage the data has been parsed already
			}
		}
		if !parseES {
			return h, nil
		}

		es, err := h.FLV.ExtractES()
		if err != nil {
			return nil, fmt.Errorf("extract es failed, err %v", err)
		}
		return es, nil

	} else if strings.HasSuffix(inputFilePath, mediaformat.AsExtension(mediaformat.H264)) {
		h := annexbes.New(inputFilePath)
		if err := h.Parse(); err != nil {
			if err != io.EOF {
				glog.Warningf("Parse ES failed but ignore to leverage the data has been parsed already, err %v", err)
				// exit.Fail()	// ignore the error so that able to leverage the data has been parsed already
			}
		}
		return &h.ElementaryStream, nil
	}

	return nil, fmt.Errorf("unknown format for input %s", inputFilePath)
}
