package bitwriter

import "io"

type BitWriter struct {
	inner          io.Writer
	curByteBuilder byte
	bitCount       int
}

func New(w io.Writer) BitWriter {
	return BitWriter{
		inner:          w,
		curByteBuilder: 0,
		bitCount:       0,
	}
}

func (w *BitWriter) WriteBit(bit bool) error {
	if w.bitCount == 8 {
		if _, err := w.inner.Write([]byte{w.curByteBuilder}); err != nil {
			return err
		}
		w.bitCount = 0
		w.curByteBuilder = 0
	}

	if bit {
		w.curByteBuilder |= 1 << (7 - w.bitCount)
	}
	w.bitCount++
	return nil
}

func (w *BitWriter) WriteBits(bits []bool) error {
	for _, b := range bits {
		if err := w.WriteBit(b); err != nil {
			return err
		}
	}
	return nil
}

func (w *BitWriter) Flush() error {
	if w.bitCount > 0 {
		for w.bitCount < 8 {
			w.curByteBuilder <<= 1
			w.bitCount++
		}
		if _, err := w.inner.Write([]byte{w.curByteBuilder}); err != nil {
			return err
		}
	}
	w.curByteBuilder = 0
	w.bitCount = 0
	return nil
}
