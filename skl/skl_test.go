package skl

import (
	"bytes"
	"testing"
)

func TestAll(t *testing.T) {
	caseList := []struct {
		CMD    string
		ArgK   []byte
		ArgV   []byte
		Expect []byte
	}{
		{
			CMD:    "GET",
			ArgK:   []byte("k1"),
			Expect: nil,
		},
		{
			CMD:  "ADD",
			ArgK: []byte("k1"),
			ArgV: []byte("v1"),
		},
		{
			CMD:    "GET",
			ArgK:   []byte("k1"),
			Expect: []byte("v1"),
		},
		{
			CMD:  "ADD",
			ArgK: []byte("k2"),
			ArgV: []byte("v2"),
		},
		{
			CMD:  "ADD",
			ArgK: []byte("k3"),
			ArgV: []byte("v3"),
		},
		{
			CMD:    "GET",
			ArgK:   []byte("k2"),
			Expect: []byte("v2"),
		},
		{
			CMD:  "DEL",
			ArgK: []byte("k2"),
		},
		{
			CMD:    "GET",
			ArgK:   []byte("k2"),
			Expect: nil,
		},
		{
			CMD:    "GET",
			ArgK:   []byte("k3"),
			Expect: []byte("v3"),
		},
	}
	sl := NewSkipList()
	for idx, c := range caseList {
		t.Logf("[%v] run CMD %v\n", idx+1, c.CMD)
		switch c.CMD {
		case "GET":
			r, _ := sl.Get(c.ArgK)
			if !bytes.Equal(r, c.Expect) {
				t.Errorf("get sl.ArgK(%s) ->'%v'  expect %v", string(c.ArgK), string(r), string(c.Expect))
				t.FailNow()
			}
		case "ADD":
			sl.Add(c.ArgK, c.ArgV)
		case "DEL":
			sl.Delete(c.ArgK)
		}
	}
}
