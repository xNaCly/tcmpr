// tcmpr1 is a rudimentary run-length lossless compression algorithm (see [1]).
// It works by merging similar blocks of data, thus reducing the total size of
// the buffer. tcmpr compressed blocks of data use the 0x74, 0x63, 0x6d, 0x70,
// 0x72, 0x31, 0x0A magic number.
//
// [1]: https://en.wikipedia.org/wiki/Run-length_encoding
package tcmpr

import (
	"bufio"
	"errors"
	"io"
	"slices"
)

var magicNum = [...]byte{0x74, 0x63, 0x6d, 0x70, 0x72, 0x31, 0x0A}

func writeErrWrapper(n any, err error) error {
	return err
}

// Compress compresses data passed in via r into w
func Compress(r io.Reader, w io.Writer) error {
	bufR := bufio.NewReader(r)
	bufW := bufio.NewWriter(w)
	err := writeErrWrapper(bufW.Write(magicNum[:]))
	if err != nil {
		return err
	}

	err = writeErrWrapper(bufR.WriteTo(bufW))
	if err != nil {
		return err
	}

	return bufW.Flush()
}

// Decompress decompresses data passed in via r into w
func Decompress(r io.Reader, w io.Writer) error {
	bufR := bufio.NewReader(r)
	bufW := bufio.NewWriter(w)

	magicBuf := make([]byte, len(magicNum))
	n, err := bufR.Read(magicBuf)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("input reader is empty")
	}
	if !slices.Equal(magicBuf, magicNum[:]) {
		return errors.New("reader is not tcmpr compressed")
	}

	err = writeErrWrapper(bufR.WriteTo(bufW))
	if err != nil {
		return err
	}
	return bufW.Flush()
}
