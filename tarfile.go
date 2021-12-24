package godeb

import (
	"archive/tar"
	"errors"
	"fmt"
	"github.com/dsnet/compress/bzip2"
	"github.com/kjk/lzma"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/xi2/xz"
	"io"
)

// DecompressFunc 解压函数
type DecompressFunc func(io.Reader) (io.ReadCloser, error)

func gzipNewReader(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}

func xzNewReader(r io.Reader) (io.ReadCloser, error) {
	r2, err := xz.NewReader(r, 0)
	return io.NopCloser(r2), err
}

func bzipNewReader(r io.Reader) (io.ReadCloser, error) {
	return bzip2.NewReader(r, nil)
}

func lzmaNewReader(r io.Reader) (io.ReadCloser, error) {
	r2 := lzma.NewReader(r)
	return io.NopCloser(r2), nil
}

func lz4NewReader(r io.Reader) (io.ReadCloser, error) {
	r2 := lz4.NewReader(r)
	return io.NopCloser(r2), nil
}

func zstdNewReader(r io.Reader) (io.ReadCloser, error) {
	r2, err := zstd.NewReader(r)
	return io.NopCloser(r2), err
}

var decompressors = map[string]DecompressFunc{
	".gz":   gzipNewReader,
	".xz":   xzNewReader,
	".bz2":  bzipNewReader,
	".lzma": lzmaNewReader,
	".lz4":  lz4NewReader,
	".zst":  zstdNewReader,
}

// Decompress 解压器
func Decompress(ext string) (DecompressFunc, error) {
	if fn, ok := decompressors[ext]; ok {
		return fn, nil
	}

	return nil, errors.New(fmt.Sprintf("not supported Extension:%s", ext))
}

// ArItem ar归档中的条目
type ArItem struct {
	Name   string            `json:"name"`
	Ext    string            `json:"ext"`
	Data   *io.SectionReader `json:"-"`
	Offset int64             `json:"offset"`
	Size   int64             `json:"size"`
}

// TarFile 将ArItem转换成tar实例，以便迭代读取压缩包内容
func (ar *ArItem) TarFile() (*tar.Reader, io.Closer, error) {
	fn, err := Decompress(ar.Ext)
	if err != nil {
		return nil, nil, err
	}
	readCloser, err := fn(ar.Data)
	return tar.NewReader(readCloser), readCloser, err
}
