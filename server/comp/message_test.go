package comp_test

import (
	"bytes"
	"github.com/appcrash/media/server/comp"
	"testing"
)

func TestMessageHeader(t *testing.T) {
	m := &comp.RawByteMessage{}
	m.SetHeader("key1", []byte("value1"))
	m.SetHeader("key2", []byte("value2"))
	if bytes.Compare(m.GetHeader("key1"), []byte("value1")) != 0 {
		t.Fatal("set key wrong")
	}
	if bytes.Compare(m.GetHeader("key2"), []byte("value2")) != 0 {
		t.Fatal("set key wrong")
	}

	m.Meta = []byte("abcd=abc;dabc=abc;zyxabcabc=abc/de;abc=/correct Value/;")
	value := m.GetHeader("abc")
	if bytes.Compare(value, []byte("/correct Value/")) != 0 {
		t.Fatalf("get key wrong: %v", string(value))
	}

	m.Meta = []byte("abc=;;;;cabc=x;")
	value = m.GetHeader("abc")
	if value != nil {
		t.Fatalf("should not get the key: %v %v", string(value), len(value))
	}

	m.Meta = []byte("abc=;;abc=;abc=x;")
	value = m.GetHeader("abc")
	if len(value) != 1 || value[0] != 'x' {
		t.Fatalf("get key wrong: %v", string(value))
	}
}
