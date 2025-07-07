package bitreader

import (
	"io"
)

type BitReader struct {
	inner    io.Reader
	curByte  [1]byte
	bitCount int // amount of bits read already, resets upon hitting 8
}

func New(r io.Reader) *BitReader {
	return &BitReader{
		inner:    r,
		bitCount: 0,
	}
}

func (r *BitReader) ReadBit() (bool, error) {
	if r.bitCount == 0 || r.bitCount == 8 {
		if count, err := r.inner.Read(r.curByte[:]); err != nil {
			return false, err
		} else if count != 1 {
			return false, io.EOF
		}
		r.bitCount = 0
	}

	isSet := (r.curByte[0] & (1 << (7 - r.bitCount))) > 0
	r.bitCount++
	return isSet, nil
}
