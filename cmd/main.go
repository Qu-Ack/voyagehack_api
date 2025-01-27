package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Qu-Ack/voyagehack_api/server"
)

func main() {
	// Set MongoDB environment variables
	os.Setenv("MONGO_URI", "mongodb+srv://dakshsangal:amqp4fJNIIsZvGMW@cluster0.ek1fe.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	os.Setenv("MONGO_DB", "voyage2")

	// Create context that listens for signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start server
	go func() {
		log.Printf("🚀 Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("🛑 Shutting down server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
