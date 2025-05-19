package main

import (
	"context"
	_ "context"
	_ "fmt"
	"net/http"
	_ "net/http"
	"os"
	"os/signal"
	"time"

	"github.com/shreyansh-ML/movieapp/file/internal/storage/files"
	"github.com/shreyansh-ML/movieapp/file/internal/storage/local"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
// var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
// var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {
	//env.Parse()
	bindAddress := ":9090"
	logLevel := "debug"
	basePath := "./"

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString(logLevel),
		},
	)
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})
	stor, err := local.New(basePath, 1024*1000*5)
	if err != nil {
		l.Error("Unable to create local storage", "error", err)
		return
	}
	handle := files.NewFiles(stor, l)
	sm := mux.NewRouter()
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id}/{filename:[a-zA-Z0-9_\\-\\.]+}", handle.ServeHTTP)

	//get handler
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(basePath))),
	)
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	go func() {
		l.Info("Starting server", "bind_address", bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
