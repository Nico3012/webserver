package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

type FileCache string

func (fileCache FileCache) Get(ctx context.Context, key string) ([]byte, error) {
	str, err := Get(string(fileCache), key)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, autocert.ErrCacheMiss
		}

		return nil, err
	}

	return []byte(str), nil
}

func (fileCache FileCache) Put(ctx context.Context, key string, data []byte) error {
	err := Put(string(fileCache), key, string(data))

	return err
}

func (fileCache FileCache) Delete(ctx context.Context, key string) error {
	err := Delete(string(fileCache), key)

	return err
}

func main() {
	fileCache := FileCache("cache.json")

	manager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("liquipay.de"),
		Client: &acme.Client{
			// DirectoryURL: "https://acme-v02.api.letsencrypt.org/directory",
			DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
		},
		Cache: fileCache,
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
