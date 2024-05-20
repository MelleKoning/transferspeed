package httpclient_test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/MelleKoning/transferspeed/httpclient"
	"github.com/MelleKoning/transferspeed/httpserver"
)

func TestServerAndC(t *testing.T) {
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
	fmt.Printf("client: response body: %s\n", resBody)
}
