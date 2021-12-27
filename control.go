package godeb

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// control.tar 包含
// - contorl 元信息
// - md5sums 所有文件的MD5摘要
// - conffiles 标明包文件中哪些是配置文件，配置文件在安装过不会被覆盖，除非明确指定。
// - preinst, postinst, prerm, postrm 钩子脚本
// - config 支持debconf技术的脚本
// - shlibs 共享库依赖列表

// Control Package别名
type Control = Package

// Package 包头Control解析, 解析为 <key> <-> <value>字典后转换为Package结构体
// 未识别的字段保存在Opts中, Raw为原始数据
// 该方法保证不出错地且不校验语法地解析出所有字段
// 其<key>和<value>均删除了收尾空格(如果有)
type Package struct {
	Package       string            `json:"package"` //mandatory
	Source        string            `json:"source,omitempty"`
	Version       string            `json:"version"`      //mandatory
	Section       string            `json:"section"`      //recommended
	Priority      string            `json:"priority"`     //recommended
	Architecture  string            `json:"architecture"` //mandatory
	Essential     string            `json:"essential,omitempty"`
	Depends       string            `json:"depends,omitempty"`
	Recommends    string            `json:"recommends,omitempty"`
	Suggests      string            `json:"suggests,omitempty"`
	Enhances      string            `json:"enhances,omitempty"`
	PreDepends    string            `json:"pre_depends,omitempty"`
	InstalledSize int64             `json:"installed_size,omitempty"`
	Maintainer    string            `json:"maintainer"`  //mandatory
	Description   string            `json:"description"` //mandatory
	Homepage      string            `json:"homepage,omitempty"`
	Opts          string            `json:"opts,omitempty"`
	Raw           map[string]string `json:"raw,omitempty"`
	RawText       string            `json:"-"`
}

// Parse 解析包头Control
func Parse(r io.Reader) (*Package, error) {
	m := make(map[string]string)
	var bf bytes.Buffer
	r = io.TeeReader(r, &bf)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == ' ' || line[0] == '\t' {
			m["Description"] += "\n" + line
			continue
		}
		lineParts := strings.SplitN(line, ":", 2)
		key := lineParts[0]
		if len(lineParts) <= 1 {
			return nil, errors.New("invalid line:" + line)
		}
		value := strings.TrimSpace(lineParts[1])
		m[key] = value
	}

	pkg := fromMap(m)
	pkg.Raw = m
	pkg.RawText = bf.String()
	return pkg, nil
}

// FromMap convert map to package
func fromMap(m map[string]string) *Package {
	var pkg Package
	v := reflect.ValueOf(&pkg).Elem()
	t := reflect.TypeOf(pkg)
	for i, n := 0, t.NumField(); i < n; i++ {
		tf := t.Field(i)
		vf := v.Field(i)
		if value, ok := m[tf.Name]; ok {
			switch vf.Kind() {
			case reflect.String:
				vf.SetString(value)
			case reflect.Int64:
				i, _ := strconv.ParseInt(value, 10, 64)
				vf.SetInt(i)
			}
		}
	}
	return &pkg
}
