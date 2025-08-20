# txslice

`txslice` is a transactional slice library for Go, designed for **high-performance workloads** where millions of iterations and minimal memory allocations matter.  
It provides a slice wrapper with transaction support, rollback, snapshots, and utility getters/search methods.

## Features

-   🔄 **Transactional operations**: Push, Pop, Shift, Insert, Set, Swap, Move, with full **rollback** support.
-   📸 **Snapshots**: Capture the current state of the slice (latest or versioned).
    -   Automatic snapshots support (`isAutoLatestSnap`)
    -   Versioned snapshots for rollback or historical access
-   🔒 **Thread-safe**: All operations are protected by `sync.RWMutex` for safe concurrent usage.
-   🧩 **Indexing**: Optional key-based indexing for fast access.
    -   Supports `Add`, `Remove`, and async handling for large batches
    -   `IndexGet` for O(1)-like lookups
-   ⚡ **Efficient search**:
    -   `Find`: linear search with custom predicate
    -   `BinaryFind`: binary search on sorted slices using a key function
-   📏 **Utility getters**: `FirstElement`, `LastElement`, `MiddleElement` and more.
-   🏎 **High-performance**:
    -   Minimal allocations where possible
    -   Batch operations for large datasets
-   🔧 **Debug mode**: Logs detailed rollback/commit operations for debugging

## Installation

```bash
go get github.com/d1m4ek1/txslice
```

## Example

### Simple Transactional Slice Operations

TxSlice Initialization and Simple Transactional Slice Operations

```golang
package main

import (
    "fmt"
    "github.com/d1m4ek1/txslice"
)

type Values struct {
    Val int
}

func main() {
    tx := txslice.New([]*Values{}, txslice.Config{})

	// Push elements
	tx.Push(&Values{1})
	tx.Push(&Values{2})

	fmt.Println("First element:", *tx.FirstElement()) // {1}

	// Create a snapshot
	tx.SetSnapshot("v1")

	// Commit example
	tx.Commit()

	// Transaction with rollback
	tx.Push(&Values{3})
	tx.Rollback() // rollback last transaction

	fmt.Println("Length after rollback:", tx.Len()) // 2
}
```

### Batch Operations

```go
package main

import (
	"fmt"

	"github.com/d1m4ek1/txslice"
)

type Values struct {
    Val int
}

func main() {
		// Initialize TxSlice
	tx := txslice.New([]*Values{{1}, {2}, {3}, {4}, {5}}, txslice.Config{})

	// Perform a batch of operations atomically
	if err := tx.Batch(func(b *txslice.TxSlice[Values]) error {
		b.Push(&Values{6}, &Values{7})

		b.ModSwap(0, 2)

		if b.FirstElement().Val != 3 {
			return errors.New("something error") // this error will be returned from Batch()
		}

		return nil
	}); err != nil {
		panic(err)
	}

	res := ""

	for _, item := range tx.Slice() {
		res += fmt.Sprintf("%d ", item.Val)
	}

	fmt.Println(res) // output: 3 2 1 4 5 6 7

	b := tx.BatchStart()

	b.ModMove(0, 3)

	b.BatchAccept()

	res = ""

	for _, item := range tx.Slice() {
		res += fmt.Sprintf("%d ", item.Val)
	}

	fmt.Println(res) // output: 2 1 4 3 5 6 7

	b = tx.BatchStart()

	b.ModMove(0, 3)

	b.UndoBatch()

	res = ""

	for _, item := range tx.Slice() {
		res += fmt.Sprintf("%d ", item.Val)
	}

	fmt.Println(res) // output: 2 1 4 3 5 6 7
}
```

TxSlice is created for developers who need performance, safety, and flexibility when working with slices in
