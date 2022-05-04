// Package stsc represents Sample To Chunk Box.
package stsc

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/container/mp4/box"
	"github.com/wangyoucao577/medialib/util"
)

// ChunkEntry reprensents a chunk box entry.
type ChunkEntry struct {
	FirstChunk             uint32 `json:"first_chunk"`
	SamplesPerChunk        uint32 `json:"samples_per_chunk"`
	SampleDescriptionIndex uint32 `json:"sample_description_index"`
}

// Box represents a stsc box.
type Box struct {
	box.FullHeader `json:"full_header"`

	EntryCount uint32       `json:"entry_count"`
	Entries    []ChunkEntry `json:"entry,omitempty"`
}

// New creates a new Box.
func New(h box.Header) box.Box {
	return &Box{
		FullHeader: box.FullHeader{
			Header: h,
		},
	}
}

// ParsePayload parse payload which requires basic box already exist.
func (b *Box) ParsePayload(r io.Reader) error {
	if err := b.Validate(); err != nil {
		glog.Warningf("box %s invalid, err %v", b.Type, err)
		return nil
	}

	// parse full header additional information first
	if err := b.FullHeader.ParseVersionFlag(r); err != nil {
		return err
	}

	// start to parse payload
	var parsedBytes uint64

	data := make([]byte, 4)
	if err := util.ReadOrError(r, data); err != nil {
		return err
	} else {
		b.EntryCount = binary.BigEndian.Uint32(data)
		parsedBytes += 4
	}

	for i := 0; i < int(b.EntryCount); i++ {
		entry := ChunkEntry{}

		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			entry.FirstChunk = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}

		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			entry.SamplesPerChunk = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}

		if err := util.ReadOrError(r, data); err != nil {
			return err
		} else {
			entry.SampleDescriptionIndex = binary.BigEndian.Uint32(data)
			parsedBytes += 4
		}

		b.Entries = append(b.Entries, entry)
	}

	if parsedBytes != b.PayloadSize() {
		return fmt.Errorf("box %s parsed bytes != payload size: %d != %d", b.Type, parsedBytes, b.PayloadSize())
	}

	return nil
}
