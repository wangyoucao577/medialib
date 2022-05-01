package main

import (
	"io"
	"os"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/mp4"
	"github.com/wangyoucao577/medialib/util/exit"
	"github.com/wangyoucao577/medialib/util/marshaler"
)

func parseMP4(inputFile string, format marshaler.Format, contentType int) ([]byte, error) {

	// parse
	m := mp4.New(inputFile)
	if err := m.Parse(); err != nil {
		if err != io.EOF {
			glog.Errorf("Parse mp4 failed, err %v", err)
			exit.Fail()
		}
	}

	// parse avc/hevc es and print
	if contentType == flagContentES || contentType == flagContentRawES {
		if es, err := m.Boxes.ExtractES(0); err != nil {
			glog.Errorf("Extract ES failed, err %v", err)
			exit.Fail()
		} else {
			// print AVC ES
			if contentType == flagContentRawES {
				if _, err := es.Dump(os.Stdout); err != nil {
					glog.Errorf("Dump ES failed, err %v", err)
					exit.Fail()
				}
				return nil, nil
			} else {
				return marshaler.Marshal(es, format)
			}
		}
	}

	// parse avc/hevc annexb es and print
	if contentType == flagContentRawAnnexBES {
		if es, err := m.Boxes.ExtractAnnexBES(0); err != nil {
			glog.Errorf("Extract ES failed, err %v", err)
			exit.Fail()
		} else {
			// print AVC ES
			if _, err := es.Dump(os.Stdout); err != nil {
				glog.Errorf("Dump ES failed, err %v", err)
				exit.Fail()
			}
			return nil, nil
		}
	}

	// print mp4 boxes
	return marshaler.Marshal(m.Boxes, format)
}
