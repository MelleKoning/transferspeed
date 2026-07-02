package main

import (
	"flag"
	"log"
	"os"

	"github.com/MelleKoning/transferspeed/grpcclient"
)

func main() {
	var (
		serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
		outputFile = flag.String("out", "output.jpg", "The path to save the downloaded image")
	)
	flag.Parse()

	addr := *serverAddr
	if envAddr := os.Getenv("GRPC_SERVER_ADDRESS"); envAddr != "" {
		addr = envAddr
	}

	log.Printf("Connecting to gRPC server at %s", addr)
	client, closeClient := grpcclient.New(addr)
	defer func() {
		if err := closeClient(); err != nil {
			log.Printf("Failed to close client connection: %v", err)
		}
	}()

	log.Println("Requesting image from server...")
	// Updated to handle multiple images, defaulting to one as per original cmd/client usage
	imageMap := client.DownloadImages([]string{"../data/lionheadbig.jpg"})
	imageBytes := imageMap["../data/lionheadbig.jpg"]
	log.Printf("Received image of size: %d bytes", len(imageBytes))

	if err := os.WriteFile(*outputFile, imageBytes, 0644); err != nil {
		log.Fatalf("Failed to write image to file %s: %v", *outputFile, err)
	}
	log.Printf("Successfully saved image to %s", *outputFile)

}
