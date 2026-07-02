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

	cSrv, err := httpserver.New(serverPort)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := cSrv.Close()
		if err != nil {
			fmt.Print(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	requestURL := fmt.Sprintf("http://localhost:%d/getimage?file=%s", serverPort, testfile)
	res, err := http.Get(requestURL)
	if err != nil {
		t.Fatalf("error making http request: %s", err)
	}

	// assert that the returned code is 200 OK
	if http.StatusOK != res.StatusCode {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	// read response body..
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}
	defer res.Body.Close()

	loadfile, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(body, loadfile) {
		t.Fatalf("not the same")
	}
}
