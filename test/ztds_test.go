package test

import (
	"github.com/klauspost/compress/zstd"
	"io"
	"os"
	"testing"
)

func Test_CreateZstd(t *testing.T) {
	out, err := os.OpenFile("../testdata/demo.tar.zstd", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open write->%+v", err)
	}
	defer out.Close()

	in, err := os.OpenFile("../testdata/1.txt", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open read->%+v", err)
	}
	defer in.Close()

	w, err := zstd.NewWriter(out)
	if err != nil {
		t.Fatalf("failed to open zstd writer->%+v", err)
	}
	defer w.Close()

	_, err = io.Copy(w, in)
	if err != nil {
		t.Fatalf("failed to write zstd->%+v", err)
	}
}
