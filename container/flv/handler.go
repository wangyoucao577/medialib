package flv

import (
	"os"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// Handler represents handler for `flv` structure.
type Handler struct {
	FLV

	f        *os.File
	filePath string
}

// NewHandler creates FLV Handler.
func NewHandler(filePath string) *Handler {
	return &Handler{
		filePath: filePath,
	}
}

// Parse parses FLV file.
func (h *Handler) Parse() error {

	if err := h.open(); err != nil {
		glog.Warningf("open %s failed, err %v", h.filePath, err)
		return err
	}
	defer h.close()

	return h.FLV.Parse(h.f)
}

// Open opens FLV file.
func (h *Handler) open() error {

	if h.filePath == util.InputStdin {
		h.f = os.Stdin
	} else {
		var err error
		if h.f, err = os.Open(h.filePath); err != nil {
			return err
		}
	}

	glog.V(1).Infof("open %s succeed.\n", h.filePath)

	return nil
}

// Close closes the FLV file handler.
func (h *Handler) close() error {
	if h == nil || h.f == nil || h.f == os.Stdin {
		return nil
	}

	return h.f.Close()
}
