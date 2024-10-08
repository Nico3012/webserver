package webserver

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/Nico3012/webserver/internal/jsoncache"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func listenAndServe(addr string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(listener, handler)
}

func listenAndServeTLS(addr string, handler http.Handler, getCertificate func(hello *tls.ClientHelloInfo) (*tls.Certificate, error)) error {
	listener, err := tls.Listen("tcp", addr, &tls.Config{
		GetCertificate: getCertificate,
	})
	if err != nil {
		return err
	}

	return http.Serve(listener, handler)
}

// exported functions:

func CreateWebServer(addr string, handler http.Handler) error {
	return listenAndServe(addr, handler)
}

// this function creates a web server from certificate files
// the server automatically reacts to certificate file updates
func CreateSecureWebServer(addr string, certPath string, certKeyPath string, handler http.Handler) {
}

// this function creates a web server from certificate authority (ca) files
// the server automatically reacts to certificate authority (ca) file updates
// the server automatically renews its certificates before they expire
// compared to Let’s Encrypt, this function does not need a cache location to store certificates, because self signed certificates can be generated as often as needed
func CreateSecureWebServerAndSelfSignedCertificate(addr string, caPath string, caKeyPath string, hosts []string, handler http.Handler) {
}

// this function creates a web server from Let’s Encrypt certificate authority (ca)
// domain Autheintification is done by http-01 challenge, which requires a http server on port 80. THIS SERVER WILL NOT BE CREATED!!! Instead a function, that requires a "normal" http handler is returned, that returns the http header with the http-01 challenge included
// the certificates will be renewed automatically
// cacheJsonPath is the file path to a json file, which acts like a key- value storage for certificate information
// compared to self signed, this function needs a cache location to store certificates, because Let’s Encrypt limits the amount of certificates, that can be generated within a few days
// testing boolean indicates if the staging environment is used
func CreateSecureWebServerAndLetsEncryptCertificate(addr string, jsonCachePath string, testing bool, hosts []string, handler http.Handler) func(http.Handler) http.Handler {
	jsonCache := jsoncache.JsonCache(jsonCachePath)

	directoryURL := "https://acme-v02.api.letsencrypt.org/directory"
	if testing {
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

	// go routine
	listenAndServeTLS(addr, handler, manager.GetCertificate)

	return manager.HTTPHandler
}
