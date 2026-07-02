package grpcclient

import (
	"context"
	"io"
	"log"

	imagedata "github.com/MelleKoning/transferspeed/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ImageStreamer struct {
	conn          *grpc.ClientConn
	requestclient imagedata.ImageServiceClient
}

func New(address string) (*ImageStreamer, func() error) {
	// Set up a connection to the server
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10)))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	requestclient := imagedata.NewImageServiceClient(conn)

	c := &ImageStreamer{
		conn:          conn,
		requestclient: requestclient,
	}

	return c, c.Close
}

func (imgsvcclient *ImageStreamer) Close() error {
	return imgsvcclient.conn.Close()
}
func (imgsvcclient *ImageStreamer) DownloadImages(imageIDs []string) map[string][]byte {
	// Create a request
	request := &imagedata.GetImagesRequest{
		ImageIds: imageIDs,
	}
	// Make the gRPC request
	stream, err := imgsvcclient.requestclient.DownloadImages(
		context.Background(),
		request,
	)
	if err != nil {
		log.Fatalf("Failed to get images: %v", err)
	}

	results := make(map[string][]byte)
	var currentImageID string

	// Iterate over the response stream
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Failed to receive image chunk: %v", err)
		}

		switch p := response.Payload.(type) {
		case *imagedata.StreamResponse_Header:
			currentImageID = p.Header.ImageId
			results[currentImageID] = []byte{}
		case *imagedata.StreamResponse_Chunk:
			if currentImageID == "" {
				log.Fatal("Received chunk before header")
			}
			results[currentImageID] = append(results[currentImageID], p.Chunk.Data...)
		}
	}

	return results
}
