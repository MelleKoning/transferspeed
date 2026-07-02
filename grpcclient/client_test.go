package grpcclient_test

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/MelleKoning/transferspeed/grpcclient"
	"github.com/MelleKoning/transferspeed/grpcserver"
)

func TestClientWithServer(t *testing.T) {
	// load the server
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = grpcserver.NewImageServer(":50051")
	if err != nil {
		t.Fatal(err)
	}

	c, closeClient := grpcclient.New(":50051")
	defer closeClient()
	results := c.DownloadImages([]string{testfile})

	response, ok := results[testfile]
	if !ok {
		t.Fatalf("file not found in response")
	}

	// compare bytes response with expected response
	if !(bytes.Equal(response, loadfile)) {
		t.Fatalf("not equal")
	}

	os.WriteFile("../data/tmp.jpg", response, os.ModePerm)
}

func BenchmarkTestgRPCClientWithServer(b *testing.B) {
	// load the server
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		b.Fatal(err)
	}

	port := ":50005"
	_, closeServerFunc, err := grpcserver.NewImageServer(port)
	if err != nil {
		b.Fatalf("could not start server %v", err)
	}
	defer closeServerFunc()

	c, close := grpcclient.New(port)
	defer func() {
		err := close()
		if err != nil {
			fmt.Printf("grpc close:%v", err)
		}
	}()

	var m sync.Map

	b.ResetTimer()
	var results map[string][]byte

	var n int
	for n = 0; n < b.N; n++ {
		results = c.DownloadImages([]string{testfile})
		m.Store(n, results[testfile])
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")

	for contenttest := 0; contenttest < n; contenttest++ {
		// compare bytes response with expected response
		bytesReceived, exists := m.Load(contenttest)
		if !exists {
			fmt.Printf("contents not found for %d?", contenttest)
		}
		res, converted := bytesReceived.([]byte)
		if !converted {
			fmt.Printf("cannot convert byte array")
		}
		if !(bytes.Equal(res, loadfile)) {
			b.Fatalf("not equal")
		}
	}

}

func BenchmarkTestgRPCClientWithServerTwoFiles(b *testing.B) {
	file1 := "../data/lionface.jpg"
	file2 := "../data/lionheadbig.jpg"
	loadfile1, err := os.ReadFile(file1)
	if err != nil {
		b.Fatal(err)
	}
	loadfile2, err := os.ReadFile(file2)
	if err != nil {
		b.Fatal(err)
	}

	port := ":50006"
	_, closeServerFunc, err := grpcserver.NewImageServer(port)
	if err != nil {
		b.Fatalf("could not start server %v", err)
	}
	defer func() {
		err := closeServerFunc()
		if err != nil {
			b.Fatal(err)
		}
	}()

	c, close := grpcclient.New(port)
	defer func() {
		err := close()
		if err != nil {
			fmt.Printf("grpc close:%v", err)
		}
	}()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		results := c.DownloadImages([]string{file1, file2})

		if !bytes.Equal(results[file1], loadfile1) {
			b.Fatalf("file1 not equal")
		}
		if !bytes.Equal(results[file2], loadfile2) {
			b.Fatalf("file2 not equal")
		}
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")
}
