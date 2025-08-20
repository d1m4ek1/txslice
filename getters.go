package txslice

func (t *TxSlice[T]) Slice() []*T {
	return t.data
}

func (t *TxSlice[T]) Len() int {
	return len(t.data)
}

func (t *TxSlice[T]) FirstElement() *T {
	return t.data[0]
}

func (t *TxSlice[T]) LastElement() *T {
	return t.data[t.Len()-1]
}

func (t *TxSlice[T]) MiddleElement() *T {
	return t.data[(t.Len()-1)/2]
}

func (t *TxSlice[T]) Find(predicate func(*T) bool) (int, *T, bool) {
	for i, v := range t.data {
		if predicate(v) {
			return i, v, true
		}
	}

	return -1, nil, false
}

func (t *TxSlice[T]) IndexFind(key any) (*T, bool) {
	return t.indexing.Get(key)
}
