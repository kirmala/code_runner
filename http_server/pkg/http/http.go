package httpx

import (
	"net/http"
)

func CreateAndRunServer(h http.Handler, addr string) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: h,
	}

	return httpServer.ListenAndServe()
}
