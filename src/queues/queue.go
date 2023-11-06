package queues

import "fmt"

var (
	NotFoundError = fmt.Errorf("message not found")
	LimitError    = fmt.Errorf("limit of messages in the queue has been reached")
)

type Queue interface {
	Add(data string) error
	Get() (string, error)
	IsFull() bool
	IsEmpty() bool
}

func NewQueue(maxQueueSize int) Queue {
	if maxQueueSize > 0 {
		return newCircularQueue(maxQueueSize)
	}
	return newBasicQueue()
}
