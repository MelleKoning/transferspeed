package httpserver_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/MelleKoning/transferspeed/httpserver"
)

func TestHttpServer(t *testing.T) {
	serverPort := 3333

	srvPort := fmt.Sprintf(":%d", serverPort)
	closeServer, err := httpserver.New(srvPort)
	if err != nil {
		t.Fail()
	}
	defer func() {
		err := closeServer()
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
}
