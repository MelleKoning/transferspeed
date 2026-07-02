package grpcserver

import (
	"bytes"
	"os"
	"testing"
	"time"

	imagedata "github.com/MelleKoning/transferspeed/api/proto"
	"google.golang.org/grpc"
)

type mockImageService_DownloadImagesServer struct {
	grpc.ServerStream
	receivedData map[string][]byte
	header       *imagedata.ImageHeader
	lastImageID  string
}

func (m *mockImageService_DownloadImagesServer) Send(res *imagedata.StreamResponse) error {
	switch p := res.Payload.(type) {
	case *imagedata.StreamResponse_Header:
		m.header = p.Header
		m.lastImageID = p.Header.ImageId
		m.receivedData[m.lastImageID] = []byte{}
	case *imagedata.StreamResponse_Chunk:
		m.receivedData[m.lastImageID] = append(m.receivedData[m.lastImageID], p.Chunk.Data...)
	}
	return nil
}

func TestDownloadImages(t *testing.T) {
	testfile := "../data/lionheadbig.jpg"
	mockStream := &mockImageService_DownloadImagesServer{
		receivedData: make(map[string][]byte),
	}
	imgServer := New()

	req := &imagedata.GetImagesRequest{
		ImageIds: []string{testfile},
	}

	err := imgServer.DownloadImages(req, mockStream)
	if err != nil {
		t.Errorf("DownloadImages failed: %v", err)
	}

	fileData, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	if !bytes.Equal(fileData, mockStream.receivedData[testfile]) {
		t.Errorf("Expected data does not match received data for %s", testfile)
	}
}

func BenchmarkDownloadImagesTwoFiles(b *testing.B) {
	file1 := "../data/lionface.jpg"
	file2 := "../data/lionheadbig.jpg"
	mockStream := &mockImageService_DownloadImagesServer{
		receivedData: make(map[string][]byte),
	}
	imgServer := New()

	req := &imagedata.GetImagesRequest{
		ImageIds: []string{file1, file2},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mockStream.receivedData = make(map[string][]byte)
		_ = imgServer.DownloadImages(req, mockStream)
	}

	if _, ok := mockStream.receivedData[file1]; !ok {
		b.Errorf("missing file 1 in response")
	}
	if _, ok := mockStream.receivedData[file2]; !ok {
		b.Errorf("missing file 2 in response")
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")
}
