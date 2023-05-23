package gskv

import "bytes"

func Open() *gskv {
	return &gskv{
		logs: make([]*kvLog, 0, 1024*512),
	}
}

type kvLog struct {
	k []byte
	v []byte
	d bool
}

type gskv struct {
	logs []*kvLog
}

func (gs *gskv) Set(k, v []byte) {
	gs.logs = append(gs.logs, &kvLog{
		k: k,
		v: v,
	})
}

func (gs *gskv) Get(k []byte) []byte {
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
}
