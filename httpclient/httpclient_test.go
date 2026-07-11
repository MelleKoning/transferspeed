package httpclient_test

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MelleKoning/transferspeed/httpclient"
	"github.com/MelleKoning/transferspeed/httpserver"
)

const defaultHTTPExternalPort = 8008

// setupHTTPServer starts a local HTTP server unless USE_EXTERNAL_SERVER is set,
// in which case it is a no-op (the test connects to an already-running server).
func setupHTTPServer(tb testing.TB, port int) func() {
	if os.Getenv("USE_EXTERNAL_SERVER") == "true" {
		return func() {}
	}
	httpServer, err := httpserver.New(port)
	if err != nil {
		tb.Fatalf("could not start local server: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	return func() {
		if err := httpServer.Close(); err != nil {
			tb.Errorf("failed to close server: %v", err)
		}
	}
}

// getPort returns the port to connect to: HTTP_SERVER_PORT if set, otherwise defaultPort.
func getPort(defaultPort int) int {
	if p := os.Getenv("HTTP_SERVER_PORT"); p != "" {
		if v, err := strconv.Atoi(strings.TrimPrefix(p, ":")); err == nil {
			return v
		}
	}
	return defaultPort
}

func TestServerAndClient(t *testing.T) {
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	serverPort := getPort(3333)
	defer setupHTTPServer(t, 3333)()

	client := httpclient.GetHTTPClient()
	defer client.CloseIdleConnections()
	results := httpclient.GetImages(client, serverPort, []string{testfile})

	resBody, ok := results[testfile]
	if !ok {
		t.Fatalf("file not found in response")
	}

	if !(bytes.Equal(resBody, loadfile)) {
		t.Fatalf("not equal")
	}

	_ = os.WriteFile("../data/tmphttpclient.jpg", resBody, os.ModePerm)
}

func BenchmarkTesthttpClientWithServer(b *testing.B) {
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		b.Fatal(err)
	}

	serverPort := getPort(3333)
	defer setupHTTPServer(b, 3333)()

	client := httpclient.GetHTTPClient()
	defer client.CloseIdleConnections()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		results := httpclient.GetImages(client, serverPort, []string{testfile})
		resBody := results[testfile]

		if !(bytes.Equal(resBody, loadfile)) {
			b.Fatalf("not equal")
		}
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")
}

func BenchmarkTesthttpClientWithServerTwoFiles(b *testing.B) {
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

	serverPort := getPort(3333)
	defer setupHTTPServer(b, 3333)()

	client := httpclient.GetHTTPClient()
	defer client.CloseIdleConnections()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		results := httpclient.GetImages(client, serverPort, []string{file1, file2})

		if !bytes.Equal(results[file1], loadfile1) {
			b.Fatalf("file1 not equal")
		}
		if !bytes.Equal(results[file2], loadfile2) {
			b.Fatalf("file2 not equal")
		}
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")
}

func BenchmarkTesthttpClient10ImageRequests(b *testing.B) {
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

	serverPort := getPort(3333)
	defer setupHTTPServer(b, 3333)()

	client := httpclient.GetHTTPClient()
	defer client.CloseIdleConnections()

	expected := map[string][]byte{
		file1: loadfile1,
		file2: loadfile2,
	}

	const requestsPerBatch = 10
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		paths := make([]string, 0, requestsPerBatch)
		for i := 0; i < requestsPerBatch/2; i++ {
			paths = append(paths, file1, file2)
		}
		for _, p := range paths {
			results := httpclient.GetImages(client, serverPort, []string{p})
			if got, ok := results[p]; !ok || !bytes.Equal(got, expected[p]) {
				b.Fatalf("unexpected response for %s", p)
			}
		}
	}
	// Per iteration (batch of 10 sequential requests):
	b.ReportMetric(float64(b.Elapsed())/float64(b.N)/1e6, "ms/batch")
	// Amortized wall-clock cost per transferred file (throughput view):
	b.ReportMetric(float64(b.Elapsed())/float64(b.N*requestsPerBatch)/1e6, "ms/file")
}
