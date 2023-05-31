package skl

import (
	"bytes"
	"math"
	"math/rand"
)

const (
	maxHeight = 32
	pValue    = 1 / math.E
)

var probabilities [maxHeight]uint32

func init() {
	p := float64(1.0)
	for i := 0; i < maxHeight; i++ {
		probabilities[i] = uint32(float64(math.MaxUint32) * p)
		p *= pValue
	}
}

func NewSkipList() *SkipList {
	sl := &SkipList{
		lev:    1,
		length: 0,
		head: &skNode{
			levs: make([]*struct {
				forward *skNode
				span    uint32
			}, maxHeight),
		},
	}
	for i := 0; i < int(maxHeight); i++ {
		sl.head.levs[i] = &struct {
			forward *skNode
			span    uint32
		}{}
	}

	return sl
}

type SkipList struct {
	lev    uint32
	length uint32
	head   *skNode
}

type skNode struct {
	levs []*struct {
		forward *skNode
		span    uint32
	}
	k        []byte
	v        []byte
	backward *skNode
}

// 添加
func (sl *SkipList) Add(k, v []byte) {
	update := make([]*skNode, maxHeight)
	rank := make([]uint32, maxHeight)

	// 检测要更新的节点
	n := sl.head
	for lev := int64(sl.lev - 1); lev >= 0; lev-- {
		if lev != int64(sl.lev-1) {
			rank[lev] = rank[lev+1]
		}
		for n.levs[lev].forward != nil {
			c := bytes.Compare(k, n.levs[lev].forward.k)
			if c == 0 {
				// 存在相同的key，则直接更新后返回
				n.levs[lev].forward.v = v
				return
			}
			if c != 1 {
				break
			}
			n = n.levs[lev].forward
			rank[lev] += n.levs[lev].span
		}
		update[lev] = n
	}

	// 随机高度
	addLev := sl.randomHeight()

	// 初始化
	for i := sl.lev; i < addLev; i++ {
		update[i] = sl.head
		update[i].levs[i].span = sl.length
	}

	add := &skNode{
		levs: make([]*struct {
			forward *skNode
			span    uint32
		}, addLev),
		k: k,
		v: v,
	}
	// 写入
	for lev := int64(addLev - 1); lev >= 0; lev-- {
		add.levs[lev] = &struct {
			forward *skNode
			span    uint32
		}{}
		update[lev].levs[lev].forward, add.levs[lev].forward = add, update[lev].levs[lev].forward

		add.levs[lev].span = update[lev].levs[lev].span - (rank[0] - rank[lev])
		update[lev].levs[lev].span = rank[0] - rank[lev] + 1
	}

	for i := sl.lev; i < addLev; i++ {
		update[i].levs[i].span += 1
	}
	if update[0] != sl.head {
		add.backward = update[0]
	}
	if add.levs[0].forward != nil {
		add.levs[0].forward.backward = add
	}

	if sl.lev < addLev {
		sl.lev = addLev
	}
	sl.length++
}

// 删除
func (sl *SkipList) Delete(k []byte) {
	update := make([]*skNode, maxHeight)
	rank := make([]int64, maxHeight)

	// 检测要更新的节点
	n := sl.head
	for lev := int64(sl.lev - 1); lev >= 0; lev-- {
		if lev != int64(sl.lev-1) {
			rank[lev] = rank[lev+1]
		}
		for n.levs[lev].forward != nil {
			if bytes.Compare(k, n.levs[lev].forward.k) != 1 {
				break
			}
			n = n.levs[lev].forward
		}
		update[lev] = n
	}
	if !bytes.Equal(update[0].levs[0].forward.k, k) {
		return
	}
	// 这里说明找到了数据
	dn := update[0].levs[0].forward

	h := len(dn.levs)
	for lev := h - 1; lev >= 0; lev-- {
		update[lev].levs[lev].span += dn.levs[lev].span - 1
		if dn.levs[lev].forward == nil {
			update[lev].levs[lev].forward = nil
		} else {
			update[lev].levs[lev].forward, dn.levs[lev].forward.backward = dn.levs[lev].forward, update[lev].levs[lev].forward
		}
	}
}

// 获取
func (sl *SkipList) Get(k []byte) ([]byte, bool) {
	// 检测要更新的节点
	n := sl.head
	for lev := int64(sl.lev - 1); lev >= 0; lev-- {
		for n.levs[lev] != nil && n.levs[lev].forward != nil {
			c := bytes.Compare(k, n.levs[lev].forward.k)
			if c == 0 {
				return n.levs[lev].forward.v, true
			}
			if c == 1 {
				n = n.levs[lev].forward
			} else {
				break
			}
		}
	}
	return nil, false
}

func (sl *SkipList) randomHeight() uint32 {
	rnd := rand.Uint32()

	h := uint32(1)
	for h < uint32(maxHeight) && rnd <= probabilities[h] {
		h++
	}

	return h
}
