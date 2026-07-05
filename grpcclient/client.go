package grpcclient

import (
	"context"
	"io"
	"log"

	imagedata "github.com/MelleKoning/transferspeed/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const bufferSize = 1 * 1024 * 1024 // 1 MB
const maxMessageSize = bufferSize + 64*1024

type ImageStreamer struct {
	conn          *grpc.ClientConn
	requestclient imagedata.ImageServiceClient
}

func New(address string) (*ImageStreamer, func() error) {
	// Set up a connection to the server
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)))
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
func (imgsvcclient *ImageStreamer) GetImage(ctx context.Context, imageID string) ([]byte, error) {
	request := &imagedata.GetImageRequest{
		ImageId: imageID,
	}

	stream, err := imgsvcclient.requestclient.GetImage(
		ctx,
		request,
	)
	if err != nil {
		return nil, err
	}

	var data []byte
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		data = append(data, chunk.Data...)
	}

	return data, nil
}

func (imgsvcclient *ImageStreamer) DownloadImages(ctx context.Context, imageIDs []string) (map[string][]byte, error) {
	// Create a request
	request := &imagedata.GetImagesRequest{
		ImageIds: imageIDs,
	}
	// Make the gRPC request
	stream, err := imgsvcclient.requestclient.DownloadImages(
		ctx,
		request,
	)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		switch p := response.Payload.(type) {
		case *imagedata.StreamResponse_Header:
			currentImageID = p.Header.ImageId
			results[currentImageID] = []byte{}
		case *imagedata.StreamResponse_Chunk:
			if currentImageID == "" {
				return nil, io.ErrUnexpectedEOF
			}
			results[currentImageID] = append(results[currentImageID], p.Chunk.Data...)
		}
	}

	return results, nil
}
