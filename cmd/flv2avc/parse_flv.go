package main

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/flv"
	"github.com/wangyoucao577/medialib/util/dump"
)

func parseFLV(inputFile string, contentType dump.ContentType, output string) error {

	// parse
	h := flv.NewHandler(flags.inputFilePath)
	if err := h.Parse(); err != nil {
		if err != io.EOF {
			glog.Warningf("Parse FLV failed but ignore to leverage the data has been parsed already, err %v", err)
			// exit.Fail()	// ignore the error so that able to leverage the data has been parsed already
		}
	}

	// output
	w, closer, err := dump.CreateOutput(output)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer.Close()
	}

	// parse avc/hevc es and print
	switch contentType {
	case dump.ContentTypeRawES:
		es, err := h.FLV.ExtractES()
		if err != nil {
			return fmt.Errorf("extract es failed, err %v", err)
		}
		if _, err = es.Dump(w); err != nil {
			return fmt.Errorf("dump es failed, err %v", err)
		}
	case dump.ContentTypeRawAnnexBES:
		es, err := h.FLV.ExtractAnnexBES()
		if err != nil {
			return fmt.Errorf("extract annexb_es failed, err %v", err)
		}
		if _, err := es.Dump(w); err != nil {
			return fmt.Errorf("dump annexb_es failed, err %v", err)
		}
	}

	return nil
}
