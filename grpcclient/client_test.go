package grpcclient_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/MelleKoning/transferspeed/grpcclient"
	"github.com/MelleKoning/transferspeed/grpcserver"
)

func setupServer(t testing.TB, port string) func() {
	if os.Getenv("USE_EXTERNAL_SERVER") == "true" {
		return func() {}
	}
	_, closeFunc, err := grpcserver.NewImageServer(port)
	if err != nil {
		t.Fatalf("could not start local server: %v", err)
	}
	return func() {
		if err := closeFunc(); err != nil {
			t.Errorf("failed to close server: %v", err)
		}
	}
}

func getAddr(port string) string {
	if addr := os.Getenv("GRPC_SERVER_ADDR"); addr != "" {
		return addr
	}
	if p := os.Getenv("GRPC_SERVER_PORT"); p != "" {
		return p
	}
	return port
}

func getDataPath(filename string) string {
	baseDir := os.Getenv("DATA_DIR")
	if baseDir == "" {
		baseDir = "../data"
	}
	return filepath.Join(baseDir, filename)
}

func resolveRemotePath(localPath string) string {
	if os.Getenv("USE_EXTERNAL_SERVER") == "true" {
		filename := filepath.Base(localPath)
		return "/data/" + filename
	}
	return localPath
}

func TestClientWithServer(t *testing.T) {
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	port := ":50051"
	cleanup := setupServer(t, port)
	defer cleanup()

	c, closeClient := grpcclient.New(getAddr(port))
	defer closeClient()
	remotePath := resolveRemotePath(testfile)
	results, err := c.DownloadImages(context.Background(), []string{remotePath})
	if err != nil {
		t.Fatalf("DownloadImages failed: %v", err)
	}

	response, ok := results[remotePath]
	if !ok {
		t.Fatalf("file not found in response, expected %s", remotePath)
	}

	if !(bytes.Equal(response, loadfile)) {
		t.Fatalf("not equal")
	}

	_ = os.WriteFile("../data/tmp.jpg", response, os.ModePerm)
}

func BenchmarkTestgRPCClientWithServer(b *testing.B) {
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		b.Fatal(err)
	}

	port := ":50005"
	cleanup := setupServer(b, port)
	defer cleanup()

	c, close := grpcclient.New(getAddr(port))
	defer func() {
		err := close()
		if err != nil {
			fmt.Printf("grpc close:%v", err)
		}
	}()

	var m sync.Map

	// Pre-validation to avoid validation in benchmark loop
	remoteFile := resolveRemotePath(testfile)
	res, err := c.DownloadImages(context.Background(), []string{remoteFile})
	if err != nil || !bytes.Equal(res[remoteFile], loadfile) {
		b.Fatal("Pre-benchmark validation failed")
	}

	b.ResetTimer()
	var results map[string][]byte
	var n int
	for n = 0; n < b.N; n++ {
		results, err = c.DownloadImages(context.Background(), []string{remoteFile})
		if err != nil {
			b.Fatalf("DownloadImages failed: %v", err)
		}
		m.Store(n, results[remoteFile])
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")

	for contenttest := 0; contenttest < n; contenttest++ {
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
	cleanup := setupServer(b, port)
	defer cleanup()

	c, close := grpcclient.New(getAddr(port))
	defer func() {
		err := close()
		if err != nil {
			fmt.Printf("grpc close:%v", err)
		}
	}()

	remoteFile1 := resolveRemotePath(file1)
	remoteFile2 := resolveRemotePath(file2)

	// Pre-validation
	res, err := c.DownloadImages(context.Background(), []string{remoteFile1, remoteFile2})
	if err != nil || !bytes.Equal(res[remoteFile1], loadfile1) || !bytes.Equal(res[remoteFile2], loadfile2) {
		b.Fatal("Pre-benchmark validation failed")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		results, err := c.DownloadImages(context.Background(), []string{remoteFile1, remoteFile2})
		if err != nil {
			b.Fatalf("DownloadImages failed: %v", err)
		}

		if !bytes.Equal(results[remoteFile1], loadfile1) {
			b.Fatalf("file1 not equal")
		}
		if !bytes.Equal(results[remoteFile2], loadfile2) {
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

	port := ":50051"
	cleanup := setupServer(t, port)
	defer cleanup()

	c, closeClient := grpcclient.New(getAddr(port))
	defer func() {
		if err := closeClient(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	response, err := c.GetImage(context.Background(), resolveRemotePath(testfile))
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
	cleanup := setupServer(b, port)
	defer cleanup()

	c, close := grpcclient.New(getAddr(port))
	defer func() {
		if err := close(); err != nil {
			b.Errorf("failed to close client: %v", err)
		}
	}()

	// 1. Perform one verification BEFORE the loop
	res, err := c.GetImage(context.Background(), resolveRemotePath(testfile))
	if err != nil || !bytes.Equal(res, loadfile) {
		b.Fatal("Initial verification failed")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := c.GetImage(context.Background(), resolveRemotePath(testfile))
		if err != nil {
			b.Fatalf("GetImage failed: %v", err)
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
	cleanup := setupServer(b, port)
	defer cleanup()

	c, close := grpcclient.New(getAddr(port))
	defer func() {
		if err := close(); err != nil {
			b.Errorf("failed to close client: %v", err)
		}
	}()

	// Pre-validation
	res1, err := c.GetImage(context.Background(), resolveRemotePath(testfile1))
	if err != nil || !bytes.Equal(res1, loadfile1) {
		b.Fatal("Pre-benchmark validation failed for testfile1")
	}
	res2, err := c.GetImage(context.Background(), resolveRemotePath(testfile2))
	if err != nil || !bytes.Equal(res2, loadfile2) {
		b.Fatal("Pre-benchmark validation failed for testfile2")
	}

	// Create a mix of 10 requests using the 2 available files
	const filesPerBatch = 10
	files := []string{testfile1, testfile2, testfile1, testfile2, testfile1, testfile2, testfile1, testfile2, testfile1, testfile2}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg sync.WaitGroup
			wg.Add(filesPerBatch)

			for i := 0; i < filesPerBatch; i++ {
				go func(idx int) {
					defer wg.Done()
					_, err := c.GetImage(context.Background(), resolveRemotePath(files[idx]))
					if err != nil {
						b.Errorf("GetImage %d failed: %v", idx, err)
					}
				}(i)
			}

			wg.Wait()
		}
	})
	// Per-request latency under 10-way concurrency (wall time per batch):
	b.ReportMetric(float64(b.Elapsed())/float64(b.N)/1e6, "ms/batch")
	// Amortized wall-clock cost per transferred file (throughput view):
	b.ReportMetric(float64(b.Elapsed())/float64(b.N*filesPerBatch)/1e6, "ms/file")
}
