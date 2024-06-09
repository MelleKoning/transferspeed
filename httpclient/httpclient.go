package httpclient

import (
	"fmt"
	"net/http"
	"os"
)

// GetHTTPClient creates the http.Client, the idea
// is that when the GetImage request is done
// this same client is reused instead of having
// to re-initialize a new client for every request
func GetHTTPClient() *http.Client {
	c := &http.Client{}

	return c
}

func GetImage(c *http.Client, serverPort int) *http.Response {
	requestURL := fmt.Sprintf("http://localhost:%d/getimage", serverPort)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	req.Proto = "HTTP/2"
	res, err := c.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil
	}
	//fmt.Printf("client: got response!\n")
	//fmt.Printf("client: status code: %d\n", res.StatusCode)

	return res
}
