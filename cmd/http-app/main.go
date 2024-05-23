package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/Grealish17/parvpo/infrastructure/kafka"
	"github.com/Grealish17/parvpo/internal/routes"
	"github.com/Grealish17/parvpo/internal/sender"
	"github.com/Grealish17/parvpo/internal/server"
	"github.com/Grealish17/parvpo/internal/service"
)

const (
	port = ":9000"
)

var brokers = []string{
	"127.0.0.1:9091",
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	kafkaProducer, err := kafka.NewProducer(brokers, kafka.WithMaxOpenRequests(1), kafka.WithRandomPartitioner(), kafka.WaitForAll(),
		kafka.ReturnSuccesses(true), kafka.ReturnErrors(true), kafka.Idempotent(true), kafka.WithCompressionLevelDefault(), kafka.WithCompressionGZIP())
	if err != nil {
		fmt.Println(err)
	}
	sender := sender.NewKafkaSender(kafkaProducer, "logs")

	serv := service.NewService(sender)

	implemetation := server.NewServer(serv)

	http.Handle("/", routes.CreateRouter(implemetation))

	if err := runServer(ctx); err != nil {
		log.Fatal(err)
	}

}
