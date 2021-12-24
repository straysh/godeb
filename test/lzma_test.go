package test

import (
	"github.com/ulikunitz/xz/lzma"
	//"github.com/pierrec/lz4/v4"
	//"github.com/kjk/lzma"
	"io"
	"os"
	"testing"
)

func Test_CreateLzma(t *testing.T) {
	out, err := os.OpenFile("../testdata/demo3.tar.lzma", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open write->%+v", err)
	}
	defer out.Close()

	in, err := os.OpenFile("../testdata/1.txt", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open read->%+v", err)
	}
	defer in.Close()

	w, _ := lzma.NewWriter(out)
	defer w.Close()

	_, err = io.Copy(w, in)
	if err != nil {
		t.Fatalf("failed to write zstd->%+v", err)
	}
}
