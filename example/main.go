package main

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"github.com/d1m4ek1/txslice"
)

func main() {
	t := NewSomeSlice(1_000_000)

	sizeBytes := len(t) * int(unsafe.Sizeof(t[0]))

	sizeMB := float64(sizeBytes) / 1024 / 1024
	fmt.Printf("Size: %.6f MB\n", sizeMB)

	timeStart := time.Now()

	tx := txslice.New(t, txslice.Config{
		IsAutoLatestSnap: true,
		IsDebug:          true,
	})

	txslice.NewIndex(context.Background(), tx, func(v *some) string { return v.ID }, 2048)

	tx.Batch(func(b *txslice.TxSlice[some]) error {
		b.Push(NewSomeSlice(20)...)
		b.Pop()
		b.Shift()

		b.ModSwap(10, 60)

		b.ModMove(50, 99)

		return nil
	})

	tx.Rollback()

	fmt.Println(time.Since(timeStart), "=====> 1")

	timeStart = time.Now()

	fmt.Println(tx.IndexGet(t[500000].ID))

	fmt.Println(time.Since(timeStart), "=====> 2")

	timeStart = time.Now()

	fmt.Println(tx.Find(func(s *some) bool { return s.ID == t[500000].ID }))

	fmt.Println(time.Since(timeStart), "=====> 3")

	fmt.Println(tx.Len() == len(t), tx.Len(), len(t))
}
