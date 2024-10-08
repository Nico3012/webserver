package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"test7/cache"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

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

	listener, err := tls.Listen("tcp", ":5000", &tls.Config{
		GetCertificate: manager.GetCertificate,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting servers")

	httpHandler := http.NewServeMux()

	httpHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hallo Welt! This is http"))
	})

	go http.ListenAndServe(":80", manager.HTTPHandler(httpHandler))
	log.Fatal(http.Serve(listener, http.FileServer(http.Dir("html"))))
}
