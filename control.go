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
	Package       string            `json:"package" field:"Package"` //mandatory
	Source        string            `json:"source,omitempty" field:"Source"`
	Version       string            `json:"version" field:"Version"`           //mandatory
	Section       string            `json:"section" field:"Section"`           //recommended
	Priority      string            `json:"priority" field:"Priority"`         //recommended
	Architecture  string            `json:"architecture" field:"Architecture"` //mandatory
	Essential     string            `json:"essential,omitempty" field:"Essential"`
	Depends       string            `json:"depends,omitempty" field:"Depends"`
	Recommends    string            `json:"recommends,omitempty" field:"Recommends"`
	Suggests      string            `json:"suggests,omitempty" field:"Suggests"`
	Enhances      string            `json:"enhances,omitempty" field:"Enhances"`
	PreDepends    string            `json:"pre_depends,omitempty" field:"Pre-Depends"`
	InstalledSize int64             `json:"installed_size,omitempty" field:"Installed-Size"`
	Maintainer    string            `json:"maintainer" field:"Maintainer"`   //mandatory
	Description   string            `json:"description" field:"Description"` //mandatory
	Homepage      string            `json:"homepage,omitempty" field:"Homepage"`
	Opts          string            `json:"opts,omitempty" field:"-"`
	Raw           map[string]string `json:"raw,omitempty" field:"-"`
	RawText       string            `json:"-" field:"-"`
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
		//tf := t.Field(i)
		vf := v.Field(i)
		tag := v.Type().Field(i).Tag
		field := tag.Get("field")
		if value, ok := m[field]; ok {
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
