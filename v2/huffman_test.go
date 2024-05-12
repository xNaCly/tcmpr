package v2

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
