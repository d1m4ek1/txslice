package txslice

import (
	"context"
	"sync"
)

type TxIndexInterface[T any] interface {
	Add(v *T)
	Remove(v *T)
	Get(key any) (*T, bool)
}

type IndexKeyFunc[T any, K comparable] func(v *T) K

type indexOperation[T any, K comparable] struct {
	typ   string
	value *T
}

type Index[T any, K comparable] struct {
	mu   sync.RWMutex
	data map[K]*T

	keyFn IndexKeyFunc[T, K]
	ch    chan indexOperation[T, K]
}

func NewIndex[T any, K comparable](ctx context.Context, txsliceData *TxSlice[T], keyFn IndexKeyFunc[T, K]) {
	idx := &Index[T, K]{
		data:  make(map[K]*T, len(txsliceData.data)),
		keyFn: keyFn,
		ch:    make(chan indexOperation[T, K], 1024),
	}

	// отдельная горутина для обслуживания индекса
	go idx.loop(ctx)

	for _, v := range txsliceData.data {
		idx.data[keyFn(v)] = v
	}

	txsliceData.indexing = idx
}

func (i *Index[T, K]) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case op := <-i.ch:
			switch op.typ {
			case "add":
				i.add(i.keyFn(op.value), op.value)

			case "remove":
				i.delete(i.keyFn(op.value))
			}
		}
	}
}

func (i *Index[T, K]) add(key K, value *T) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.data[key] = value
}

func (i *Index[T, K]) delete(key K) {
	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.data, key)
}

func (i *Index[T, K]) Add(v *T) {
	i.ch <- indexOperation[T, K]{typ: "add", value: v}
}

func (i *Index[T, K]) Remove(v *T) {
	i.ch <- indexOperation[T, K]{typ: "remove", value: v}
}

func (i *Index[T, K]) Get(key any) (*T, bool) {
	k, ok := key.(K)
	if !ok {
		return nil, false
	}

	i.mu.RLock()
	defer i.mu.RUnlock()

	v, ok := i.data[k]
	return v, ok
}
