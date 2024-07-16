# Tcmpr

> _WARNING_:
>
> tcmpr is not ready for any use, let alone production.

Tcmpr is a run-length lossless compression algorithm (see [1]). It works by
merging similar blocks of data, thus reducing the total size of the buffer.
Tcmpr is best employed for data containing larger patterns, such as images or
byte streams. Tcmpr is a merger and an abbreviation of teo and compression.

[1]: https://en.wikipedia.org/wiki/Run-length_encoding

## Usage

1. Import tcmpr into your go project:

```shell
$ go get github.com/xnacly/tcmpr/v1
$ go mod tidy
```

2. Compress and Decompress:

```go
package main

import (
	"bytes"
	"os"

	"github.com/xnacly/tcmpr/v1"
)

func main() {
	inputBuffer := bytes.Buffer{}
	inputBuffer.WriteString("Hello World")

	outputBuffer := bytes.Buffer{}
	err := tcmpr.Compress(&inputBuffer, &outputBuffer)
	if err != nil {
		panic(err)
	}
	tcmpr.Decompress(&outputBuffer, os.Stdout)
	if err != nil {
		panic(err)
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
the release of the module in the `v1` namespace.
