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
	flag.IntVar(&defaultTimeout, "defaultTimeout", -1, "to disable the defaultTimeout, enter -1")
	flag.IntVar(&maxQueuesCount, "max_queues", 5, "")
	flag.IntVar(&maxQueueSize, "max_messages", 3, "")
	flag.Parse()

	server := http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
	}

	c := NewController(defaultTimeout, NewQueuesController(maxQueuesCount, maxQueueSize))

	http.HandleFunc("/queue/", c.Queue)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed start server: %v", err)
	}
}
