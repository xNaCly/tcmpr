package bitreader

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	br := bytes.NewReader([]byte{0b10101010})
	expected := []bool{true, false, true, false, true, false, true, false}
	r := New(br)
	for _, b := range expected {
		bit, err := r.ReadBit()
		if err != nil {
			assert.ErrorIs(t, err, io.EOF)
			break
		}
		assert.Equal(t, bit, b)
	}
}
