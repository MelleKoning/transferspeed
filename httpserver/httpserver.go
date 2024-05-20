package httpserver

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// getImage writes the reponse to `w` based on the incoming request `r`.
func getImage(w http.ResponseWriter, r *http.Request) {
	fmt.Print("got http image request")
	io.WriteString(w, "this is my website")
}

// setup the http server
func New(httpport string) (func() error, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/getimage", getImage)

	// initialize http server instance
	server := &http.Server{
		Addr:    httpport,
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)

		}
	}()

	return server.Close, nil
}
