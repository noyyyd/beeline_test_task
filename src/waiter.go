package main

import (
	"context"
)

type waiter struct {
	isWait       bool
	waitersCount int
	waitChan     chan struct{}
}

func newWaiter() *waiter {
	w := new(waiter)
	w.waitChan = make(chan struct{})

	return w
}

func (w *waiter) wait(ctx context.Context) {
	if !w.isWait {
		w.isWait = true
	}

	w.waitersCount++

	go func() {
		select {
		case <-ctx.Done():
			if w.isWait {
				w.stopWait()
			}
			return
		}
	}()
}

func (w *waiter) stopWait() {
	w.waitChan <- struct{}{}

	w.waitersCount--

	if w.waitersCount == 0 {
		w.isWait = false
	}
}
