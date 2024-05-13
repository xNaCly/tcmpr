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

type huffmanNode struct {
	k, f byte
}

func (h *huffmanNode) String() string {
	return fmt.Sprintf("0x%x [%c]: %d", h.k, h.k, h.f)
}

type prioQueue []*huffmanNode

func (p *prioQueue) push(b *huffmanNode) {
	*p = append(*p, b)
}

func (p *prioQueue) pull() *huffmanNode {
	if len(*p) == 0 {
		return nil
	}
	var lowest *huffmanNode = nil
	var index int = 0
	for i, v := range *p {
		if lowest == nil || v.f < lowest.f {
			lowest = v
			index = i
		}
	}

	(*p)[index] = (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return lowest
}

type huffman struct {
	key       byte
	frequency byte
	l         *huffman
	r         *huffman
}

func createTree(r *bufio.Reader) (*huffman, error) {
	b := &bytes.Buffer{}
	tee := bufio.NewReader(io.TeeReader(r, b))
	freq := map[byte]byte{}
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
	p := &prioQueue{}
	for k, v := range freq {
		p.push(&huffmanNode{k, v})
	}
	for {
		e := p.pull()
		if e == nil {
			break
		}
		fmt.Println(e)
	}
	return &huffman{}, nil
}

func Compress(r io.Reader, w io.Writer) error {
	rr := bufio.NewReader(r)
	h, err := createTree(rr)
	if err != nil {
		return err
	}
	fmt.Println(h)
	panic("Not implemented")
}

func Decompress(r io.Reader, w io.Writer) error {
	panic("Not implemented")
}
