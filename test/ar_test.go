package test

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"github.com/straysh/godeb"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_PartialDebControl(t *testing.T) {
	debfile := "../testdata/com.deepin.appstore.helloworld_amd64_5.6.8-1.deb"
	in, err := os.Open(debfile)
	if err != nil {
		t.Fatalf("failed to open debfile->%+v", err)
	}

	b2 := make([]byte, 1024)
	_, err = in.Read(b2)
	if err != nil {
		t.Fatalf("failed to read section->%+v", err)
	}
	in2 := bytes.NewReader(b2)

	deb, err := godeb.LoadDeb(in2)
	if err != nil {
		t.Fatalf("failed to LoadDeb->%+v", err)
	}
	t.Logf("control_ext=%s", deb.ControlExt)
	t.Logf("data_ext=%s", deb.DataExt)

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "  ")
	jsonEncoder.Encode(deb.Control)
	//b,_ := json.MarshalIndent(deb.Control, "", "  ")
	t.Logf("\n%s", bf.String())
	t.Logf("\n%s", deb.Control.RawText)
}

func Test_DebControl(t *testing.T) {
	debfile := "../testdata/com.deepin.appstore.helloworld_amd64_5.6.8-1.deb"
	in, err := os.Open(debfile)
	if err != nil {
		t.Fatalf("failed to open debfile->%+v", err)
	}
	deb, err := godeb.LoadDeb(in)
	if err != nil {
		t.Fatalf("failed to LoadDeb->%+v", err)
	}
	t.Logf("control_ext=%s", deb.ControlExt)
	t.Logf("data_ext=%s", deb.DataExt)

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "  ")
	jsonEncoder.Encode(deb.Control)
	//b,_ := json.MarshalIndent(deb.Control, "", "  ")
	t.Logf("\n%s", bf.String())
	t.Logf("\n%s", deb.Control.RawText)
}

func Test_DebData(t *testing.T) {
	debfile := "../testdata/com.deepin.appstore.helloworld_amd64_5.6.8-1.deb"
	in, err := os.Open(debfile)
	if err != nil {
		t.Fatalf("failed to open debfile->%+v", err)
	}
	deb, err := godeb.LoadDeb(in)
	if err != nil {
		t.Fatalf("failed to LoadDeb->%+v", err)
	}
	t.Logf("control_ext=%s", deb.ControlExt)
	t.Logf("data_ext=%s", deb.DataExt)

	r, closer, err := deb.Data.TarFile()
	if err != nil {
		t.Fatalf("failed to init tarfile->%+v", err)
	}
	defer closer.Close()

	for {
		var buf bytes.Buffer
		member, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed to read member->%+v", err)
		}
		if member.Typeflag == tar.TypeDir {
			continue
		}
		_, err = io.Copy(&buf, r)
		if err != nil {
			t.Fatalf("failed to copy->%+v", err)
		}
		t.Logf("%s\n", member.Name)
	}
}

func Test_Ar(t *testing.T) {
	debfile := "../testdata/com.deepin.appstore.helloworld_amd64_5.6.8-1.deb"
	in, err := os.Open(debfile)
	if err != nil {
		t.Fatalf("failed to open debfile->%+v", err)
	}
	contents, err := godeb.LoadAr(in)
	if err != nil {
		t.Fatalf("failed to LoadAr->%+v", err)
	}
	for name, aritem := range contents {
		t.Logf("ariten.name->%s", name)
		switch {
		case strings.Contains(name, "control."):
			t.Logf("===> %s", name)
			r, closer, err := aritem.TarFile()
			if err != nil {
				t.Fatalf("failed to init tarfile->%+v", err)
			}

			for {
				var buf bytes.Buffer
				member, err := r.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("failed to read member->%+v", err)
				}
				if member.Typeflag == tar.TypeDir {
					continue
				}
				_, err = io.Copy(&buf, r)
				if err != nil {
					t.Fatalf("failed to copy->%+v", err)
				}
				t.Logf("%s\n%s\n", member.Name, buf.String())
			}
			_ = closer.Close()
		}
	}
}
