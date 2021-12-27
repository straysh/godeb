package godeb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/blakesmith/ar"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Debian包是标准的Unix ar格式归档，其含两个tar包——一个保存control信息；另一个含有安装文件。
// Debian包内容按顺序如下：
// 1. `debian-binary`: 只有一行文本，记录了包格式版本号。现行版本号为`2.0`
// 2. `control`归档: 名称为`control.tar`的归档文件，包含脚本以及元信息(包名，版本号，依赖，维护人等)。使用`gzip`或`xz`格式压缩。其后缀名表明了压缩格式。
// 3. `data`归档: 名为`data.tar`的归档文件，包含实际安装的文件。使用`gzip`、`bzip2`、`lzma`或`xz`压缩格式。其后缀名表明了压缩格式。

// Deb deb包结构
type Deb struct {
	Control    *Control
	Data       *ArItem
	ControlExt string `json:"control_ext"`
	DataExt    string `json:"data_ext"`
	ArItem     map[string]*ArItem
}

type Reader interface {
	io.Reader
	io.ReaderAt
}

// LoadDeb 加载deb包，返回*Deb结构
func LoadDeb(in Reader) (*Deb, error) {
	contents, err := LoadAr(in)
	if err != nil {
		return nil, err
	}

	deb := &Deb{
		ArItem: contents,
	}
	for name, aritem := range contents {
		switch {
		case strings.HasPrefix(name, "control."):
			deb.ControlExt = filepath.Ext(name)
			deb.Control, err = loadDebControl(aritem)
			if err != nil {
				return nil, err
			}
		case strings.HasPrefix(name, "data."):
			deb.DataExt = filepath.Ext(name)
			deb.Data = aritem
		}
	}

	return deb, nil
}

// LoadAr 加载ar归档返回map[string]*ArItem
func LoadAr(in Reader) (map[string]*ArItem, error) {
	contents := make(map[string]*ArItem)
	arReader := ar.NewReader(in)
	offset := int64(8)
	for {
		header, err := arReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		header.Name = strings.TrimSuffix(strings.TrimSpace(header.Name), "/")
		contents[header.Name] = &ArItem{
			Name:   header.Name,
			Ext:    filepath.Ext(header.Name),
			Data:   io.NewSectionReader(in, offset+ar.HEADER_BYTE_SIZE, header.Size),
			Offset: offset,
			Size:   header.Size,
		}
		offset = offset + ar.HEADER_BYTE_SIZE + header.Size + header.Size%2
	}

	aritem, ok := contents["debian-binary"]
	if !ok {
		return nil, errors.New("package invalid: missing debina-binary")
	}
	versionBytes, err := ioutil.ReadAll(aritem.Data)
	if err != nil {
		return nil, errors.New("failed to read aritem")
	}
	version := string(bytes.TrimSpace(versionBytes))
	if version != "2.0" {
		return nil, errors.New(fmt.Sprintf("unsupported debian-binary:%s", version))
	}

	return contents, nil
}

func loadDebControl(aritem *ArItem) (*Control, error) {
	r, closer, err := aritem.TarFile()
	if err != nil {
		return nil, fmt.Errorf("failed to init tarfile->%+v", err)
	}
	defer closer.Close()

	for {
		member, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read member->%+v", err)
		}

		name := strings.TrimLeft(member.Name, "./")
		if name == "control" {
			var control *Control
			control, err = Parse(r)
			if err != nil {
				return nil, err
			}
			return control, nil
		}
	}
	return nil, nil
}
