package main

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4"
	"github.com/wangyoucao577/medialib/util/dump"
)

func parseMP4(inputFile string, contentType dump.ContentType, output string) error {

	// parse
	m := mp4.New(inputFile)
	if err := m.Parse(); err != nil {
		if err != io.EOF {
			glog.Warningf("Parse mp4 failed but ignore to leverage the data has been parsed already, err %v", err)
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

	// parse avc es and print
	switch contentType {
	case dump.ContentTypeRawES:
		es, err := m.Boxes.ExtractES(0)
		if err != nil {
			return fmt.Errorf("extract es failed, err %v", err)
		}
		if _, err := es.Dump(w); err != nil {
			return fmt.Errorf("dump es failed, err %v", err)
		}
	case dump.ContentTypeRawAnnexBES:
		es, err := m.Boxes.ExtractAnnexBES(0)
		if err != nil {
			return fmt.Errorf("extract annexb_es failed, err %v", err)
		}
		if _, err := es.Dump(w); err != nil {
			return fmt.Errorf("dump annexb_es failed, err %v", err)
		}
	}

	return nil
}
