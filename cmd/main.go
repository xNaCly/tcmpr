package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"unicode"

	"log/slog"

	v1 "github.com/xnacly/tcmpr/v1"
	v2 "github.com/xnacly/tcmpr/v2"
)

func must[T any](data T, err error) T {
	if err != nil {
		slog.Error("panic", "err", err)
		os.Exit(1)
	}
	return data
}

func annotatedHexdump(buf *bytes.Buffer) {
	data := buf.Bytes()

	// Header: magic number (2 bytes) + length of keys (1 byte)
	if len(data) < 3 {
		fmt.Println("Buffer too small")
		return
	}
	magic := data[:2]
	fmt.Printf("magic=%02x %02x\n", magic[0], magic[1])

	keyCount := int(data[2])
	fmt.Printf("keys=%d\n", keyCount)

	// Keys start at offset 3
	keysStart := 3
	keysEnd := keysStart + keyCount
	if keysEnd > len(data) {
		fmt.Println("Data too short for keys")
		return
	}
	keys := data[keysStart:keysEnd]
	freqsStart := keysEnd
	freqsEnd := freqsStart + 4*keyCount
	if freqsEnd > len(data) {
		fmt.Println("Data too short for frequencies")
		return
	}
	freqs := data[freqsStart:freqsEnd]
	for i, k := range keys {
		freq := binary.BigEndian.Uint32(freqs[i*4 : i*4+4])
		fmt.Printf("[%c|0x%02x]: %d\n", k, k, freq)
	}

	// The rest is compressed data
	compStart := freqsEnd
	compData := data[compStart:]
	fmt.Printf("%d bytes compressed:\n", len(compData))

	const bytesPerLine = 16
	for i := 0; i < len(compData); i += bytesPerLine {
		end := i + bytesPerLine
		if end > len(compData) {
			end = len(compData)
		}
		line := compData[i:end]

		fmt.Printf("%08x  ", compStart+i)

		// hex bytes
		for j := 0; j < bytesPerLine; j++ {
			if j < len(line) {
				fmt.Printf("%02x ", line[j])
			} else {
				fmt.Printf("   ")
			}
			if j == 7 {
				fmt.Printf(" ")
			}
		}

		// ascii representation
		fmt.Printf(" |")
		for _, b := range line {
			if unicode.IsPrint(rune(b)) {
				fmt.Printf("%c", b)
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("|\n")
	}
}

func main() {
	wantsV1 := flag.Bool("v1", false, "enable v1 compression")
	hexdump := flag.Bool("hex", false, "dump hex of input")
	flag.Parse()
	in := flag.Arg(0)
	file := must(os.Open(in))
	if *wantsV1 {
		out := must(os.Create(in + ".tv1"))
		if err := v1.Compress(file, out); err != nil {
			slog.Error("panic", "err", err)
			os.Exit(1)
		}
	} else {
		if *hexdump {
			b := bytes.Buffer{}
			must(b.ReadFrom(file))
			annotatedHexdump(&b)
		} else {
			out := must(os.Create(in + ".tv2"))
			if err := v2.Compress(file, out); err != nil {
				slog.Error("panic", "err", err)
				os.Exit(1)
			}
		}
	}
}
