package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// Only serving a static files for now
func main() {
	var port string
	var wait time.Duration

	flag.StringVar(&port, "PORT", "3000", "the port to listen for requests on")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*25, "the duration for which the server will gracefully wait for existing connections to finish - eg 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	// default home path is "/"
	r.HandleFunc("/api/health", HealthCheckHandler)
	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	r.Use(loggingMiddleware)
	addr := fmt.Sprintf(":%s", port)
	log.Printf("bind address: %s", addr)
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT, or SIGTERM (Ctrl_/) will not be caught
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal
	<-c

	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	srv.Shutdown(ctx)

	// Optionally you could run srv.Shutdown in a goroutine and block on <-ctx.Done() if your applications should wait
	// for other services to finalize based on context cancellation
	log.Println("shutting down")
	os.Exit(0)
}

type HealthCheckStruct struct {
	Alive bool
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	health := HealthCheckStruct{
		Alive: true,
	}
	json.NewEncoder(w).Encode(health)
}
