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

func New(testfile string) *ImgServer {
	// Open the file to be streamed
	file, err := os.ReadFile(testfile)
	if err != nil {
		fmt.Print(err)
	}

	return &ImgServer{
		fileBuf: file,
	}
}

func (s *ImgServer) GetImage(req *imagedata.ImageRequest, stream imagedata.ImageService_GetImageServer) error {
	reader := bytes.NewReader(s.fileBuf)

	// how much to send in one stream message
	bufferSize := int64(10 * 1024 * 1024)

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

func NewImageServer(port string, servefile string) (imagedata.ImageServiceServer, func() error, error) {
	// Create a new gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, nil, err
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*10),
		grpc.MaxSendMsgSize(1024*1024*10),
	)
	imgServer := New(servefile)
	imagedata.RegisterImageServiceServer(s, imgServer)

	go func() {
		// Start the server
		log.Println("Server started on port :50051")
		if err := s.Serve(lis); err != nil {
			log.Printf("Closing serverport: %v", err)
		}
		s.GracefulStop()
		lis.Close()
	}()

	return imgServer, lis.Close, nil
}
