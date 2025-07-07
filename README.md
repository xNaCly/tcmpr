# Tcmpr

Tcmpr contains two compression algorithm approaches:

1. v1 is a run-length lossless compression algorithm (see [1]). It works by
   merging similar blocks of data, thus reducing the total size of the buffer.
   Tcmpr is best employed for data containing larger patterns, such as images or
   byte streams. Tcmpr is a merger and an abbreviation of teo and compression.

2. v2 is a huffman coding lossless compression algorithm (see [2]). It works by
   computing frequencies of bytes in the input and packing them efficiently

[1]: https://en.wikipedia.org/wiki/Run-length_encoding
[2]: https://en.wikipedia.org/wiki/Huffman_coding

## Usage

1. Import tcmpr into your go project:

```shell
$ go get github.com/xnacly/tcmpr/v2
$ go mod tidy
```

2. Compress and Decompress:

```go
package main

import (
	"bytes"
	"os"

	v2 "github.com/xnacly/tcmpr/v2"
)

func main() {
	textThatShouldBeCompressed := strings.Repeat("HelloWorld", 1024)
	input := strings.NewReader(textThatShouldBeCompressed)
	compressedBuffer := bytes.Buffer{}
	err := v2.Compress(input, &compressedBuffer)
	if err != nil {
		panic("failed to compress input: " + err.Error())
	}

	fmt.Printf("compressed=%d;original=%d;ratio=%f\n",
		compressedBuffer.Len(),
		len(textThatShouldBeCompressed),
		float32(len(textThatShouldBeCompressed))/float32(compressedBuffer.Len()),
	)

	decompressedBuffer := bytes.Buffer{}
	err = v2.Decompress(&compressedBuffer, &decompressedBuffer)
	if err != nil {
		panic("failed to decompress buffer: " + err.Error())
	}
}
```

## Benchmarks

> Benchmarks will be available once the tests grow, the project is more mature
> and a significant data set for meaningful benchmarks is aggregated.

## Api

The current api follows the go way of working with data by accepting an
`io.Reader` and an `io.Writer`. The api will be extended once the core of the
project is implemented. I plan on somewhat mirroring existing compression
packages in the go standard library.

## Versioning

Even though Tcmpr is not yet ready for a v1 release the current code
architecture allows for easy versioning, thus enabling me to break
compatibility at will by simply incrementing the version. This is the case for
the release of the module in the `v1` namespace and its successor: `v2`.
