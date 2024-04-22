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

func rwErrWrapper(n any, err error) error {
	return err
}

func writeCompressed(w *bufio.Writer, counter int, char byte) {
	// fixed bug: handle edge case of counter being bigger than a byte (c>255): split
	// it up into multiple writes of COUNT BYTE all writing at most 255, thus
	// we would split 614 c as follows: 255c255c104c
	for counter > 255 {
		counter -= 255
		w.WriteByte(255)
		w.WriteByte(char)
	}
	w.WriteByte(byte(counter))
	w.WriteByte(char)
}

// Compress compresses data passed in via r into w
func Compress(r io.Reader, w io.Writer) error {
	bufR := bufio.NewReader(r)
	bufW := bufio.NewWriter(w)
	err := rwErrWrapper(bufW.Write(magicNum[:]))
	if err != nil {
		return err
	}

	lc, err := bufR.ReadByte()
	if err != nil {
		return err
	}

	var counter int = 1
	for {
		b, err := bufR.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// hitting eof should just stop us iterating and we of course
				// need to write the last buffered byte and its counter
				writeCompressed(bufW, counter, lc)
				break
			} else {
				return err
			}
		}

		if b != lc {
			writeCompressed(bufW, counter, lc)
			lc = b
			counter = 1
		} else {
			counter++
		}
	}

	err = rwErrWrapper(bufR.WriteTo(bufW))
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

	err = rwErrWrapper(bufR.WriteTo(bufW))
	if err != nil {
		return err
	}
	return bufW.Flush()
}
