package gskv

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/rwasayc/gskv/skl"
	"github.com/rwasayc/gskv/vfs"
)

func Open() *gskv {
	fs := vfs.NewOSFS()
	f, err := fs.OpenFile("./def.gs", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	sk := skl.NewSkipList()
	var offset int64
	header := make([]byte, 9)
	for {
		n, err := f.ReadAt(header, offset)
		offset += 9
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF || n == 0 {
			break
		}
		if n != len(header) {
			panic("n != len(header)")
		}

		kl := binary.LittleEndian.Uint32(header[:4])
		vl := binary.LittleEndian.Uint32(header[4:8])
		delete := header[8] == 1

		if delete {
			kv := make([]byte, kl+vl)
			kvn, err := f.ReadAt(kv, offset)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if err == io.EOF || uint32(kvn) != kl+vl {
				panic(fmt.Errorf("err == io.EOF || uint32(kvn) != kl+vl  kvn(%d) kl(%d) vl(%d)  err(%v)", kvn, kl, vl, err))
			}
			sk.Delete(kv[:kl])
		}
		if !delete {
			kv := make([]byte, kl+vl)
			kvn, err := f.ReadAt(kv, offset)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if err == io.EOF || uint32(kvn) != kl+vl {
				panic(fmt.Errorf("err == io.EOF || uint32(kvn) != kl+vl  kvn(%d) kl(%d) vl(%d)  err(%v)", kvn, kl, vl, err))
			}
			sk.Add(kv[:kl], kv[kl:kl+vl])
		}
		offset += int64((kl + vl))
	}

	return &gskv{
		sk:   sk,
		flog: f,
	}
}

type gskv struct {
	flog vfs.File
	sk   *skl.SkipList
}

func (gs *gskv) Set(k, v []byte) {
	data := make([]byte, 9, len(k)+len(v)+9)
	binary.LittleEndian.PutUint32(data, uint32(len(k)))
	binary.LittleEndian.PutUint32(data[4:], uint32(len(v)))
	data = append(data, k...)
	data = append(data, v...)
	_, err := gs.flog.Write(data)
	if err != nil {
		panic(err)
	}
	gs.sk.Add(k, v)
}

func (gs *gskv) Get(k []byte) []byte {
	v, exist := gs.sk.Get(k)
	if exist {
		return v
	}
	return nil
}

func (gs *gskv) Delete(k []byte) {
	_, exist := gs.sk.Get(k)
	if !exist {
		return
	}

	data := make([]byte, 9, 9+len(k))
	binary.LittleEndian.PutUint32(data, uint32(len(k)))
	data[8] = 1
	data = append(data, k...)
	_, err := gs.flog.Write(data)
	if err != nil {
		panic(err)
	}

	gs.sk.Delete(k)
}

func (gs *gskv) Close() {
	gs.flog.Close()
}
