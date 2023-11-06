package use_cases

import "context"

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
	w.isWait = true
	w.waitersCount++

	select {
	case <-ctx.Done():
	case <-w.waitChan:
	}

	w.changeWaitStatus()
}

func (w *waiter) changeWaitStatus() {
	w.waitersCount--
	if w.waitersCount == 0 {
		w.isWait = false
	}
}

func (w *waiter) stopWait() {
	w.waitChan <- struct{}{}
}
