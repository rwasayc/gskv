package skl

import (
	"bytes"
	"fmt"
	"testing"
)

func TestIterator(t *testing.T) {
	caseList := []struct {
		CMD       string
		ArgBatchN int
		ArgK      []byte
		ArgV      []byte
		ExpectK   []byte
		ExpectV   []byte
	}{
		{
			CMD: "NEXT",
		},
		{
			CMD: "FIRST",
		},
		{
			CMD:       "BATCH_RAND_ADD",
			ArgBatchN: 3,
		},
		{
			CMD:     "FIRST",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 1)),
			ExpectV: []byte("v_1"),
		},
		{
			CMD:     "NEXT",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 2)),
			ExpectV: []byte("v_2"),
		},
		{
			CMD:     "NEXT",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 3)),
			ExpectV: []byte("v_3"),
		},
		{
			CMD: "NEXT",
		},
		{
			CMD:       "BATCH_RAND_ADD",
			ArgBatchN: 100,
		},
		{
			CMD:     "FIRST",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 1)),
			ExpectV: []byte("v_1"),
		},
		{
			CMD:     "NEXT",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 2)),
			ExpectV: []byte("v_2"),
		},
		{
			CMD:     "NEXT",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 3)),
			ExpectV: []byte("v_3"),
		},
		{
			CMD:     "FIRST",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 1)),
			ExpectV: []byte("v_1"),
		},
		{
			CMD:     "NEXT",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 2)),
			ExpectV: []byte("v_2"),
		},
		{
			CMD:  "DEL",
			ArgK: []byte(fmt.Sprintf("k_%05d", 1)),
		},
		{
			CMD:     "FIRST",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 2)),
			ExpectV: []byte("v_2"),
		},
		{
			CMD:     "NEXT",
			ExpectK: []byte(fmt.Sprintf("k_%05d", 3)),
			ExpectV: []byte("v_3"),
		},
	}
	skl := NewSkipList()
	iter := NewIterator(skl)

	for idx, c := range caseList {
		switch c.CMD {
		case "NEXT":
			k, v := iter.Next()
			if !bytes.Equal(c.ExpectK, k) {
				t.Errorf("[idx:%v]!bytes.Equal(c.ExpectK, k) c.ExpectK:%v  k:%v", idx, string(c.ExpectK), string(k))
				t.FailNow()
			}
			if !bytes.Equal(c.ExpectV, v) {
				t.Errorf("[idx:%v]!bytes.Equal(c.ExpectV, v) c.ExpectV:%v  v:%v", idx, string(c.ExpectV), string(v))
				t.FailNow()
			}
		case "FIRST":
			k, v := iter.First()
			if !bytes.Equal(c.ExpectK, k) {
				t.Errorf("[idx:%v]!bytes.Equal(c.ExpectK, k) c.ExpectK:%v  k:%v", idx, string(c.ExpectK), string(k))
				t.FailNow()
			}
			if !bytes.Equal(c.ExpectV, v) {
				t.Errorf("[idx:%v]!bytes.Equal(c.ExpectV, v) c.ExpectV:%v  v:%v", idx, string(c.ExpectV), string(v))
				t.FailNow()
			}
		case "ADD":
			skl.Add(c.ArgK, c.ArgV)
		case "DEL":
			skl.Delete(c.ArgK)
		case "BATCH_RAND_ADD":
			for i := 1; i <= c.ArgBatchN; i++ {
				skl.Add([]byte(fmt.Sprintf("k_%05d", i)), []byte(fmt.Sprintf("v_%d", i)))
			}
		}
	}

}
