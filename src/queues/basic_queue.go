package queues

import "log"

type basicQueue struct {
	queue []string
}

func newBasicQueue() *basicQueue {
	return &basicQueue{
		queue: make([]string, 0),
	}
}

func (q *basicQueue) IsFull() bool {
	return false
}

func (q *basicQueue) IsEmpty() bool {
	return len(q.queue) == 0
}

func (q *basicQueue) Add(data string) error {
	q.queue = append(q.queue, data)
	log.Println(q.queue)
	return nil
}

func (q *basicQueue) Get() (string, error) {
	if q.IsEmpty() {
		return "", NotFoundError
	}

	defer func() {
		q.queue = q.queue[1:]
		log.Println(q.queue)
	}()
	log.Println(q.queue)
	return q.queue[0], nil
}
