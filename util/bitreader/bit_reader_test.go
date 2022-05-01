package bitreader

import (
	"bytes"
	"testing"
)

func TestReadBit(t *testing.T) {
	cases := []struct {
		count int
		in    []byte
		out   []byte
	}{
		{count: 8, in: []byte{0xAA}, out: []byte{0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0}},
		{count: 8, in: []byte{0x55}, out: []byte{0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1}},
		{count: 10, in: []byte{0x55, 0xFF}, out: []byte{0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1, 0x1, 0x1}},
	}

	for _, c := range cases {
		r := bytes.NewReader(c.in)
		br := New(r)

		for i := 0; i < c.count; i++ {
			if nextBit, err := br.ReadBit(); err != nil {
				t.Error(err)
			} else if nextBit != c.out[i] {
				t.Errorf("read bit %d expect 0x%x but got 0x%x", i, c.out[i], nextBit)
			}
		}
	}
}

func TestReadBits(t *testing.T) {
	cases := []struct {
		count int
		in    []byte
		out   []byte
	}{
		{count: 1, in: []byte{0xAA}, out: []byte{0x1}},
		{count: 2, in: []byte{0xAA}, out: []byte{0x2}},
		{count: 3, in: []byte{0xAA}, out: []byte{0x5}},
		{count: 4, in: []byte{0xAA}, out: []byte{0xA}},
		{count: 5, in: []byte{0xAA}, out: []byte{0x15}},
		{count: 6, in: []byte{0xAA}, out: []byte{0x2A}},
		{count: 7, in: []byte{0xAA}, out: []byte{0x55}},
		{count: 8, in: []byte{0xAA}, out: []byte{0xAA}},

		{count: 1, in: []byte{0xAA, 0xFF}, out: []byte{0x1}},
		{count: 8, in: []byte{0xAA, 0xFF}, out: []byte{0xAA}},
		{count: 9, in: []byte{0xAA, 0xFF}, out: []byte{0xAA, 0x1}},
	}

	for _, c := range cases {
		r := bytes.NewReader(c.in)
		br := New(r)

		if nextBits, err := br.ReadBits(uint(c.count)); err != nil {
			t.Error(err)
		} else if !bytes.Equal(nextBits, c.out) {
			t.Errorf("read %d bits on %v expect %v but got %v", c.count, c.in, c.out, nextBits)
		}
	}
}

func TestReadBytes(t *testing.T) {
	cases := []struct {
		count int
		in    []byte
		out   []byte
	}{
		{count: 1, in: []byte{0xAA}, out: []byte{0xAA}},
		{count: 4, in: []byte{0x00, 0x00, 0x00, 0xC0, 0x00}, out: []byte{0x00, 0x00, 0x00, 0xC0}},
		{count: 8, in: []byte{0x00, 0x00, 0x00, 0xC0, 0x1, 0x00, 0x00, 0x00}, out: []byte{0x00, 0x00, 0x00, 0xC0, 0x1, 0x00, 0x00, 0x00}},
	}

	for _, c := range cases {
		r := bytes.NewReader(c.in)
		br := New(r)

		if nextBytes, err := br.ReadBytes(uint(c.count)); err != nil {
			t.Error(err)
		} else if !bytes.Equal(nextBytes, c.out) {
			t.Errorf("read bytes expect count %d value %v but got %v", c.count, c.out, nextBytes)
		}
	}
}
