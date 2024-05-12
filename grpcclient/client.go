package grpcclient

import (
	"bytes"
	"context"
	"io"
	"log"

	imagedata "github.com/MelleKoning/transferspeed/api/proto"
	"google.golang.org/grpc"
)

type ImageSvcClient struct {
	client imagedata.ImageServiceClient
}

func New(address string) *ImageSvcClient {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(nil))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	client := imagedata.NewImageServiceClient(conn)
	return &ImageSvcClient{
		client: client,
	}
}

func (imgsvcclient *ImageSvcClient) GetImage() []byte {

	// Create a request
	request := &imagedata.ImageRequest{
		// Set your request parameters here
	}

	// Make the gRPC request
	stream, err := imgsvcclient.client.GetImage(context.Background(), request)
	if err != nil {
		log.Fatalf("Failed to get image: %v", err)
	}

	// Create a buffer to store the concatenated chunks
	var buffer bytes.Buffer

	// Iterate over the response stream
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Failed to receive image chunk: %v", err)
		}

		// Append the received chunk to the buffer
		if _, err := buffer.Write(response.ChunkData); err != nil {
			log.Fatalf("Failed to append chunk data: %v", err)
		}

	}

	return buffer.Bytes()
}
