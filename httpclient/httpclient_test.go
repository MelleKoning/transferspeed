package httpclient_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/MelleKoning/transferspeed/httpclient"
	"github.com/MelleKoning/transferspeed/httpserver"
)

func TestServerAndC(t *testing.T) {
	// arrange loading the testfile for comparison
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		t.Fail()
	}

	// arrange a listening server
	serverPort := 3333
	httpServer, err := httpserver.New(serverPort, testfile)
	if err != nil {
		t.Fail()
	}
	defer func() {
		err := httpServer.Close()
		if err != nil {
			fmt.Print(err)
		}
	}()

	// provide the server some time to start
	// listening before doing a client call
	time.Sleep(100 * time.Millisecond)

	// Act - do the http call
	client := httpclient.GetHTTPClient()
	defer client.CloseIdleConnections()
	res := httpclient.GetImage(client, serverPort)

	if res.StatusCode != 200 {
		t.Fail()
	}

	// write the returned response
	resBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	// compare bytes response with expected response
	if !(bytes.Equal(resBody, loadfile)) {
		t.Fatalf("not equal")
	}

	// for human check of the received file
	// write out the file to disk
	_ = os.WriteFile("../data/tmphttpclient.jpg", resBody, os.ModePerm)

}

func BenchmarkTesthttpClientWithServer(b *testing.B) {
	// load the server
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		b.Fail()
	}

	// arrange a listening server
	serverPort := 3333
	httpServer, err := httpserver.New(serverPort, testfile)
	if err != nil {
		b.Fail()
	}
	defer func() {
		err := httpServer.Close()
		if err != nil {
			fmt.Print(err)
		}
	}()

	// provide the server some time to start
	// listening before doing a client call
	time.Sleep(100 * time.Millisecond)

	// Act - create client for reuse
	client := httpclient.GetHTTPClient()
	defer client.CloseIdleConnections()

	b.ResetTimer()
	var resBody []byte

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			//	for n := 0; n < b.N; n++ {
			res := httpclient.GetImage(client, serverPort)

			if res.StatusCode != 200 {
				b.Fail()
			}

			// write the returned response
			resBody, err = io.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				fmt.Printf("client: could not read response body: %s\n", err)
				os.Exit(1)
			}

		}
	}) // parallel

	// compare bytes response with expected response
	if !(bytes.Equal(resBody, loadfile)) {
		b.Fatalf("not equal")
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")

}
