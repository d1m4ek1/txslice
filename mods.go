package txslice

func (t *TxSlice[T]) ModSwap(index1, index2 int) {
	if index1 < 0 || index2 < 0 || index1 >= t.Len() || index2 >= t.Len() {
		return
	}

	t.data[index1], t.data[index2] = t.data[index2], t.data[index1]

	t.journal = append(t.journal, &operation[T]{
		typ:     modSwap,
		indexes: []int{index2, index1},
	})
}

func (t *TxSlice[T]) ModMove(from, to int) {
	if from < 0 || from >= t.Len() || to < 0 || to >= t.Len() {
		return
	}

	if from == to {
		return
	}

	val := t.data[from]

	t.data = append(t.data[:from], t.data[from+1:]...)

	if to >= len(t.data) {
		t.data = append(t.data, val)
	} else {
		t.data = append(t.data[:to], append([]*T{val}, t.data[to:]...)...)
	}

	t.pushJournal(&operation[T]{
		typ:     modMove,
		indexes: []int{to, from},
	})
}
