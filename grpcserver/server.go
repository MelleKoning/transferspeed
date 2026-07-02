package grpcserver

import (
	"bytes"
	"errors"
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
}

func New() *ImgServer {
	return &ImgServer{}
}

func (s *ImgServer) DownloadImages(
	req *imagedata.GetImagesRequest,
	stream imagedata.ImageService_DownloadImagesServer,
) error {
	// how much to send in one stream message
	bufferSize := int64(10 * 1024 * 1024)

	for _, imageID := range req.ImageIds {
		// Read the file
		fileData, err := os.ReadFile(imageID)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", imageID, err)
		}

		// Send header
		header := &imagedata.ImageHeader{
			ImageId:    imageID,
			MimeType:   "image/jpeg", // Assuming jpeg for now based on data/
			TotalBytes: uint64(len(fileData)),
		}
		if err := stream.Send(&imagedata.StreamResponse{
			Payload: &imagedata.StreamResponse_Header{Header: header},
		}); err != nil {
			return err
		}

		// Send chunks
		reader := bytes.NewReader(fileData)
		buffer := make([]byte, bufferSize)

		for {
			n, err := reader.Read(buffer)
			if n > 0 {
				chunk := &imagedata.ImageChunk{
					Data: buffer[:n],
				}
				if err := stream.Send(&imagedata.StreamResponse{
					Payload: &imagedata.StreamResponse_Chunk{Chunk: chunk},
				}); err != nil {
					return err
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
	}

	return nil
}

func NewImageServer(port string) (imagedata.ImageServiceServer, func() error, error) {
	// Create a new gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, nil, err
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024*10),
		grpc.MaxSendMsgSize(1024*1024*10),
	)
	imgServer := New()
	imagedata.RegisterImageServiceServer(s, imgServer)

	go func() {
		// Start the server
		if err := s.Serve(lis); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Printf("Serve returned error: %v", err)
			}
		}
		s.GracefulStop()
		lis.Close()
	}()

	return imgServer, lis.Close, nil
}
