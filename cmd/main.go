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

	"github.com/gorilla/mux"
	"github.com/kavehmz/math/serve"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port    = kingpin.Flag("port", "port").Envar("PORT").Default("8080").Int64()
	timeout = kingpin.Flag("timeout", "timeout in ms").Envar("TIMEOUT").Default("1000").Int64()
	ttl     = kingpin.Flag("ttl", "cache ttl in seconds").Envar("TTL").Default("3600").Int64()
)

func main() {
	ctx := context.Background()
	kingpin.Parse()

	mux := mux.NewRouter()
	mux.PathPrefix("/metrics").Handler(promhttp.Handler())
	mux.PathPrefix("/math").Handler(serve.Math(time.Duration(*timeout)*time.Millisecond, time.Duration(*ttl)*time.Millisecond))

	server := &http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: mux}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Panic("Serve failed", err)
		}
	}()

	log.Println("Listening at", *port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT)
	sig := <-sigChan
	log.Printf("Received signal '%v', shutting down\n", sig)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Handler exited with: %v", err)
	}
}
