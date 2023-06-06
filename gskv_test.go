package gskv

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
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

	db.Close()
	os.Remove("./def.gs")
}

func TestCloseOpen(t *testing.T) {
	db := Open()
	for i := 0; i < 1000; i++ {
		db.Set([]byte(fmt.Sprintf("k_%d", i)), []byte(fmt.Sprintf("v_%d", i)))
	}
	for i := 0; i < 100; i++ {
		db.Delete([]byte(fmt.Sprintf("k_%d", i)))
	}

	db.Close()

	db = Open()
	for i := 0; i < 1000; i++ {
		k := fmt.Sprintf("k_%d", i)
		v := db.Get([]byte(k))
		if i < 100 {
			if len(v) != 0 {
				t.Errorf("get deleted k(%s)  v(%s)", k, string(v))
			}
		} else {
			if !bytes.Equal([]byte(fmt.Sprintf("v_%d", i)), v) {
				t.Errorf("not equal k(%s)  v(%s)", k, string(v))
			}
		}
	}
	os.Remove("./def.gs")
}

func BenchmarkGet100000(b *testing.B) {
	db := Open()
	Get100000(b, db)
	os.Remove("./def.gs")
}

func BenchmarkRadnomGet100000(b *testing.B) {
	db := Open()
	RandomGet100000(b, db)
	os.Remove("./def.gs")
}

func RandomGet100000(b *testing.B, db *gskv) {
	b.StopTimer()
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

func Get100000(b *testing.B, db *gskv) {
	b.StopTimer()
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
