package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Grealish17/parvpo/internal/api/server"
	"github.com/Grealish17/parvpo/internal/model"
	"github.com/Grealish17/parvpo/internal/receiver"
)

func runServer(ctx context.Context) error {
	var (
		srv = &http.Server{
			Addr: port,
		}
		wg = sync.WaitGroup{}
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	msgChan := make(chan model.Message, 1000)
	wg.Add(1)
	go func() {
		defer wg.Done()
		receiver.ConsumerGroupLogging(brokers, "responses", msgChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Distribute(msgChan)
	}()

	fmt.Println("START")

	log.Printf("listening on %s", port)
	<-ctx.Done()

	wg.Wait()

	log.Println("shutting down server gracefully")

	shutdownCtx := context.Context(context.Background())

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	return nil
}
