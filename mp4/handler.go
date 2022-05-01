//Package mp4 defines mp4 handlers and structures.
package mp4

import (
	"os"

	"github.com/golang/glog"
)

// Handler represents handler for `mp4` structure.
type Handler struct {
	Boxes

	f        *os.File
	filePath string
}

// New creates mp4 Handler.
func New(filePath string) *Handler {
	return &Handler{
		Boxes: newBoxes(),

		filePath: filePath,
	}
}

// Parse parses mp4 file.
func (h *Handler) Parse() error {

	if err := h.open(); err != nil {
		glog.Warningf("open %s failed, err %v", h.filePath, err)
		return err
	}
	defer h.close()

	return h.Boxes.ParsePayload(h.f)
}

// Open opens mp4 file.
func (h *Handler) open() error {

	var err error
	if h.f, err = os.Open(h.filePath); err != nil {
		return err
	}
	glog.V(1).Infof("open %s succeed.\n", h.filePath)

	return nil
}

// Close closes the mp4 file handler.
func (h *Handler) close() error {
	if h == nil || h.f == nil {
		return nil
	}

	return h.f.Close()
}
