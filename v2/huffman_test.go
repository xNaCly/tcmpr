package v2

// TODO: compare with other compressions via the communist manifest
// (https://www.gutenberg.org/cache/epub/61/pg61.txt)

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tDisplay(tree *huffman) string {
	out, _ := json.MarshalIndent(tree, "", "\t")
	return string(out)
}

func TestPriorityQueue(t *testing.T) {
	p := prioQueue{}
	p.push(&huffman{Key: 0x0, Frequency: 1})
	p.push(&huffman{Key: 0x1, Frequency: 2})
	p.push(&huffman{Key: 0xA, Frequency: 3})
	l := 3
	compare := func(key byte, f int32) {
		h := p.pull()
		assert.Equal(t, h.Key, key)
		assert.Equal(t, h.Frequency, f)
		l -= 1
		assert.Len(t, p, l)
	}
	compare(0x0, 1)
	compare(0x1, 2)
	compare(0xA, 3)
}

func TestFrequency(t *testing.T) {
	in := bufio.NewReader(strings.NewReader("BCAADDDCCACACAC"))
	f := frequency{}
	_, err := f.compute(in)
	assert.NoError(t, err)
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
	_, err := f.compute(in)
	assert.NoError(t, err)
	tree := f.tree()
	fmt.Println(tDisplay(tree))
}

func TestTreeTable(t *testing.T) {
	in := bufio.NewReader(strings.NewReader("BCAADDDCCACACAC"))
	f := frequency{}
	_, err := f.compute(in)
	assert.NoError(t, err)
	table := f.tree().walk()
	fmt.Printf("%#+v\n", table)
}

func TestHuffman(t *testing.T) {
	input := []string{
		"ABC",
		"BCAADDDCCACACAC",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"ABABABABABABABABABABABABABABABABABAB",
		"The quick brown fox jumps over the lazy dog",
		"0123456789!@#$%^&*()_+-=[]{}|;':,.<>/?",
		"",
		"A",
		"漢字テスト日本語",
		strings.Repeat("GoLangIsAwesome!", 100),
	}

	for _, test := range input {
		t.Run(test[:min(50, len(test))], func(t *testing.T) {
			inBuf := strings.NewReader(test)
			outBuf := bytes.Buffer{}
			err := Compress(inBuf, &outBuf)
			assert.NoError(t, err, "Failed to compress buffer")
			outBuf2 := bytes.Buffer{}
			err = Decompress(&outBuf, &outBuf2)
			assert.NoError(t, err)

			got := outBuf2.Bytes()
			expected := []byte(test)
			if got == nil {
				got = []byte{}
			}
			assert.Equal(t, expected, got)
		})
	}
}

func TestReadme(t *testing.T) {
	textThatShouldBeCompressed := strings.Repeat("HelloWorld", 1024)
	input := strings.NewReader(textThatShouldBeCompressed)
	compressedBuffer := bytes.Buffer{}
	err := Compress(input, &compressedBuffer)
	if err != nil {
		panic("failed to compress input: " + err.Error())
	}

	fmt.Printf("compressed=%d;original=%d;ratio=%f\n",
		compressedBuffer.Len(),
		len(textThatShouldBeCompressed),
		float32(len(textThatShouldBeCompressed))/float32(compressedBuffer.Len()),
	)

	decompressedBuffer := bytes.Buffer{}
	err = Decompress(&compressedBuffer, &decompressedBuffer)
	if err != nil {
		panic("failed to decompress buffer: " + err.Error())
	}
}
