// tcmpr1 is a rudimentary run-length lossless compression algorithm (see [1]).
// It works by merging similar blocks of data, thus reducing the total size of
// the buffer. tcmpr compressed blocks of data use the 0x74, 0x63, 0x6d, 0x70,
// 0x72, 0x31, 0x0A magic number (tcmpr1).
//
// [1]: https://en.wikipedia.org/wiki/Run-length_encoding
package v1

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"slices"
)

var magicNum = [...]byte{0x74, 0x63, 0x6d, 0x70, 0x72, 0x31, 0x0A}

func rwErrWrapper(n any, err error) error {
	return err
}

func writeCompressed(w *bufio.Writer, counter int, char byte) error {
	// fixed bug: handle edge case of counter being bigger than a byte (c>255): split
	// it up into multiple writes of COUNT BYTE all writing at most 255, thus
	// we would split 614 c as follows: 255c255c104c
	if counter > 255 {
		// compute the times we have to write byte max (255)
		var m = counter / 255
		// compute the remaining count of bytes to write
		var r = counter % 255
		// repeat byte repetition with correct byte
		err := rwErrWrapper(w.Write(bytes.Repeat([]byte{255, char}, m)))
		if err != nil {
			return err
		}
		counter = r
	}
	err := w.WriteByte(byte(counter))
	if err != nil {
		return err
	}
	return w.WriteByte(char)
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
				err := writeCompressed(bufW, counter, lc)
				if err != nil {
					return err
				}
				break
			} else {
				return err
			}
		}

		if b != lc {
			err := writeCompressed(bufW, counter, lc)
			if err != nil {
				return err
			}
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
	for {
		c, err := bufR.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		b, err := bufR.ReadByte()
		if err != nil {
			return err
		}
		buf := make([]byte, c)
		for i := byte(0); i < c; i++ {
			buf[i] = b
		}
		bufW.Write(buf)
	}
	return bufW.Flush()
}
