package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"test7/cache"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func startHttp(addr string, handler http.Handler) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal(err)
	}
}

func startHttps(addr string, handler http.Handler, getCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)) {
	listener, err := tls.Listen("tcp", addr, &tls.Config{
		GetCertificate: getCertificate,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cache := cache.Cache("./cache.json")

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("liquipay.de"),
		Client: &acme.Client{
			// DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory",
			DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
		},
		Cache: cache,
	}

	log.Println("Starting servers")

	go startHttp(":80", manager.HTTPHandler(http.RedirectHandler("https://liquipay.de:5000", http.StatusTemporaryRedirect)))
	startHttps(":5000", http.FileServer(http.Dir("html")), manager.GetCertificate)
}
