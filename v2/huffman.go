/*
tcmpr2 is a huffman coding lossless compression algorithm (see [1]).

tcmpr compressed blocks of data use the 0x74, 0x63, 0x6d, 0x70, 0x72, 0x32,
0x0A magic number (tcmpr2). The resulting format is represented as follows:

	<magic number><map frequency keys>0x0<map frequency values>0x0<encoded data>

[1]: https://en.wikipedia.org/wiki/Huffman_coding
*/
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
	key       byte
	frequency byte
	l         *huffman
	r         *huffman
}

type prioQueue []*huffman

func (p *prioQueue) push(b *huffman) {
	*p = append(*p, b)
}

func (p *prioQueue) pull() *huffman {
	if len(*p) == 0 {
		return nil
	}
	var lowest *huffman = nil
	var index int = 0
	for i, v := range *p {
		if lowest == nil || v.frequency < lowest.frequency {
			lowest = v
			index = i
		}
	}

	(*p)[index] = (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return lowest
}

type frequency struct {
	M map[byte]byte
}

// Write dumps the frequency map into w, keys as a list of bytes, values as a
// list of bytes, the separator is the 0x0 byte
func (f *frequency) Dump(w io.Writer) error {
	keys := make([]byte, 0, len(f.M))
	values := make([]byte, 0, len(f.M))
	for k, v := range f.M {
		keys = append(keys, k)
		values = append(values, v)
	}
	_, err := w.Write(keys)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{0x0})
	if err != nil {
		return err
	}
	_, err = w.Write(values)
	if err != nil {
		return err
	}
	return nil
}

// TODO: implement this
func (f *frequency) ComputeFromDump(r *bufio.Reader) error {
	return nil
}

// Compute produces the frequency map from a list of bytes
func (f *frequency) Compute(r *bufio.Reader) error {
	f.M = map[byte]byte{}
	for {
		c, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return err
			}
		}
		if val, ok := f.M[c]; ok {
			f.M[c] = val + 1
		} else {
			f.M[c] = 1
		}
	}
	return nil
}

func createTree(freq map[byte]byte) (*huffman, error) {
	p := &prioQueue{}
	for k, v := range freq {
		p.push(&huffman{
			key:       k,
			frequency: v,
		})
	}
	for {
		e := p.pull()
		if e == nil {
			break
		}
	}
	// TODO: say youre going to make a tree and then just tree all over the
	// place
	return &huffman{}, nil
}

func Compress(r io.Reader, w io.Writer) error {
	w.Write(magicNum[:])
	b := &bytes.Buffer{}
	tee := bufio.NewReader(io.TeeReader(r, b))
	f := frequency{}
	err := f.Compute(tee)
	if err != nil {
		return nil
	}
	err = f.Dump(w)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{0x0})
	if err != nil {
		return err
	}
	h, err := createTree(f.M)
	if err != nil {
		return err
	}
	fmt.Println(h)
	panic("Not implemented")
}

func Decompress(r io.Reader, w io.Writer) error {
	// TODO: check magic num
	// TODO: compute frequency map from byte array
	// TODO: decode data using the tree
	panic("Not implemented")
}
