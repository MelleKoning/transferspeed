package httpclient_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/MelleKoning/transferspeed/httpclient"
	"github.com/MelleKoning/transferspeed/httpserver"
)

func TestServerAndClient(t *testing.T) {
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	serverPort := 3333
	httpServer, err := httpserver.New(serverPort)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := httpServer.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

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

	serverPort := 3333
	httpServer, err := httpserver.New(serverPort)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		err := httpServer.Close()
		if err != nil {
			b.Fatal(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

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

	serverPort := 3333
	httpServer, err := httpserver.New(serverPort)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		err := httpServer.Close()
		if err != nil {
			b.Fatal(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

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
