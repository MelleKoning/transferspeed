package grpcclient_test

import (
	"bytes"
	"os"
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
		t.Fail()
	}

	grpcserver.NewImageServer(":50051", testfile)

	c, closeClient := grpcclient.New(":50051")
	defer closeClient()
	response := c.GetImage()

	// compare bytes response with expected response
	if !(bytes.Equal(response, loadfile)) {
		t.Fatalf("not equal")
	}

	os.WriteFile("../data/tmp.jpg", response, os.ModePerm)
}

func BenchmarkTestClientWithServer(b *testing.B) {
	// load the server
	testfile := "../data/lionheadbig.jpg"
	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		b.Fail()
	}

	port := ":50005"
	_, closeServerFunc, err := grpcserver.NewImageServer(port, testfile)
	if err != nil {
		b.Fatalf("could not start server %v", err)
	}
	defer closeServerFunc()

	//time.Sleep(1 * time.Second)
	c, close := grpcclient.New(port)
	defer close()

	b.ResetTimer()
	var response []byte

	for n := 0; n < b.N; n++ {
		response = c.GetImage()
	}
	// compare bytes response with expected response
	if !(bytes.Equal(response, loadfile)) {
		b.Fatalf("not equal")
	}
	b.ReportMetric(float64(b.Elapsed()/time.Duration(b.N))/float64(1e6), "ms/op")

}
