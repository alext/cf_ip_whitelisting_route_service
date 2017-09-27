package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	addr := ":" + os.Getenv("PORT")

	if os.Getenv("SKIP_SSL_VALIDATION") != "" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	ipset, err := NewIPSet(strings.Split(os.Getenv("WHITELIST_ADDRS"), ","))
	if err != nil {
		log.Fatal(err)
	}

	xffOffset, err := parseXFFOffset(os.Getenv("XFF_OFFSET"))

	proxy := NewAuthProxy(ipset, xffOffset)

	err = http.ListenAndServe(addr, proxy)
	if err != nil {
		log.Fatal(err)
	}
}

func parseXFFOffset(input string) (int, error) {
	if input == "" {
		return 0, nil
	}
	return strconv.Atoi(input)
}
