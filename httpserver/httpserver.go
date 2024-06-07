package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

// CustomServer extends http.Server with additional fields
type CustomServer struct {
	*http.Server
	filebuf []byte
}

// setup the http server with custom data
func New(httpport string, testfile string) (*CustomServer, error) {
	mux := http.NewServeMux()

	// Open the file to be streamed
	file, err := os.ReadFile(testfile)
	if err != nil {
		fmt.Print(err)
	}

	// initialize custom http server instance
	cSrv := &CustomServer{
		Server: &http.Server{
			Addr:    httpport,
			Handler: mux,
		},
		filebuf: file,
	}
	mux.HandleFunc("/getimage", cSrv.getImage)

	go func() {
		err := cSrv.Server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
		}
	}()

	return cSrv, nil
}

// getImage writes the reponse to `w` based on the incoming request `r`.
func (c *CustomServer) getImage(w http.ResponseWriter, r *http.Request) {
	fmt.Print("got http image request")
	_, err := w.Write(c.filebuf)
	if err != nil {
		fmt.Printf("can not write data %v", err)
	}
}
