package test

import (
	"github.com/klauspost/compress/gzip"
	"io"
	"os"
	"testing"
)

func Test_CreateGz(t *testing.T) {
	out, err := os.OpenFile("../testdata/demo2.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open write->%+v", err)
	}
	defer out.Close()

	in, err := os.OpenFile("../testdata/1.txt", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open read->%+v", err)
	}
	defer in.Close()

	w := gzip.NewWriter(out)
	defer w.Close()

	_, err = io.Copy(w, in)
	if err != nil {
		t.Fatalf("failed to write zstd->%+v", err)
	}
}
