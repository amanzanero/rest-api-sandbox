package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if len(port) == 1 {
		port = ":8080"
	}
	r := chi.NewRouter()
	done := make(chan os.Signal, 1)

	// setup routes here
	r.Get("/", helloWorld)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("server started on %s", port)

	<-done
	gracefullyShutdown(srv)
}

func gracefullyShutdown(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("server exited properly")
}
