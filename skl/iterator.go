package skl

func NewIterator(skl *SkipList) *iterator {
	return &iterator{
		skl: skl,
	}
}

type iterator struct {
	c   *skNode
	skl *SkipList
}

func (it *iterator) First() ([]byte, []byte) {
	it.c = it.skl.head
	return it.Next()
}

func (it *iterator) Next() ([]byte, []byte) {
	if it.c == nil || it.c.levs[0] == nil {
		return nil, nil
	}
	it.c = it.c.levs[0].forward
	if it.c == nil {
		return nil, nil
	}
	return it.c.k, it.c.v
}
