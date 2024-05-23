package main

import (
	"sync"

	"github.com/Grealish17/parvpo/internal/logger"
)

var brokers = []string{
	"127.0.0.1:9091",
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.ConsumerGroupLogging(brokers)
	}()
	wg.Wait()
}
