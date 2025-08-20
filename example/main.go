package main

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"github.com/d1m4ek1/txslice"
)

func main() {
	t := NewSomeSlice(1_00)

	sizeBytes := len(t) * int(unsafe.Sizeof(t[0]))

	sizeMB := float64(sizeBytes) / 1024 / 1024
	fmt.Printf("Size: %.6f MB\n", sizeMB)

	timeStart := time.Now()

	tx := txslice.New(t, txslice.Config{
		IsAutoLatestSnap: true,
	})

	txslice.NewIndex(context.Background(), tx, func(v *some) string { return v.ID })

	fmt.Println(time.Since(timeStart), "=====> 1")

	timeStart = time.Now()

	tx.IndexFind(t[50].ID)

	fmt.Println(time.Since(timeStart), "=====> 2")

	timeStart = time.Now()

	tx.Find(func(s *some) bool { return s.ID == t[50].ID })

	fmt.Println(time.Since(timeStart), "=====> 3")

	fmt.Println(tx.Len() == len(t), tx.Len(), len(t))
}
