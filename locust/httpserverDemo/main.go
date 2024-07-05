package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a new server
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received request %s - %s - %s - %s\n", time.Now(), r.UserAgent(), r.Header.Get("Connection"), r.Header.Get("Keep-Alive"))
			time.Sleep(200 * time.Millisecond)
			fmt.Fprintln(w, "OK!")
		}),
	}

	// Channel to listen for interrupt or terminate signals from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)
	// signal.Notify(stop, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()
	log.Printf("Server is ready to handle requests at %s", server.Addr)

	// Block until a signal is received
	s := <-stop
	log.Println("get stop signal", s.String())
	log.Println("Server is shutting down...")

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// time.Sleep(45 * time.Second)
	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
