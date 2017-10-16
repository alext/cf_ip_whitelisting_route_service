package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	addr := ":" + os.Getenv("PORT")

	if os.Getenv("SKIP_SSL_VALIDATION") != "" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	ipset, err := NewIPSet(parseCommaSeparated(os.Getenv("WHITELIST_ADDRS")))
	if err != nil {
		log.Fatal(err)
	}

	trustedRouters := parseCommaSeparated(os.Getenv("TRUSTED_ROUTERS"))

	proxy := NewAuthProxy(ipset, trustedRouters)

	err = http.ListenAndServe(addr, proxy)
	if err != nil {
		log.Fatal(err)
	}
}

func parseCommaSeparated(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	addrs := strings.Split(input, ",")
	for i, addr := range addrs {
		addrs[i] = strings.TrimSpace(addr)
	}
	return addrs
}
