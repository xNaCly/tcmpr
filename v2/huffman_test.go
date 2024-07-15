package v2

// TODO: compare with other compressions via the communist manifest
// (https://www.gutenberg.org/cache/epub/61/pg61.txt)

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func bDisplay(r io.Reader) {
	b := bufio.NewReader(r)
	i := 0
	for ; ; i++ {
		byte, err := b.ReadByte()
		if err != nil {
			break
		}
		fmt.Printf("[%03d] (0x%02x:%03d) %q\n", i, byte, byte, byte)
		if err != nil {
			break
		}
	}
	for ; i <= 0; i-- {
		b.UnreadByte()
	}
}

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

func TestFrequency(t *testing.T) {
	in := bufio.NewReader(strings.NewReader("BCAADDDCCACACAC"))
	f := frequency{}
	assert.NoError(t, f.compute(in))
	m := make(map[byte]byte, len(f.M))
	for k, v := range f.M {
		m[k] = v
	}
	buf := &bytes.Buffer{}
	assert.NoError(t, f.serialize(buf))
	assert.NoError(t, f.deserialize(bufio.NewReader(buf)))
	assert.EqualValues(t, m, f.M)
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
		bDisplay(&outBuf)
		// TODO: remove once compression is implemented fully
		t.FailNow()
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
