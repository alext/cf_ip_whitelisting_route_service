package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := ":" + os.Getenv("PORT")

	if os.Getenv("SKIP_SSL_VALIDATION") != "" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	proxy := NewAuthProxy()

	err := http.ListenAndServe(addr, proxy)
	if err != nil {
		log.Fatal(err)
	}
}
