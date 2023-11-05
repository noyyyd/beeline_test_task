package main

type queue struct {
	queue []string
	recvi int
	sendi int
}

func NewQueue(maxQueueSize int) *queue {
	return &queue{
		recvi: -1,
		sendi: -1,
		queue: make([]string, maxQueueSize),
	}
}

func (q *queue) IsFull() bool {
	return (q.sendi == q.recvi+1) || (q.sendi == 0 && q.recvi == len(q.queue)-1)
}

func (q *queue) IsEmpty() bool {
	if q == nil {
		return true
	}

	return q.sendi == -1
}

func (q *queue) Add(data string) {
	if q.sendi == -1 {
		q.sendi = 0
	}
	q.incRecvx()
	q.queue[q.recvi] = data
}

func (q *queue) incRecvx() {
	if q.recvi == len(q.queue)-1 {
		q.recvi = 0
	} else {
		q.recvi++
	}
}

func (q *queue) Get() string {
	defer q.incSendx()
	return q.queue[q.sendi]
}

func (q *queue) incSendx() {
	if q.sendi == q.recvi {
		q.sendi = -1
		q.recvi = -1
	} else if q.sendi == len(q.queue)-1 {
		q.sendi = 0
	} else {
		q.sendi++
	}
}
