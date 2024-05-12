package v1

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func debugTest(in *bytes.Buffer) {
	b := make([]byte, in.Len())
	_, err := in.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#+v\n", b)
}

func TestRLE(t *testing.T) {
	input := []struct {
		in  string
		exp []byte
	}{
		{
			// taken from: https://en.wikipedia.org/wiki/Run-length_encoding#Example
			in: "WWWWWWWWWWWWBWWWWWWWWWWWWBBBWWWWWWWWWWWWWWWWWWWWWWWWBWWWWWWWWWWWWWW",
			// "12W1B12W3B24W1B14W"
			exp: append(magicNum[:], []byte{12, 'W', 1, 'B', 12, 'W', 3, 'B', 24, 'W', 1, 'B', 14, 'W'}...),
		},
		{
			in:  "00000000000000000001111111111111111111111222222222222222222",
			exp: append(magicNum[:], []byte{19, '0', 22, '1', 18, '2'}...),
		},
		// TODO: test with https://netpbm.sourceforge.net/doc/pgm.html and other image formats
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

func TestEdgeRLE(t *testing.T) {
	input := []struct {
		name string
		in   string
		exp  []byte
	}{
		{
			name: "counter for byte larger than 255",
			in:   strings.Repeat("c", 614),
			exp:  append(magicNum[:], []byte{255, 'c', 255, 'c', 104, 'c'}...),
		},
	}
	for _, test := range input {
		inBuf := strings.NewReader(test.in)
		outBuf := bytes.Buffer{}
		err := Compress(inBuf, &outBuf)
		assert.NoError(t, err, "Failed to compress buffer")

		assert.Equal(t, test.exp, outBuf.Bytes())

		outBuf2 := bytes.Buffer{}
		err = Decompress(&outBuf, &outBuf2)
		assert.NoError(t, err, "Failed to decompress buffer")
		assert.Equal(t, []byte(test.in), outBuf2.Bytes())
	}

}
