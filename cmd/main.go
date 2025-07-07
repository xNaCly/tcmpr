package main

import (
	"flag"
	"os"

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

func main() {
	wantsV1 := flag.Bool("v1", false, "enable v1 compression")
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
		out := must(os.Create(in + ".tv2"))
		if err := v2.Compress(file, out); err != nil {
			slog.Error("panic", "err", err)
			os.Exit(1)
		}
	}
}
