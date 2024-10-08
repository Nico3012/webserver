package unsecure

import (
	"net"
	"net/http"
)

func CreateWebServer(addr string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	return http.Serve(listener, handler)
}
