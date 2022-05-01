package tag

import "io"

// Tag represents tag mandantory functions.
type Tag interface {
	ParsePayload(io.Reader) error

	GetTagHeader() Header
}
