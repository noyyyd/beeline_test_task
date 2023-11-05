package main

import (
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
	defaultTimeout int
	queues         *QueuesController
}

func NewController(defaultTimeout int, queues *QueuesController) *Controller {
	return &Controller{
		defaultTimeout: defaultTimeout,

		queues: queues,
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

		// добавить заранее проверку на полноту

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

		if message.Message != nil {
			if err := c.queues.Push(queueName, *message.Message); err != nil {
				log.Printf("failed add data in queue %s: %v", queueName, err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("%+v", message)

	case http.MethodGet:
		queueName := strings.TrimLeft(r.URL.Path, "/queue/")

		ctx, cancel, err := c.createContext(r)
		defer cancel()
		if err != nil {
			log.Printf("failed create context: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := c.queues.Pop(ctx, queueName)
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

	ctx, cancel := context.WithCancel(context.Background())

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
