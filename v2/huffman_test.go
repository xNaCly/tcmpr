package v2

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {
	p := prioQueue{}
	p.push(&huffman{key: 0x0, frequency: 0x12})
	p.push(&huffman{key: 0x1, frequency: 0x3})
	p.push(&huffman{key: 0xA, frequency: 0x25})
	assert.Len(t, p, 3)
	h := p.pull()
	assert.Equal(t, h.key, byte(0x1))
	assert.Equal(t, h.frequency, byte(0x3))
	assert.Len(t, p, 2)
	h = p.pull()
	assert.Equal(t, h.key, byte(0x0))
	assert.Equal(t, h.frequency, byte(0x12))
	assert.Len(t, p, 1)
	h = p.pull()
	assert.Equal(t, h.key, byte(0xA))
	assert.Equal(t, h.frequency, byte(0x25))
	assert.Len(t, p, 0)
}

func TestHuffman(t *testing.T) {
	input := []struct {
		in  string
		exp string
	}{
		{"BCAADDDCCACACAC", ""},
	}

	for _, test := range input {
		inBuf := strings.NewReader(test.in)
		outBuf := bytes.Buffer{}
		err := Compress(inBuf, &outBuf)
		assert.NoError(t, err, "Failed to compress buffer")

		b := outBuf.Bytes()
		fmt.Printf("exp: %#+v\nout: %#+v\n len in: %d, len out: %d\n", test.exp, b, len([]byte(test.in)), len(b))
		assert.Equal(t, test.exp, outBuf.Bytes())

		outBuf2 := bytes.Buffer{}
		err = Decompress(&outBuf, &outBuf2)
		assert.NoError(t, err, "Failed to decompress buffer")
		assert.Equal(t, []byte(test.in), outBuf2.Bytes())
	}
}
