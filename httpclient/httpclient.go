package httpclient

import (
	"fmt"
	"net/http"
	"os"
)

func GetHTTPClient(serverPort int) *http.Response {
	requestURL := fmt.Sprintf("http://localhost:%d/getimage", serverPort)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	c := http.Client{
		// in future we try HTTP/2
		//Transport: &http2.Transport{},
	}

	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil
	}
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	return res
}
