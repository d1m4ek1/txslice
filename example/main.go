package main

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"txslice"
)

func main() {
	t := NewSomeSlice(1_0)

	sizeBytes := len(t) * int(unsafe.Sizeof(t[0]))

	sizeMB := float64(sizeBytes) / 1024 / 1024
	fmt.Printf("Size: %.6f MB\n", sizeMB)

	timeStart := time.Now()

	tx := txslice.New(t, txslice.Config{
		IsAutoLatestSnap: true,
	})

	if err := tx.Batch(func(b *txslice.TxSlice[Some]) error {
		g := NewSomeSlice(1_0)

		for _, item := range g {
			b.Push(item)
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	tx.Rollback()

	fmt.Println(time.Since(timeStart))

	fmt.Println(tx.Len() == len(t), tx.Len(), len(t))
}
