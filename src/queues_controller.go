package main

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type QueuesController struct {
	maxQueuesCount int
	maxQueueSize   int

	queues  map[string]*queue
	waiters map[string]*waiter

	mx sync.RWMutex
}

func NewQueuesController(maxQueuesCount, maxQueueSize int) *QueuesController {
	return &QueuesController{
		maxQueueSize:   maxQueueSize,
		maxQueuesCount: maxQueuesCount,
		queues:         make(map[string]*queue, maxQueuesCount),
		waiters:        make(map[string]*waiter, maxQueuesCount),
	}
}

func (q *QueuesController) Pop(ctx context.Context, name string) (string, error) {
	if waitChan := q.needWait(ctx, name); waitChan != nil {
		<-waitChan
	}

	q.mx.RLock()
	defer q.mx.RUnlock()

	if q.queues[name].IsEmpty() {
		return "", fmt.Errorf("message not found")
	}

	log.Println(q.queues[name].sendi, q.queues[name].recvi)
	log.Println(q.queues[name].queue, len(q.queues[name].queue))

	return q.queues[name].Get(), nil
}

func (q *QueuesController) needWait(ctx context.Context, name string) chan struct{} {
	q.mx.RLock()
	defer q.mx.RUnlock()

	if q.queues[name].IsEmpty() {
		w, ok := q.waiters[name]
		if !ok {
			w = newWaiter()
			w.wait(ctx)
			q.waiters[name] = w
		} else {
			w.wait(ctx)
		}
		return w.waitChan
	}

	return nil
}

func (q *QueuesController) Push(name, value string) error {
	q.mx.Lock()
	defer func() {
		if w, ok := q.waiters[name]; ok && w.isWait {
			w.stopWait()
		}
		q.mx.Unlock()
	}()

	if len(q.queues) == q.maxQueuesCount {
		return fmt.Errorf("limit on the number of queues(%d) has been reached", q.maxQueuesCount)
	}

	if _, ok := q.queues[name]; !ok {
		q.queues[name] = NewQueue(q.maxQueueSize)
		q.queues[name].Add(value)
	} else {
		if q.queues[name].IsFull() {
			return fmt.Errorf("limit of messages(%d) in the %s queue has been reached", q.maxQueuesCount, name)
		}

		q.queues[name].Add(value)
	}

	return nil
}
