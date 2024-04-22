package tcmpr

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func debugTest(in *bytes.Buffer) {
	b := make([]byte, in.Len())
	_, err := in.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#+v\n", b)
}

func TestSimple(t *testing.T) {
	input := []struct {
		in  string
		exp string
	}{
		// taken from: https://en.wikipedia.org/wiki/Run-length_encoding#Example
		{
			in:  "WWWWWWWWWWWWBWWWWWWWWWWWWBBBWWWWWWWWWWWWWWWWWWWWWWWWBWWWWWWWWWWWWWW",
			exp: "12W1B12W3B24W1B14W",
		},
	}
	for _, test := range input {
		inBuf := strings.NewReader(test.in)
		outBuf := bytes.Buffer{}
		err := Compress(inBuf, &outBuf)
		if err != nil {
			debugTest(&outBuf)
			t.Errorf("Failed to compress buffer: %s", err)
		}
		outBuf2 := bytes.Buffer{}
		err = Decompress(&outBuf, &outBuf2)
		if err != nil {
			debugTest(&outBuf)
			t.Errorf("Failed to decompress buffer: %s", err)
		}
		if outBuf2.String() != test.exp {
			t.Errorf("Expected %q, got %q", test.exp, outBuf2.String())
		}
	}
}
