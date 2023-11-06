package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var (
		port           int
		defaultTimeout int
		maxQueuesCount int
		maxQueueSize   int
	)

	flag.IntVar(&port, "port", 9999, "")
	flag.IntVar(&defaultTimeout, "default_timeout", -1, "to disable the defaultTimeout, enter -1")
	flag.IntVar(&maxQueuesCount, "max_queues", 0, "")
	flag.IntVar(&maxQueueSize, "max_queue_size", 0, "")
	flag.Parse()

	server := http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
	}

	c := NewController(defaultTimeout, maxQueuesCount, maxQueueSize)

	http.HandleFunc("/queue/", c.Queue)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed start server: %v", err)
	}
}
