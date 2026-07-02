package httpserver

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// CustomServer extends http.Server with additional fields
type CustomServer struct {
	*http.Server
}

// setup the http server
func New(httpport int) (*CustomServer, error) {
	mux := http.NewServeMux()

	// initialize custom http server instance
	cSrv := &CustomServer{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", httpport),
			Handler: mux,
		},
	}
	mux.HandleFunc("/getimage", cSrv.getImage)

	go func() {
		err := cSrv.Server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			// silent close
		} else if err != nil {
			log.Printf("error starting server: %s\n", err)
		}
	}()

	return cSrv, nil
}

// getImage writes the reponse to `w` based on the incoming request `r`.
func (c *CustomServer) getImage(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}

	// Use http.ServeFile to efficiently stream the file
	http.ServeFile(w, r, fileName)
}
