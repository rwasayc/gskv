package gskv

import (
	"bytes"

	"github.com/rwasayc/gskv/skl"
)

func Open() *gskv {
	return &gskv{
		logs: make([]*kvLog, 0, 1024*512),
	}
}

func OpenSK() *gskv {
	return &gskv{
		logs: make([]*kvLog, 0, 1024*512),
		sk:   skl.NewSkipList(),
	}
}

type kvLog struct {
	k []byte
	v []byte
	d bool
}

type gskv struct {
	logs []*kvLog
	sk   *skl.SkipList
}

func (gs *gskv) Set(k, v []byte) {
	gs.logs = append(gs.logs, &kvLog{
		k: k,
		v: v,
	})
	if gs.sk != nil {
		gs.sk.Add(k, v)
	}
}

func (gs *gskv) Get(k []byte) []byte {
	if gs.sk != nil {
		v, exist := gs.sk.Get(k)
		if exist {
			return v
		}
	}

	for idx := len(gs.logs) - 1; idx >= 0; idx-- {
		if !bytes.Equal(gs.logs[idx].k, k) {
			continue
		}
		if gs.logs[idx].d {
			return nil
		}
		return gs.logs[idx].v
	}
	return nil
}

func (gs *gskv) Delete(k []byte) {
	gs.logs = append(gs.logs, &kvLog{
		k: k,
		d: true,
	})
	if gs.sk != nil {
		gs.sk.Delete(k)
	}
}
