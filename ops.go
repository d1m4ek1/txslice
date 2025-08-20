package txslice

func (t *TxSlice[T]) Push(n ...*T) {
	t.data = append(t.data, n...)

	t.pushJournal(&operation[T]{
		typ:           opPush,
		countAppended: len(n),
		values:        n,
	})
}

func (t *TxSlice[T]) Pop() *T {
	lastIndex := t.Len() - 1
	item := t.data[lastIndex]

	t.data = t.data[:lastIndex]

	t.pushJournal(&operation[T]{
		typ:     opPop,
		indexes: []int{lastIndex},
		values:  []*T{item},
	})

	return item
}

func (t *TxSlice[T]) Shift() *T {
	item := t.data[0]

	t.data = t.data[1:]

	t.pushJournal(&operation[T]{
		typ:    opShift,
		values: []*T{item},
	})

	return item
}

func (t *TxSlice[T]) Insert(index int, item *T) {
	if index < 0 || t.Len() < index {
		return
	}

	t.data = append(t.data, item)
	copy(t.data[index+1:], t.data[index:])

	t.pushJournal(&operation[T]{
		typ:     opInsert,
		indexes: []int{index},
	})
}

func (t *TxSlice[T]) Set(index int, item *T) {
	if index < 0 || index > t.Len() {
		return
	}

	t.data[index] = item

	t.pushJournal(&operation[T]{
		typ:     opSet,
		indexes: []int{index},
		values:  []*T{item},
	})
}
