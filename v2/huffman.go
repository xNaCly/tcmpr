/*
tcmpr2 is a huffman coding lossless compression algorithm (see [1]).

tcmpr compressed blocks of data use the 0x74, 0x0 magic number (t). The
resulting format is represented as follows:

	<magic number>
	<list of frequency keys>
	<list frequency values>
	<length of original input>
	<encoded data>

[1]: https://en.wikipedia.org/wiki/Huffman_coding
*/
package v2

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/xnacly/tcmpr/v2/bitreader"
	"github.com/xnacly/tcmpr/v2/bitwriter"
)

// []{'t', 0x0}
var magicNum []byte = []byte{0x74, 0x0}

type huffman struct {
	HasKey    bool
	Key       byte
	Frequency int32
	L         *huffman
	R         *huffman
}

func dfs(table map[byte][]bool, node *huffman, path []bool) {
	if node == nil {
		return
	}

	if node.HasKey {
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

// walk pre computes a table of byte encodings, thus making huffman access O(1)
func (h *huffman) walk() map[byte][]bool {
	table := make(map[byte][]bool, 256)
	dfs(table, h, nil)
	return table
}

// walkWith returns the byte encoded inside the huffman tree by consuming
// exactly as many bits as necessary. Errors on no valid encoding endpoint
// found or with io.EOF if reader is exhausted
func (h *huffman) walkWith(br *bitreader.BitReader) (byte, error) {
	node := h

	for {
		if node == nil {
			return 0x0, errors.New("invalid bit stream, reached nul node before decoding a byte")
		}

		if node.HasKey {
			break
		}

		bit, err := br.ReadBit()
		if err != nil {
			return 0x0, err
		}

		if bit {
			node = node.R
		} else {
			node = node.L
		}
	}

	return node.Key, nil
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
func (f *frequency) compute(r *bufio.Reader) (int32, error) {
	f.M = map[byte]int32{}
	l := int32(0)
	for {
		c, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return 0, err
			}
		}
		l++
		if val, ok := f.M[c]; ok {
			f.M[c] = val + 1
		} else {
			f.M[c] = 1
		}
	}
	return l, nil
}

// tree builds a huffman tree out of the frequency, the root being the return value
func (f *frequency) tree() *huffman {
	p := &prioQueue{}

	// this is needed to make the map access deterministic
	keys := make([]byte, 0, len(f.M))
	for k := range f.M {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		// Primary: lower frequency first
		if f.M[keys[i]] != f.M[keys[j]] {
			return f.M[keys[i]] < f.M[keys[j]]
		}
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		p.push(&huffman{
			HasKey:    true,
			Key:       k,
			Frequency: f.M[k],
		})
	}

	for len(*p) > 1 {
		l := p.pull()
		r := p.pull()

		p.push(&huffman{
			HasKey:    false,
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
	size, err := f.compute(tee)
	if err != nil {
		return nil
	}
	err = f.serialize(w)
	if err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, int32(size)); err != nil {
		return err
	}

	h := f.tree().walk()
	bWriter := bitwriter.New(w)
	for _, b := range b.Bytes() {
		path := h[b]
		bWriter.WriteBits(path)
	}

	return bWriter.Flush()
}

func debug[T any](t T) T {
	out, _ := json.MarshalIndent(t, "", "\t")
	fmt.Println(string(out))
	return t
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

	var originalLength int32
	if err := binary.Read(buf, binary.BigEndian, &originalLength); err != nil {
		return err
	}

	if originalLength == 0 {
		_, err := w.Write([]byte{})
		return err
	}

	tree := f.tree()
	br := bitreader.New(buf)
	byteBuffer := bytes.Buffer{}
	for range originalLength {
		b, err := tree.walkWith(br)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		byteBuffer.WriteByte(b)
	}

	actualWriteSize, err := byteBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	if actualWriteSize != int64(originalLength) {
		return errors.New("Failed to write the whole decompressed buffer to the writer")
	}

	return nil
}
