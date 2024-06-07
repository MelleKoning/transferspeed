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

	serverPort := 3333

	srvPort := fmt.Sprintf(":%d", serverPort)
	closeServer, err := httpserver.New(srvPort, testfile)
	if err != nil {
		t.Fail()
	}
	defer func() {
		err := closeServer.Close()
		if err != nil {
			fmt.Print(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	res := httpclient.GetHTTPClient(serverPort)

	if res.StatusCode != 200 {
		t.Fail()
	}

	// write the returned response
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	// compare bytes response with expected response
	if !(bytes.Equal(resBody, loadfile)) {
		t.Fatalf("not equal")
	}

	_ = os.WriteFile("../data/tmphttpclient.jpg", resBody, os.ModePerm)

}
