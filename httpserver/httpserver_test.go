package httpserver_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/MelleKoning/transferspeed/httpserver"
)

func TestHttpServer(t *testing.T) {
	serverPort := 3333
	testfile := "../data/lionheadbig.jpg"

	cSrv, err := httpserver.New(serverPort, testfile)
	if err != nil {
		t.Fail()
	}
	defer func() {
		err := cSrv.Close()
		if err != nil {
			fmt.Print(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	requestURL := fmt.Sprintf("http://localhost:%d/getimage", serverPort)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)
	// assert that the returned code is 200 OK
	if 200 != res.StatusCode {
		t.Fail()
	}

	// read response body..
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)
	}
	defer res.Body.Close()

	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		t.Fail()
	}

	if !bytes.Equal(body, loadfile) {
		t.Fatalf("not the same")

	}
}
