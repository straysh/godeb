package test

import (
	"bytes"
	"github.com/straysh/godeb"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func Test_DecompressFile(t *testing.T) {
	files, err := filepath.Glob("../testdata/*.tar.*")
	if err != nil {
		t.Fatalf("failed to read dir->%+v", err)
	}

	for _, f := range files {
		t.Run(f, func(t2 *testing.T) {
			r, err := os.Open(f)
			if err != nil {
				t2.Fatalf("failed to read file:%s->%+v", f, err)
			}

			ext := filepath.Ext(f)
			decompressor, err := godeb.Decompress(ext)
			if err != nil {
				t2.Fatalf("failed to get decompressor->%+v", err)
			}
			b, err := decompressor(r)
			if err != nil {
				t2.Fatalf("failed to decompress->%+v", err)
			}

			var buf bytes.Buffer
			_, err = io.Copy(&buf, b)
			if err != nil {
				t2.Fatalf("failed to copy->%+v", err)
			}
			t2.Logf("got:%s", buf.String())
		})
	}
}
