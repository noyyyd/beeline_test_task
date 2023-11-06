package use_cases

import (
	"beeline_test_task/queues"
	"context"
	"fmt"
	"log"
	"sync"
)

const (
	WithoutLimit = -1
)

type QueueController struct {
	maxQueuesCount int
	maxQueueSize   int

	queues  map[string]queues.Queue
	waiters map[string]*waiter

	// решил использовать RWMutex, а не sync.Map так как по идее тут проблемы cache connection не должно быть
	mx sync.RWMutex
}

func NewQueueController(maxQueuesCount, maxQueueSize int) *QueueController {
	c := new(QueueController)

	if maxQueueSize > 0 {
		c.maxQueueSize = maxQueueSize
	} else {
		c.maxQueueSize = WithoutLimit
	}

	// делаем так чтобы избежать эвакуации данных при добавлении очередей, если заранее задан максимум
	if maxQueuesCount > 0 {
		c.maxQueuesCount = maxQueuesCount
		c.queues = make(map[string]queues.Queue, maxQueuesCount)
		c.waiters = make(map[string]*waiter, maxQueuesCount)
	} else {
		c.maxQueuesCount = WithoutLimit
		c.queues = make(map[string]queues.Queue)
		c.waiters = make(map[string]*waiter)
	}

	return c
}

func (q *QueueController) Pop(ctx context.Context, name string) (string, error) {
	if w, need := q.needWait(name); need {
		w.wait(ctx)
	}

	q.mx.RLock()
	defer q.mx.RUnlock()

	queue, ok := q.queues[name]
	if !ok {
		return "", queues.NotFoundError
	}

	return queue.Get()
}

func (q *QueueController) needWait(name string) (*waiter, bool) {
	q.mx.RLock()
	defer q.mx.RUnlock()

	if q.IsEmpty(name) {
		w, ok := q.waiters[name]
		if !ok {
			w = newWaiter()
			q.waiters[name] = w
		}
		return w, true
	}

	return nil, false
}

func (q *QueueController) Push(name, value string) error {
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

	if existQueue, ok := q.queues[name]; !ok {
		queue := queues.NewQueue(q.maxQueueSize)

		if err := queue.Add(value); err != nil {
			log.Printf("failed add value in queue %s: %v", name, err)
			return err
		}

		q.queues[name] = queue
	} else {
		if err := existQueue.Add(value); err != nil {
			log.Printf("failed add value in queue %s: %v", name, err)
			return err
		}
	}

	return nil
}

func (q *QueueController) IsEmpty(name string) bool {
	q.mx.RLock()
	defer q.mx.RUnlock()

	if existQueue, ok := q.queues[name]; ok && existQueue != nil {
		return existQueue.IsEmpty()
	}

	return true
}

func (q *QueueController) IsFull(name string) bool {
	q.mx.RLock()
	defer q.mx.RUnlock()

	if existQueue, ok := q.queues[name]; ok && existQueue != nil {
		return existQueue.IsFull()
	}

	return false
}
