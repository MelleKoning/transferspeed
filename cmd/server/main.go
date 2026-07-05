package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MelleKoning/transferspeed/grpcserver"
	"github.com/MelleKoning/transferspeed/httpserver"
	"github.com/kouhin/envflag"
)

var httport = flag.Int("httport", 8008, "HTTP server port")

func main() {
	if err := envflag.Parse(); err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	port := ":50051"
	if envPort := os.Getenv("GRPC_PORT"); envPort != "" {
		port = envPort
	}

	_, _, err := grpcserver.NewImageServer(port)
	if err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}

	// Start HTTP server
	httpSrv, err := httpserver.New(*httport)
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}

	exitCtx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Keep the main goroutine alive
	<-exitCtx.Done()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	shutDownErr := httpSrv.Shutdown(shutdownCtx)
	if shutDownErr != nil {
		log.Fatalf("failed to shutdown http server: %v", shutDownErr)
	}
}
