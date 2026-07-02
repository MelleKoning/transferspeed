package httpclient

import (
	"fmt"
	"io"
	"net/http"
)

// GetHTTPClient creates the http.Client, the idea
// is that when the GetImage request is done
// this same client is reused instead of having
// to re-initialize a new client for every request
func GetHTTPClient() *http.Client {
	c := &http.Client{}

	return c
}

func GetImages(c *http.Client, serverPort int, imageIDs []string) map[string][]byte {
	results := make(map[string][]byte)

	for _, imageID := range imageIDs {
		requestURL := fmt.Sprintf("http://localhost:%d/getimage?file=%s", serverPort, imageID)

		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			fmt.Printf("error creating request: %s\n", err)
			continue
		}

		req.Proto = "HTTP/2"
		body, err := DoRequest(c, req)
		if err != nil {
			fmt.Printf("error reading response body: %s\n", err)
			continue
		}
		results[imageID] = body
	}

	return results
}

func DoRequest(c *http.Client, req *http.Request) ([]byte, error) {
	res, err := c.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error making http request: %s", err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %s", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)

	}
	body, err := io.ReadAll(res.Body)
	return body, err
}
