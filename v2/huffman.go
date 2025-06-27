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
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/xnacly/tcmpr/v2/bitwriter"
)

// []{'t', 0x0}
var magicNum []byte = []byte{0x74, 0x0}

type huffman struct {
	hasKey    bool
	Key       byte
	Frequency int32
	L         *huffman
	R         *huffman
}

func dfs(table map[byte][]bool, node *huffman, path []bool) {
	if node == nil {
		return
	}

	if node.hasKey {
		if len(path) == 0 {
			table[node.Key] = []bool{false}
		} else {
			table[node.Key] = append([]bool{}, path...)
		}
		return
	}

	dfs(table, node.L, append(path, false))
	dfs(table, node.R, append(path, true))
}

// walk computes a table of byte encodings, thus making huffman access O(1)
func (h *huffman) walk() map[byte][]bool {
	table := make(map[byte][]bool, 256)
	dfs(table, h, nil)
	return table
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
	M map[byte]int32
}

// serialize the frequency map into w, keys as a list of bytes, values as a
// list of bytes
func (f *frequency) serialize(w io.Writer) error {
	keys := make([]byte, 0, len(f.M))
	values := make([]int32, 0, len(f.M))
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
	for _, v := range values {
		if err := binary.Write(w, binary.BigEndian, v); err != nil {
			return err
		}
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
	f.M = make(map[byte]int32, length)
	keys := make([]byte, 0, length)

	for i := 0; i < length; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		keys = append(keys, b)
	}

	for i := 0; i < length; i++ {
		var b int32
		if err := binary.Read(r, binary.BigEndian, &b); err != nil {
			return err
		}
		f.M[keys[i]] = b
	}

	return nil
}

// computes the frequency map from a list of bytes
func (f *frequency) compute(r *bufio.Reader) error {
	f.M = map[byte]int32{}
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

	for len(*p) > 1 {
		l := p.pull()
		r := p.pull()

		p.push(&huffman{
			hasKey:    false,
			Frequency: l.Frequency + r.Frequency,
			L:         l,
			R:         r,
		})
	}

	return p.pull()
}

// Compresses bytes in r into w
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
	h := f.tree().walk()
	bWriter := bitwriter.New(w)
	for _, b := range b.Bytes() {
		path := h[b]
		bWriter.WriteBits(path)
	}
	bWriter.Flush()
	return nil
}

func Decompress(r io.Reader, w io.Writer) error {
	buf := bufio.NewReader(r)

	shouldBeMagicNum, err := buf.ReadBytes(magicNum[len(magicNum)-1])
	if err != nil {
		return err
	}
	if !bytes.Equal(magicNum[:], shouldBeMagicNum) {
		return fmt.Errorf("Magic number incorrect %#+v vs %#+v", shouldBeMagicNum, magicNum)
	}

	f := frequency{}
	if err := f.deserialize(buf); err != nil {
		return err
	}

	// tree := f.tree()

	return nil
}
