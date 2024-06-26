package webserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivmello/kakebo-go-api/internal/provider"
)

type webserver struct {
	provider *provider.Provider
}

func New(provider *provider.Provider) *webserver {
	return &webserver{
		provider,
	}
}

func (w *webserver) Start() {
	mux := http.NewServeMux()
	router := w.registerRoutes(mux)

	server := &http.Server{
		Addr:         ":9999",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on :9999")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :9999: %v\n", err)
		}
	}()

	<-done
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
