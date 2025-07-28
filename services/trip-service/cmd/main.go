package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
)

func main() {
	inmemRepo := repository.NewInmenRepository()
	tripService := service.NewService(inmemRepo)
	mux := http.NewServeMux()

	httpHandler := h.HttpHandler{
		Service: tripService,
	}

	mux.HandleFunc("POST /preview", httpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server listening on port %s\n", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Server error: %v\n", err)
	case sig := <-shutdown:
		log.Printf("Shutting down server due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v\n", err)
			server.Close()
		} else {
			log.Println("Server gracefully stopped")
		}
		os.Exit(0)
	}
}
