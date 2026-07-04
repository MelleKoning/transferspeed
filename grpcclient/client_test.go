package grpcclient_test

import (
	"bytes"
	"context"
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
	results, err := c.DownloadImages(context.Background(), []string{testfile})
	if err != nil {
		t.Fatalf("DownloadImages failed: %v", err)
	}

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
		results, err = c.DownloadImages(context.Background(), []string{testfile})
		if err != nil {
			b.Fatalf("DownloadImages failed: %v", err)
		}
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
		results, err := c.DownloadImages(context.Background(), []string{file1, file2})
		if err != nil {
			b.Fatalf("DownloadImages failed: %v", err)
		}

		if !bytes.Equal(results[file1], loadfile1) {
			b.Fatalf("file1 not equal")
		}
		if !bytes.Equal(results[file2], loadfile2) {
			b.Fatalf("file2 not equal")
		}
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")
}

func TestGetImage(t *testing.T) {
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	port := ":50052"
	_, closeServerFunc, err := grpcserver.NewImageServer(port)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := closeServerFunc(); err != nil {
			t.Errorf("failed to close server: %v", err)
		}
	}()

	c, closeClient := grpcclient.New(port)
	defer func() {
		if err := closeClient(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	response, err := c.GetImage(context.Background(), testfile)
	if err != nil {
		t.Fatalf("GetImage failed: %v", err)
	}

	if !(bytes.Equal(response, loadfile)) {
		t.Fatalf("data not equal")
	}
}

func BenchmarkGetImage(b *testing.B) {
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		b.Fatal(err)
	}

	port := ":50007"
	_, closeServerFunc, err := grpcserver.NewImageServer(port)
	if err != nil {
		b.Fatalf("could not start server %v", err)
	}

	defer func() {
		if err := closeServerFunc(); err != nil {
			b.Errorf("failed to close server: %v", err)
		}
	}()

	c, close := grpcclient.New(port)
	defer func() {
		if err := close(); err != nil {
			b.Errorf("failed to close client: %v", err)
		}
	}()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		res, err := c.GetImage(context.Background(), testfile)
		if err != nil {
			b.Fatalf("GetImage failed: %v", err)
		}
		if !(bytes.Equal(res, loadfile)) {
			b.Fatalf("not equal")
		}
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")
}

func BenchmarkGetImageTenParallel(b *testing.B) {
	testfile1 := "../data/lionface.jpg"
	testfile2 := "../data/lionheadbig.jpg"
	loadfile1, err := os.ReadFile(testfile1)
	if err != nil {
		b.Fatal(err)
	}
	loadfile2, err := os.ReadFile(testfile2)
	if err != nil {
		b.Fatal(err)
	}

	port := ":50008"
	_, closeServerFunc, err := grpcserver.NewImageServer(port)
	if err != nil {
		b.Fatalf("could not start server %v", err)
	}
	defer func() {
		if err := closeServerFunc(); err != nil {
			b.Errorf("failed to close server: %v", err)
		}
	}()

	c, close := grpcclient.New(port)
	defer func() {
		if err := close(); err != nil {
			b.Errorf("failed to close client: %v", err)
		}
	}()

	// Create a mix of 10 requests using the 2 available files
	files := []string{testfile1, testfile2, testfile1, testfile2, testfile1, testfile2, testfile1, testfile2, testfile1, testfile2}
	loadFiles := [][]byte{loadfile1, loadfile2, loadfile1, loadfile2, loadfile1, loadfile2, loadfile1, loadfile2, loadfile1, loadfile2}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg sync.WaitGroup
			wg.Add(10)

			for i := 0; i < 10; i++ {
				go func(idx int) {
					defer wg.Done()
					res, err := c.GetImage(context.Background(), files[idx])
					if err != nil {
						b.Errorf("GetImage %d failed: %v", idx, err)
					}
					if !bytes.Equal(res, loadFiles[idx]) {
						b.Errorf("file %d data not equal", idx)
					}
				}(i)
			}

			wg.Wait()
		}
	})
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")
}
