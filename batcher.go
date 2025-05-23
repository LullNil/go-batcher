package batcher

import (
	"context"
	"sync"
	"time"
)

type DoFunc[T any] func(batch []T)

type Batcher[T any] struct {
	buffer    int
	timeout   time.Duration
	do        DoFunc[T]
	inputChan chan T

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New creates a new Batcher with context, timeout (can be 0), and batch size.
func New[T any](buffer int, timeout time.Duration, do DoFunc[T]) *Batcher[T] {
	if buffer <= 0 {
		panic("buffer must be > 0")
	}
	ctx, cancel := context.WithCancel(context.Background())

	b := &Batcher[T]{
		buffer:    buffer,
		timeout:   timeout,
		do:        do,
		inputChan: make(chan T),
		ctx:       ctx,
		cancel:    cancel,
	}

	b.wg.Add(1)
	go b.run()

	return b
}

// Add sends an item to the batcher (non-blocking).
func (b *Batcher[T]) Add(values ...T) {
	for _, v := range values {
		select {
		case <-b.ctx.Done():
			return
		case b.inputChan <- v:
		}
	}
}

// Close flushes remaining items and stops the batcher gracefully.
func (b *Batcher[T]) Close() {
	b.cancel()
	b.wg.Wait()
}

func (b *Batcher[T]) run() {
	defer b.wg.Done()

	var (
		batch  []T
		timer  *time.Timer
		timerC <-chan time.Time // nil by default
	)

	for {
		select {
		case <-b.ctx.Done():
			if len(batch) > 0 {
				b.do(batch)
			}
			if timer != nil {
				timer.Stop()
			}
			return

		case item := <-b.inputChan:
			batch = append(batch, item)

			if len(batch) == 1 && b.timeout > 0 {
				// The first element is to start the timer
				timer = time.NewTimer(b.timeout)
				timerC = timer.C
			}

			if len(batch) >= b.buffer {
				b.do(batch)
				batch = nil

				if timer != nil {
					timer.Stop()
					timer = nil
					timerC = nil
				}
			}

		case <-timerC:
			if len(batch) > 0 {
				b.do(batch)
				batch = nil
			}
			// Reset timer
			if timer != nil {
				timer.Stop()
				timer = nil
				timerC = nil
			}
		}
	}
}
