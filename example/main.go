package main

import (
	"fmt"
	"time"

	"github.com/LullNil/go-batcher"
)

func main() {
	b := batcher.New(3, 500*time.Millisecond, func(batch []string) {
		fmt.Printf("Batch processed: %v\n", batch)
	})

	b.Add("alpha", "beta", "gamma")
	b.Add("delta")

	b.Close()
}
