// tcmpr2 is a huffman conding lossless compression algorithm (see [1]). tcmpr
// compressed blocks of data use the 0x74, 0x63, 0x6d, 0x70, 0x72, 0x32, 0x0A
// magic number (tcmpr2).
//
// [1]: https://en.wikipedia.org/wiki/Huffman_coding
package v2

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

var magicNum = [...]byte{0x74, 0x63, 0x6d, 0x70, 0x72, 0x32, 0x0A}

type huffman struct {
	// TODO: think about serializing the huffmann coding tree and prefixing the serialized bytes with it.
	key byte
	l   *huffman
	r   *huffman
}

func createTree(r *bufio.Reader) (*huffman, error) {
	b := &bytes.Buffer{}
	tee := bufio.NewReader(io.TeeReader(r, b))
	freq := map[byte]int{}
	for {
		c, err := tee.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, err
			}
		}
		if val, ok := freq[c]; ok {
			freq[c] = val + 1
		} else {
			freq[c] = 1
		}
	}
	fmt.Println(freq)
	return &huffman{}, nil
}

func Compress(r io.Reader, w io.Writer) error {
	h, err := createTree(bufio.NewReader(r))
	if err != nil {
		return err
	}
	fmt.Println(h)
	panic("Not implemented")
}

func Decompress(r io.Reader, w io.Writer) error {
	panic("Not implemented")
}
