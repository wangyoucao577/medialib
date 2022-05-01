package tag

import "io"

// Tag represents tag mandantory functions.
type Tag interface {
	ParsePayload(io.Reader) error

	GetTagHeader() Header
	Size() int64 // total bytes of the tag, equal to (HeaderSize(11bytes) + DataSize)
}
