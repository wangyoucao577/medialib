package rtmp

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/golang/glog"
)

func (h *Handler) handshark() error {

	// handshark states
	const (
		stateUninitialized = "uninitialized"
		stateVersionSent   = "version_sent"
		stateAckSent       = "ack_sent"
		stateHandsharkDone = "handshark_done"
	)
	state := stateUninitialized
	glog.V(1).Infof("state %s", state)

	// send c0+c1
	const c0Version = 0x3 // fixed rtmp version 3
	c0 := []byte{c0Version}
	timestamp := time.Now().UnixMilli()
	timestampBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(timestampBytes, uint32(timestamp))
	const c1RandBytes = 1528
	c1RandomData := make([]byte, c1RandBytes)
	if n, err := rand.Read(c1RandomData); err != nil {
		return err
	} else if n != c1RandBytes {
		glog.Warningf("rand bytes error, expect %d but got %d bytes", c1RandBytes, n)
	}
	c1 := append(timestampBytes, []byte{0x00, 0x00, 0x00, 0x00}...) // timestamp + 4 zero bytes
	c1 = append(c1, c1RandomData...)                                // + random 1528 bytes
	c0c1Len := len(c0) + len(c1)                                    // should be 1+1536

	glog.V(1).Infof("send c0 version 0x%x, c1 timestamp %d", c0Version, timestamp)
	if n, err := h.conn.Write(append(c0, c1...)); err != nil {
		return err
	} else if n != len(c0)+len(c1) {
		return fmt.Errorf("connection write %d bytes but expect %d", n, len(c0)+len(c1))
	}
	state = stateVersionSent
	glog.V(1).Infof("state %s", state)

	// recv s0+s1
	s0s1 := make([]byte, c0c1Len)
	var readS0S1Bytes int
	for {
		if readS0S1Bytes >= c0c1Len {
			break
		}

		if n, err := h.conn.Read(s0s1[readS0S1Bytes:]); err != nil {
			return err
		} else {
			readS0S1Bytes += n
		}
	}
	s0s1ReadTimestamp := time.Now().UnixMilli()
	s0s1ReadTimestampBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(s0s1ReadTimestampBytes, uint32(s0s1ReadTimestamp))

	// check s0,s1
	s0Version := s0s1[0]
	s1Timestamp := binary.BigEndian.Uint32(s0s1[1:5])
	s2Randomdata := s0s1[9:]
	glog.V(1).Infof("recved s0 version 0x%x, s1 timestamp %d, s0+s1 read timestamp %d", s0Version, s1Timestamp, s0s1ReadTimestamp)
	if s0Version != c0Version {
		return fmt.Errorf("unsupported server rtmp version 0x%x", s0Version)
	}

	// send c2
	c2 := []byte{}
	c2 = append(c2, s0s1[1:5]...)              // append s1 timestamp bytes
	c2 = append(c2, s0s1ReadTimestampBytes...) // append timestamp when s1 read
	c2 = append(c2, s2Randomdata...)           // append s1 random data
	glog.V(1).Infof("send c2 timestamp %d timestamp2 %d", s1Timestamp, s0s1ReadTimestamp)
	if n, err := h.conn.Write(c2); err != nil {
		return err
	} else if n != len(c2) {
		return fmt.Errorf("connection write %d bytes but expect %d", n, len(c2))
	}
	state = stateAckSent
	glog.V(1).Infof("state %s", state)

	// recv s2
	s2Len := len(c2)
	s2 := make([]byte, s2Len)
	var readS2Bytes int
	for {
		if readS2Bytes >= s2Len {
			break
		}

		if n, err := h.conn.Read(s2[readS2Bytes:]); err != nil {
			return err
		} else {
			readS2Bytes += n
		}
	}

	// check s2
	s2Timestamp := binary.BigEndian.Uint32(s2[:4])
	s2Timestamp2 := binary.BigEndian.Uint32(s2[4:8])
	s2RandomData := s2[8:]
	glog.V(1).Infof("recved s2 timestamp %d timestamp2 %d", s2Timestamp, s2Timestamp2)
	if !bytes.Equal(c1RandomData, s2RandomData) { // don't need check strictly
		glog.Warningf("c1 and s2 random data not equal")
	}
	state = stateHandsharkDone
	glog.V(1).Infof("state %s", state)

	return nil
}
