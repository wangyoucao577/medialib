package annexbes

import (
	"os"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util"
)

// Handler represents AnnexB format Elementary Stream handler.
type Handler struct {
	ElementaryStream

	f        *os.File
	filePath string
}

// New creates Elementary Stream Handler.
func New(filePath string) *Handler {
	return &Handler{
		filePath: filePath,
	}
}

// Parse parses Elementary Stream file.
func (h *Handler) Parse() error {

	if err := h.open(); err != nil {
		glog.Warningf("open %s failed, err %v", h.filePath, err)
		return err
	}
	defer h.close()

	_, err := h.ElementaryStream.Parse(h.f, 0)
	return err
}

// Open opens Elementary Stream file.
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

// Close closes the Elementary Stream file handler.
func (h *Handler) close() error {
	if h == nil || h.f == nil || h.f == os.Stdin {
		return nil
	}

	return h.f.Close()
}
