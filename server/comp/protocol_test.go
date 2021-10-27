package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"testing"
)

func TestWithString(t *testing.T) {
	s1 := "cmd1 arg_a argb"
	if r, err := comp.WithString(s1); err != nil {
		t.Fatal(err)
	} else {
		if r[0] != "cmd1" || r[1] != "arg_a" || r[2] != "argb" {
			t.Errorf("parse (%v) error", s1)
		}
	}

	s2 := " cmd2 \"arga\" \"arg b\" \"\" "
	if r, err := comp.WithString(s2); err != nil {
		t.Fatal(err)
	} else {
		if r[0] != "cmd2" || r[1] != "arga" || r[2] != "arg b" || r[3] != "" || len(r) != 4 {
			t.Errorf("parse (%v) error", s2)
		}
	}

	s3 := "cmd3 print  始めまして "
	if r, err := comp.WithString(s3); err != nil {
		t.Fatal(err)
	} else {
		if r[0] != "cmd3" || r[1] != "print" || r[2] != "始めまして" || len(r) != 3 {
			t.Errorf("parse (%v) error", s3)
		}
	}

	s4 := ` \cmd4  arg\s  "contain\"quote"  `
	if r, err := comp.WithString(s4); err != nil {
		t.Fatal(err)
	} else {
		if r[0] != "\\cmd4" || r[1] != "arg\\s" || r[2] != "contain\"quote" || len(r) != 3 {
			t.Errorf("parse (%v) error", s4)
		}
	}
}
