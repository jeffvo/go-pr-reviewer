package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jeffvo/go-pr-reviewer/handlers"
	"github.com/jeffvo/go-pr-reviewer/internals/adapters"
	"github.com/jeffvo/go-pr-reviewer/usecases"
)

func main() {
	mux := http.NewServeMux()

	githubAdapter := adapters.NewGithubAdapter("")
	geminiAdapter := adapters.NewGeminiAdapter("")
	usecase := usecases.NewWebhookProcessor(githubAdapter, geminiAdapter)
	webhookHandler := handlers.NewWebhookHandler(usecase)

	mux.HandleFunc("/", webhookHandler.Handle)

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Starting server on :3000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :3000: %v\n", err)
		}
	}()

	<-stop // Wait for interrupt signal
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
