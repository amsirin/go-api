package main

import (
	"context"
	"example/api/handlers"
	"example/api/version"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Printf("Starting the service...\ncommit: %s, build time: %s, release: %s", version.Commit, version.BuildTime, version.Release)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not set.")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	r := handlers.Router(version.BuildTime, version.Commit, version.Release)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// this channel is for graceful shutdown:
	// if we receive an error, we can send it here to notify the server to be stopped
	shutdown := make(chan struct{}, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			shutdown <- struct{}{}
			log.Printf("%v", err)
		}
	}()
	log.Print("The services is ready to listen and serve.")

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			log.Print("Got SIGINT...")
		case syscall.SIGTERM:
			log.Print("Got SIGTERM")
		}
	case <-shutdown:
		log.Printf("Got an error...")
	}

	log.Print("Ther service is shutting down")
	srv.Shutdown(context.Background())
	log.Print("Done")

}
