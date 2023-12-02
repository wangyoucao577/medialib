package box

import (
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// ParseBox tries to parse a box from a mount of data.
// return ErrUnknownBoxType if doesn't know, ErrInsufficientSize if no enough data,
// otherwise fatal error.
func ParseBox(r io.Reader, pb ParentBox, bytesAvailable uint64) (uint64, error) {
	boxHeader := Header{}
	if err := boxHeader.Parse(r, bytesAvailable); err != nil {
		if err == io.EOF {
			glog.V(1).Info("EOF")
			return boxHeader.HeaderSize(), err
		} else if err == ErrInsufficientSize { // ignore remain available bytes
			bytesToIgnore := bytesAvailable - boxHeader.HeaderSize()
			glog.Warningf("%v when parse header, ignore %s size %d", err, boxHeader.Type, bytesToIgnore)

			if bytesToIgnore > 0 {
				if err := util.ReadOrError(r, make([]byte, bytesToIgnore)); err != nil {
					return bytesAvailable, err
				}
			}
			return bytesAvailable, err
		}
		// glog.Warningf("parse box header failed, err %v", err)
		return boxHeader.HeaderSize(), err
	}
	bytesAvailable -= boxHeader.HeaderSize()

	b, err := pb.CreateSubBox(boxHeader)
	if err != nil {
		if err == ErrUnknownBoxType {
			bytesToIgnore := boxHeader.PayloadSize()
			if bytesToIgnore > bytesAvailable {
				bytesToIgnore = bytesAvailable
			}
			glog.Warningf("ignore %v when create sub box, type %s payload size %d (available %d)", err, boxHeader.Type, boxHeader.PayloadSize(), bytesAvailable)

			if bytesToIgnore > 0 {
				if err := util.ReadOrError(r, make([]byte, bytesToIgnore)); err != nil {
					return boxHeader.HeaderSize() + bytesToIgnore, err
				}
			}
			return boxHeader.HeaderSize() + bytesToIgnore, err
		}
		return boxHeader.HeaderSize(), err
	}

	if err := b.ParsePayload(r); err != nil {
		glog.Warningf("parse box type %s payload failed, err %v", string(boxHeader.Type[:]), err)
		return boxHeader.BoxSize(), err
	}

	return boxHeader.BoxSize(), nil
}
