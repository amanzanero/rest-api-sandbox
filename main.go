package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/{userId}", handler).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Print("Server Started on 8080")

	<-done
	log.Print("Server Stopped")
	gracefullyShutdown(srv)

}

func handler(w http.ResponseWriter, r *http.Request) {
	// mux.Vars(r) returns all values captured in the request URL.
	vars := mux.Vars(r)

	_, err := fmt.Fprintf(w, "User %s\n", vars["userId"])
	if err != nil {
		fmt.Println(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
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
	log.Print("Server Exited Properly")
}