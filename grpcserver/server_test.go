package grpcserver

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	imagedata "github.com/MelleKoning/transferspeed/api/proto"
)

type mockImageService_GetImageServer struct {
	imagedata.ImageService_GetImageServer
	receivedData []byte
}

func (m *mockImageService_GetImageServer) Send(chunk *imagedata.ImageChunk) error {
	m.receivedData = append(m.receivedData, chunk.ChunkData...)
	return nil
}

// In your test function
func TestGetImage(t *testing.T) {
	testfile := "../data/lionheadbig.jpg"
	mockStream := &mockImageService_GetImageServer{}
	imgServer := New(testfile)

	// Assuming you have a way to set the file path or mock os.Open
	err := imgServer.GetImage(nil, mockStream)
	if err != nil {
		t.Errorf("GetImage failed: %v", err)
	}

	// the received data should now match the file contents
	fileData, err := os.ReadFile(testfile)
	if err != nil {
		fmt.Print(err)
	}

	if !bytes.Equal(fileData, mockStream.receivedData) {
		t.Errorf("Expected data does not match received data")
	}
}

func BenchmarkGetImage(b *testing.B) {
	testfile := "../data/lionheadbig.jpg"

	mockStream := &mockImageService_GetImageServer{}
	imgServer := New(testfile)

	// the received data should now match the file contents
	fileData, err := os.ReadFile(testfile)
	if err != nil {
		fmt.Print(err)
	}

	b.ResetTimer()
	// Assuming you have a way to set the file path or mock os.Open
	for n := 0; n < b.N; n++ {
		// each run we have to reset the received stream to check
		mockStream.receivedData = nil

		_ = imgServer.GetImage(nil, mockStream)
	}
	// the check on read bytes is only done once to keep it out of the timing
	// of the benchmark test loop
	if !bytes.Equal(fileData, mockStream.receivedData) {
		b.Errorf("Expected data does not match received data")
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")
}
