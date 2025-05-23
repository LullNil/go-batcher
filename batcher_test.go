package batcher

import (
	"sync"
	"testing"
	"time"
)

func TestBatcher_BatchSizeTrigger(t *testing.T) {
	var mu sync.Mutex
	var batches [][]int

	b := New(3, 0, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		batches = append(batches, batch)
	})

	b.Add(1)
	b.Add(2)
	b.Add(3) // => trigger

	time.Sleep(100 * time.Millisecond)
	b.Close()

	mu.Lock()
	defer mu.Unlock()
	if len(batches) != 1 || len(batches[0]) != 3 {
		t.Errorf("Expected 1 batch of 3, got %+v", batches)
	}
}

func TestBatcher_TimeoutTrigger(t *testing.T) {
	var mu sync.Mutex
	var batches [][]int

	b := New(10, 100*time.Millisecond, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		batches = append(batches, batch)
	})

	b.Add(42)

	time.Sleep(200 * time.Millisecond)
	b.Close()

	mu.Lock()
	defer mu.Unlock()
	if len(batches) != 1 || len(batches[0]) != 1 {
		t.Errorf("Expected timeout-triggered batch, got %+v", batches)
	}
}

func TestBatcher_MultipleFlushes(t *testing.T) {
	var mu sync.Mutex
	var batches [][]int

	b := New(2, 0, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		batches = append(batches, batch)
	})

	b.Add(1)
	b.Add(2) // flush
	b.Add(3)
	b.Add(4) // flush

	time.Sleep(50 * time.Millisecond)
	b.Close()

	mu.Lock()
	defer mu.Unlock()
	if len(batches) != 2 {
		t.Errorf("Expected 2 batches, got %+v", batches)
	}
}

func TestBatcher_MultipleItem(t *testing.T) {
	var (
		mu     sync.Mutex
		result []int
	)

	b := New(10, 50*time.Millisecond, func(items []int) {
		mu.Lock()
		defer mu.Unlock()
		result = append(result, items...)
	})

	b.Add(1, 2, 3)
	b.Add(4, 5)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(result) != 5 {
		t.Errorf("Expected 5 items, got %d", len(result))
	}
}

func BenchmarkBatcher_Add(b *testing.B) {
	batcher := New(1000, 0, func(batch []int) {})

	for i := 0; i < b.N; i++ {
		batcher.Add(i)
	}

	batcher.Close()
}
