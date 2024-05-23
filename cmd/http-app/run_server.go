package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
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
