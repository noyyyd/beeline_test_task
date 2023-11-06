package queues

type circularQueue struct {
	queue []string
	recvi int
	sendi int
}

func newCircularQueue(maxQueueSize int) *circularQueue {
	return &circularQueue{
		recvi: -1,
		sendi: -1,
		queue: make([]string, maxQueueSize),
	}
}

func (q *circularQueue) IsFull() bool {
	return (q.sendi == q.recvi+1) || (q.sendi == 0 && q.recvi == len(q.queue)-1)
}

func (q *circularQueue) IsEmpty() bool {
	return q.sendi == -1
}

func (q *circularQueue) Add(data string) error {
	if q.IsFull() {
		return LimitError
	}

	if q.sendi == -1 {
		q.sendi = 0
	}
	q.incRecvi()
	q.queue[q.recvi] = data

	return nil
}

func (q *circularQueue) Get() (string, error) {
	if q.IsEmpty() {
		return "", NotFoundError
	}

	defer q.incSendi()
	return q.queue[q.sendi], nil
}

func (q *circularQueue) incRecvi() {
	q.recvi++
}

func (q *circularQueue) incSendi() {
	if q.sendi == q.recvi {
		q.sendi = -1
		q.recvi = -1
	} else if q.sendi == len(q.queue)-1 {
		q.sendi = 0
	} else {
		q.sendi++
	}
}
