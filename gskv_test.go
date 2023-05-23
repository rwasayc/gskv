package gskv

import (
	"bytes"
	"fmt"
	"math/rand"
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
			CMD:  "SET",
			ArgK: []byte("k1"),
			ArgV: []byte("v1"),
		},
		{
			CMD:    "GET",
			ArgK:   []byte("k1"),
			Expect: []byte("v1"),
		},
		{
			CMD:  "DELETE",
			ArgK: []byte("k1"),
		},
		{
			CMD:    "GET",
			ArgK:   []byte("k1"),
			Expect: nil,
		},
		{
			CMD:  "SET",
			ArgK: []byte("k2"),
			ArgV: []byte("v2"),
		},
		{
			CMD:  "SET",
			ArgK: []byte("k3"),
			ArgV: []byte("v3"),
		},

		{
			CMD:  "DELETE",
			ArgK: []byte("k4"),
		},
		{
			CMD:    "GET",
			ArgK:   []byte("k3"),
			Expect: []byte("v3"),
		},
		{
			CMD:  "DELETE",
			ArgK: []byte("k2"),
		},
		{
			CMD:  "DELETE",
			ArgK: []byte("k3"),
		},
		{
			CMD:    "GET",
			ArgK:   []byte("k3"),
			Expect: nil,
		},
	}

	db := Open()

	for _, c := range caseList {
		switch c.CMD {
		case "GET":
			r := db.Get(c.ArgK)
			if !bytes.Equal(r, c.Expect) {
				t.Errorf("get c.ArgK(%s) ->%v  expect %v", string(c.ArgK), string(r), string(c.Expect))
				t.FailNow()
			}
		case "SET":
			db.Set(c.ArgK, c.ArgV)
		case "DELETE":
			db.Delete(c.ArgK)
		}
	}
}

func BenchmarkGet100000(b *testing.B) {
	b.StopTimer()
	db := Open()
	limit := 100000
	for i := 0; i <= limit; i++ {
		db.Set([]byte(fmt.Sprintf("k%d", i)), []byte(fmt.Sprintf("v%d", i)))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		idx := i % limit
		db.Get([]byte(fmt.Sprintf("k%d", idx)))
	}
}

func BenchmarkRandomGet100000(b *testing.B) {
	b.StopTimer()
	db := Open()
	limit := 100000
	for i := 0; i <= limit*10; i++ {
		v := rand.Int31n(int32(limit))
		db.Set([]byte(fmt.Sprintf("k%d", v)), []byte(fmt.Sprintf("v%d", v)))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		idx := i % limit
		db.Get([]byte(fmt.Sprintf("k%d", idx)))
	}
}
