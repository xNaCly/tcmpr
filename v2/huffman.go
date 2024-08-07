/*
tcmpr2 is a huffman coding lossless compression algorithm (see [1]).

tcmpr compressed blocks of data use the 0x74, 0x0 magic number (t). The
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

var magicNum = [...]byte{0x74, 0x0}

type huffman struct {
	hasKey    bool
	Key       byte
	Frequency byte
	L         *huffman
	R         *huffman
}

// walk walks the tree until b with f is found, returns the encoded value for b
// and true if found
func walk(tree *huffman) [256]int8 {
	cache := [256]int8{}
	// TODO:
	return cache
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
		if lowest == nil || v.Frequency < lowest.Frequency {
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
	_, err := w.Write([]byte{byte(len(keys))})
	if err != nil {
		return err
	}
	_, err = w.Write(keys)
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
	lengthRaw, err := r.ReadByte()
	if err != nil {
		return err
	}
	length := int(lengthRaw)
	f.M = make(map[byte]byte, 64)
	keys := make([]byte, 0, 64)
	values := make([]byte, 0, 64)

	for i := 0; i < length; i++ {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		keys = append(keys, b)
	}
	for i := 0; i < length; i++ {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		values = append(values, b)
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

func (f *frequency) tree() *huffman {
	p := &prioQueue{}
	for k, v := range f.M {
		p.push(&huffman{
			hasKey:    true,
			Key:       k,
			Frequency: v,
		})
	}
	var root *huffman
	for {
		if len(*p) == 1 {
			root = p.pull()
		}
		l := p.pull()
		if l == nil {
			break
		}
		r := p.pull()
		if r == nil {
			break
		}
		p.push(&huffman{
			Frequency: l.Frequency + r.Frequency,
			L:         l,
			R:         r,
		})
	}
	return root
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
	h := f.tree()
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
