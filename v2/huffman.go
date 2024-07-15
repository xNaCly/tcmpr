/*
tcmpr2 is a huffman coding lossless compression algorithm (see [1]).

tcmpr compressed blocks of data use the 0x74, 0x0A magic number (t). The
resulting format is represented as follows:

	<magic number><map frequency keys>0xA<map frequency values>0xA<encoded data>

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

var magicNum = [...]byte{0x74, 0x0A}

type huffman struct {
	key       byte
	frequency byte
	l         *huffman
	r         *huffman
}

// walk walks the tree until b is found, returns the encoded value for b and true if found
func (h *huffman) walk(b byte) (byte, bool) {
	// TODO:
	return 0x0, false
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

// serialize the frequency map into w, keys as a list of bytes, values as a
// list of bytes, the separator between the list of keys/bytes and their
// values/occurences, is the 0xA byte
func (f *frequency) serialize(w io.Writer) error {
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
	_, err = w.Write([]byte{0xA})
	if err != nil {
		return err
	}
	_, err = w.Write(values)
	if err != nil {
		return err
	}
	return nil
}

// deserialize reads bytes in form exported by serialize and forms a frequency struct
func (f *frequency) deserialize(r *bufio.Reader) error {
	f.M = make(map[byte]byte, 64)
	// preallocation to 64, idk if thats a good value
	keys := make([]byte, 0, 64)
	values := make([]byte, 0, 64)
	workingKeys := true
	var err error
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if b == 0xA {
			workingKeys = false
			continue
		}

		if workingKeys {
			keys = append(keys, b)
		} else {
			values = append(values, b)
		}
	}
	if len(keys) != len(values) {
		return errors.New("key and value list not equally sized")
	}
	for i := 0; i < len(keys); i++ {
		f.M[keys[i]] = values[i]
	}
	return err
}

// computes the frequency map from a list of bytes
func (f *frequency) compute(r *bufio.Reader) error {
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
	err := f.compute(tee)
	if err != nil {
		return nil
	}
	err = f.serialize(w)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{0xA})
	if err != nil {
		return err
	}
	h, err := createTree(f.M)
	if err != nil {
		return err
	}
	fmt.Println(h)
	// panic("Not implemented")
	return nil
}

func Decompress(r io.Reader, w io.Writer) error {
	// TODO: check magic num
	// TODO: compute frequency map from byte array
	// TODO: decode data using the tree
	panic("Not implemented")
}
