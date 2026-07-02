package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MelleKoning/transferspeed/grpcserver"
)

func main() {
	port := ":50051"
	if envPort := os.Getenv("GRPC_PORT"); envPort != "" {
		port = envPort
	}

	_, _, err := grpcserver.NewImageServer(port)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	exitCtx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Keep the main goroutine alive

	<-exitCtx.Done()

}
