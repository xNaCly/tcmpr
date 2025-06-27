package v2

// TODO: compare with other compressions via the communist manifest
// (https://www.gutenberg.org/cache/epub/61/pg61.txt)

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func tDisplay(tree *huffman) string {
	out, _ := json.MarshalIndent(tree, "", "\t")
	return string(out)
}

func TestPriorityQueue(t *testing.T) {
	p := prioQueue{}
	p.push(&huffman{Key: 0x0, Frequency: 0x12})
	p.push(&huffman{Key: 0x1, Frequency: 0x3})
	p.push(&huffman{Key: 0xA, Frequency: 0x25})
	assert.Len(t, p, 3)
	h := p.pull()
	assert.Equal(t, h.Key, byte(0x1))
	assert.Equal(t, h.Frequency, byte(0x3))
	assert.Len(t, p, 2)
	h = p.pull()
	assert.Equal(t, h.Key, byte(0x0))
	assert.Equal(t, h.Frequency, byte(0x12))
	assert.Len(t, p, 1)
	h = p.pull()
	assert.Equal(t, h.Key, byte(0xA))
	assert.Equal(t, h.Frequency, byte(0x25))
	assert.Len(t, p, 0)
}

func TestFrequency(t *testing.T) {
	in := bufio.NewReader(strings.NewReader("BCAADDDCCACACAC"))
	f := frequency{}
	assert.NoError(t, f.compute(in))
	m := make(map[byte]int32, len(f.M))
	for k, v := range f.M {
		m[k] = v
	}
	fmt.Println(m)
	buf := &bytes.Buffer{}
	assert.NoError(t, f.serialize(buf))
	assert.NoError(t, f.deserialize(bufio.NewReader(buf)))
	assert.EqualValues(t, m, f.M)
}

func TestTree(t *testing.T) {
	in := bufio.NewReader(strings.NewReader("BCAADDDCCACACAC"))
	f := frequency{}
	assert.NoError(t, f.compute(in))
	tree := f.tree()
	fmt.Println(tDisplay(tree))
}

func TestTreeTable(t *testing.T) {
	in := bufio.NewReader(strings.NewReader("BCAADDDCCACACAC"))
	f := frequency{}
	assert.NoError(t, f.compute(in))
	table := f.tree().walk()
	fmt.Printf("%#+v\n", table)
}

func TestHuffman(t *testing.T) {
	input := []string{
		"ABC",
		"BCAADDDCCACACAC",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		// strings.Repeat("1111111111111111111111111111111111", 256),
	}

	for _, test := range input {
		inBuf := strings.NewReader(test)
		outBuf := bytes.Buffer{}
		err := Compress(inBuf, &outBuf)
		assert.NoError(t, err, "Failed to compress buffer")
		fmt.Printf("len in: %d, len out: %d\n", len([]byte(test)), outBuf.Len())

		outBuf2 := bytes.Buffer{}
		err = Decompress(&outBuf, &outBuf2)
		assert.NoError(t, err)
		// assert.Equal(t, []byte(test), outBuf2.Bytes())
	}
}
