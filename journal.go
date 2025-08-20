package txslice

type (
	opType int
)

const (
	opPush opType = iota
	opPop
	opShift
	opUnshift
	opInsert
	opSplice
	opSet
	modSwap
	modMove
	batch
)

const (
	defaultJournalMinCap = 16
	defaultJournalStep   = 16
)

type operation[T any] struct {
	typ           opType // тип операции
	countAppended int
	indexes       []int // индексы
	values        []*T  // значения, которые надо восстановить при откате
	nested        []*operation[T]
}

func newJournal[T any](capacity int) []*operation[T] {
	return make([]*operation[T], 0, capacity)
}

func (t *TxSlice[T]) pushJournal(op *operation[T]) {
	// Проверка на переполнение capacity
	if len(t.journal) == cap(t.journal) {
		step := max(t.journalStep, cap(t.journal)/4)

		newCap := cap(t.journal) + step

		newJournal := make([]*operation[T], len(t.journal), newCap)

		copy(newJournal, t.journal)

		t.journal = newJournal
	}

	t.journal = append(t.journal, op)
}

func (t *TxSlice[T]) SetJournalStep(step int) {
	if step > 0 {
		t.journalStep = step
	}
}

func (t *TxSlice[T]) JournalInfo() (length int, capacity int) {
	return len(t.journal), cap(t.journal)
}
