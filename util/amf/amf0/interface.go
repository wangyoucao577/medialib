package amf0

import "io"

type decoder interface {
	Decode(io.Reader) (int, error)
}
