package webserver

import (
	"crypto/tls"
	"net/http"

	"github.com/Nico3012/webserver/jsoncache"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func listenAndServeTLS(addr string, getCertificate func(hello *tls.ClientHelloInfo) (*tls.Certificate, error), handler http.Handler) error {
	listener, err := tls.Listen("tcp", addr, &tls.Config{
		GetCertificate: getCertificate,
	})
	if err != nil {
		return err
	}

	return http.Serve(listener, handler)
}

// exported functions:

// this function creates a web server from certificate files
// the server automatically reacts to certificate file updates
func CreateWebServer(addr string, certPath string, certKeyPath string, handler http.Handler) {
}

// this function creates a web server from certificate authority (ca) files
// the server automatically reacts to certificate authority (ca) file updates
// the server automatically renews its certificates before they expire
// compared to Let’s Encrypt, this function does not need a cache location to store certificates, because self signed certificates can be generated as often as needed
func CreateWebServerWithSelfSignedCertificate(addr string, caPath string, caKeyPath string, hosts []string, handler http.Handler) {
}

// this function creates a web server from Let’s Encrypt certificate authority (ca)
// domain Autheintification is done by http-01 challenge, which requires a http server on port 80.
// the certificates will be renewed automatically
// cacheJsonPath is the file path to a json file, which acts like a key- value storage for certificate information
// compared to self signed, this function needs a cache location to store certificates, because Let’s Encrypt limits the amount of certificates, that can be generated within a few days
// testing boolean indicates if the staging environment is used
func CreateWebServerWithLetsEncryptCertificate(addr string, jsonCachePath string, testing bool, hosts []string, handler http.Handler, httpHandler http.Handler) {
	jsonCache := jsoncache.JsonCache(jsonCachePath)

	directoryURL := "https://acme-v02.api.letsencrypt.org/directory"
	if testing {
		directoryURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
		Client: &acme.Client{
			DirectoryURL: directoryURL,
		},
		Cache: jsonCache,
	}
}
