package grpcserver

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	imagedata "github.com/MelleKoning/transferspeed/api/proto"

	"google.golang.org/grpc"
)

type ImgServer struct {
	imagedata.UnimplementedImageServiceServer
	fileBuf []byte // in memory file
}

func New() *ImgServer {
	// Open the file to be streamed
	file, err := os.ReadFile("../data/lionface.jpg")
	if err != nil {
		fmt.Print(err)
	}

	return &ImgServer{
		fileBuf: file,
	}
}

func (s *ImgServer) GetImage(req *imagedata.ImageRequest, stream imagedata.ImageService_GetImageServer) error {
	reader := bytes.NewReader(s.fileBuf)

	// Create a buffer of size 1Mb
	bufferSize := int64(1 * 1024 * 1024) // 1Mb

	// Create a buffer to read data from the file
	buffer := make([]byte, bufferSize)

	// Read data from the file in chunks and send them over the stream
	for {
		// Read a chunk of data from the file
		n, err := reader.Read(buffer)
		if n > 0 {
			// Send the chunk over the stream
			sendErr := stream.Send(&imagedata.ImageChunk{
				ChunkData: buffer[:n],
			})
			if sendErr != nil {
				return sendErr
			}
		}
		if err != nil {
			if err == io.EOF {
				break // Exit the loop if we reach the end of the file
			}

			return err
		}
	}

	return nil
}

func NewImageServer() imagedata.ImageServiceServer {
	// Create a new gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	imgServer := &ImgServer{}
	imagedata.RegisterImageServiceServer(s, imgServer)

	// Start the server
	log.Println("Server started on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	return imgServer
}
