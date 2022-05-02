package main

import (
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/flv"
	"github.com/wangyoucao577/medialib/util/dump"
	"github.com/wangyoucao577/medialib/util/exit"
)

func parseFLV(inputFile string, format dump.Format, contentType dump.ContentType, output string) error {

	// parse
	h := flv.NewHandler(flags.inputFilePath)
	if err := h.Parse(); err != nil {
		if err != io.EOF {
			glog.Errorf("Parse ES failed, err %v", err)
			exit.Fail()
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
	if contentType == dump.ContentTypeES || contentType == dump.ContentTypeRawES {
		es, err := h.FLV.ExtractES()
		if err != nil {
			return fmt.Errorf("extract es failed, err %v", err)
		}

		// print AVC ES
		if contentType == dump.ContentTypeRawES {
			_, err = es.Dump(w)
		} else {
			err = dump.DumpToWriter(es, format, w)
		}
		if err != nil {
			return fmt.Errorf("dump es failed, err %v", err)
		}
		return nil
	}

	// parse avc/hevc annexb es and print
	if contentType == dump.ContentTypeRawAnnexBES {
		es, err := h.FLV.ExtractAnnexBES()
		if err != nil {
			return fmt.Errorf("extract annexb_es failed, err %v", err)
		}

		// print AVC ES
		if _, err := es.Dump(w); err != nil {
			return fmt.Errorf("dump annexb_es failed, err %v", err)
		}
		return nil
	}

	// print flv boxes
	if err = dump.DumpToWriter(h.FLV, format, w); err != nil {
		return fmt.Errorf("dump flv tags failed, err %v", err)
	}
	return nil
}
