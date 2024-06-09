package grpcclient

import (
	"bytes"
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
	conn, err := grpc.Dial(address,
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
func (imgsvcclient *ImageStreamer) GetImage() []byte {
	// Create a request
	request := &imagedata.ImageRequest{
		// Set your request parameters here
	}
	// Make the gRPC request
	stream, err := imgsvcclient.requestclient.GetImage(
		context.Background(),
		request,
	)
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
				//log.Println("image received")
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
