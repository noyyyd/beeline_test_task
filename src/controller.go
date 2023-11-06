package main

import (
	"beeline_test_task/use_cases"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	defaultTimeout   int
	queuesController *use_cases.QueueController
}

func NewController(defaultTimeout, maxQueuesCount, maxQueueSize int) *Controller {
	return &Controller{
		defaultTimeout: defaultTimeout,

		queuesController: use_cases.NewQueueController(maxQueuesCount, maxQueueSize),
	}
}

type Message struct {
	// Message здесь является указателем чтобы мы могли проверить по формату тело или нет
	// без этого мы не сможем понять отправили нам "" или поле приняло это значение по умолчанию
	Message *string `json:"message"`
}

func (c *Controller) Queue(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		queueName := strings.TrimLeft(r.URL.Path, "/queue/")

		// заранее проверяем полна ли очередь, чтобы не приходилось обрабатывать заранее неуспешлый запрос
		if c.queuesController.IsFull(queueName) {
			w.WriteHeader(http.StatusBadRequest)
		}

		message := new(Message)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed read req body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, message); err != nil {
			log.Printf("failed unmarshal body %v: %v", body, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if message.Message == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := c.queuesController.Push(queueName, *message.Message); err != nil {
			log.Printf("failed add data in queue %s: %v", queueName, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodGet:
		queueName := strings.TrimLeft(r.URL.Path, "/queue/")

		ctx, cancel, err := c.createContext(r)
		defer func() {
			cancel()
		}()
		if err != nil {
			log.Printf("failed create context: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := c.queuesController.Pop(ctx, queueName)
		if err != nil {
			log.Printf("failed get data from queue %s: %v", queueName, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		message := &Message{
			Message: &data,
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("failed marshal data %v: %v", message, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(messageBytes); err != nil {
			log.Printf("failed write resp: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) createContext(r *http.Request) (context.Context, context.CancelFunc, error) {
	timeout := c.defaultTimeout

	ctx, cancel := context.WithCancel(r.Context())

	timeoutStr := r.URL.Query().Get("timeout")
	if timeoutStr != "" {
		var err error
		timeout, err = strconv.Atoi(timeoutStr)
		if err != nil {
			log.Printf("failed get timeout value from %v", timeoutStr)
			return ctx, cancel, err
		}
	}

	if timeout != -1 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	}

	return ctx, cancel, nil
}
