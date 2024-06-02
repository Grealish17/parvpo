package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Grealish17/parvpo/infrastructure/kafka"
	"github.com/Grealish17/parvpo/infrastructure/logger"
	"github.com/Grealish17/parvpo/internal/app/db"
	"github.com/Grealish17/parvpo/internal/app/repository"
	"github.com/Grealish17/parvpo/internal/app/server"
	"github.com/Grealish17/parvpo/internal/app/service"
	"github.com/Grealish17/parvpo/internal/model"
	"github.com/Grealish17/parvpo/internal/receiver"
	"github.com/Grealish17/parvpo/internal/sender"
)

const (
	envFile = "../../prod.env"
)

var brokers = []string{
	"kafka1:9091",
	"kafka2:9092",
	"kafka3:9093",
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	config, err := db.LoadEnv(envFile)
	if err != nil {
		//log.Fatal(err)
		logger.Fatal(err)
	}
	logger.Debug("Load config from environmental")

	database, err := db.NewDb(ctx, config)
	if err != nil {
		//log.Fatal(err)
		logger.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	ticketsRepo := repository.NewTickets(database)

	serv := service.NewTicketsService(ticketsRepo)

	msgChan := make(chan model.Message, 1000)

	kafkaProducer, err := kafka.NewProducer(brokers, kafka.WithMaxOpenRequests(1), kafka.WithRandomPartitioner(), kafka.WaitForAll(),
		kafka.ReturnSuccesses(true), kafka.ReturnErrors(true), kafka.Idempotent(true), kafka.WithCompressionLevelDefault(), kafka.WithCompressionGZIP())
	if err != nil {
		//fmt.Println(err)
		logger.Error(err)
	}
	sender := sender.NewKafkaSender(kafkaProducer, "responses")

	implemetation := server.NewServer(serv, sender, msgChan)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		receiver.ConsumerGroupLogging(brokers, "requests", msgChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		implemetation.Listen(ctx)
	}()

	wg.Wait()
}
