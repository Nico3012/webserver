package secure

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/Nico3012/webserver/internal/jsoncache"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func listenAndServeTLS(addr string, handler http.Handler, getCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)) error {
	listener, err := tls.Listen("tcp", addr, &tls.Config{
		GetCertificate: getCertificate,
	})

	if err != nil {
		return err
	}

	return http.Serve(listener, handler)
}

func CreateWebServer(addr string, handler http.Handler, certFile string, certKeyFile string) error {
	return listenAndServeTLS(addr, handler, func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cert, err := tls.LoadX509KeyPair(certFile, certKeyFile)

		if err != nil {
			return &tls.Certificate{}, err
		}

		return &cert, nil
	})
}

func CreateWebServerWithLetsEncryptCertificate(addr string, handler http.Handler, jsonCacheFile string, hosts []string, staging bool) (func(http.Handler) http.Handler, error) {
	jsonCache := jsoncache.JsonCache(jsonCacheFile)

	directoryURL := "https://acme-v02.api.letsencrypt.org/directory"
	if staging {
		directoryURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	client := &acme.Client{
		DirectoryURL: directoryURL,
	}

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Client:     client,
		Cache:      jsonCache,
	}

	go func() {
		err := listenAndServeTLS(addr, handler, manager.GetCertificate)

		log.Fatal(err)
	}()

	return manager.HTTPHandler, nil
}
