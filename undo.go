package txslice

func (t *TxSlice[T]) undoPush(op *operation[T]) {
	t.data = t.data[:t.Len()-op.countAppended]
}

func (t *TxSlice[T]) undoPop(op *operation[T]) {
	t.data = append(t.data, op.values...)
}

func (t *TxSlice[T]) undoShift(op *operation[T]) {
	t.data = append(op.values, t.data...)
}

func (t *TxSlice[T]) undoInsert(op *operation[T]) {
	t.data = append(t.data[:op.indexes[0]], t.data[op.indexes[0]+1:]...)
}

func (t *TxSlice[T]) undoSet(op *operation[T]) {
	t.data[op.indexes[0]] = op.values[0]
}

func (t *TxSlice[T]) undoModSwap(op *operation[T]) {
	t.data[op.indexes[0]], t.data[op.indexes[1]] = t.data[op.indexes[1]], t.data[op.indexes[0]]
}

func (t *TxSlice[T]) undoModMove(op *operation[T]) {
	val := t.data[op.indexes[0]]

	t.data = append(t.data[:op.indexes[0]], t.data[op.indexes[0]+1:]...)

	if op.indexes[1] >= t.Len() {
		t.data = append(t.data, val)
	} else {
		t.data = append(t.data[:op.indexes[1]], append([]*T{val}, t.data[op.indexes[1]:]...)...)
	}
}
