package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	CF_FORWARDED_URL_HEADER = "X-CF-Forwarded-Url"
)

type AuthProxy struct {
	ipset          *IPSet
	trustedRouters map[string]bool
	backend        http.Handler
}

func NewAuthProxy(ipset *IPSet, trustedRouters []string) *AuthProxy {
	a := &AuthProxy{
		ipset:   ipset,
		backend: buildBackendProxy(),
	}
	a.setTrustedRouters(trustedRouters)
	return a
}

func (a *AuthProxy) setTrustedRouters(routers []string) {
	a.trustedRouters = make(map[string]bool)
	for _, r := range routers {
		a.trustedRouters[r] = true
	}
}

func buildBackendProxy() http.Handler {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			forwardedURL := req.Header.Get(CF_FORWARDED_URL_HEADER)
			if forwardedURL == "" {
				// This should never happen due to the check in AuthProxy.ServeHTTP
				panic("missing forwarded URL")
			}
			url, err := url.Parse(forwardedURL)
			if err != nil {
				// This should never happen due to the check in AuthProxy.ServeHTTP
				panic("Invalid forwarded URL: " + err.Error())
			}

			req.URL = url
			req.Host = url.Host
		},
	}
}

func (a *AuthProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqIP := a.extractRequestIP(req)
	if !a.ipset.Contains(reqIP) {
		log.Printf("Denied access for IP '%s' (X-Forwarded-For: %s)", reqIP, req.Header.Get("X-Forwarded-For"))
		http.Error(w, "Permission denied.", 403)
		return
	}

	forwardedURL := req.Header.Get(CF_FORWARDED_URL_HEADER)
	if forwardedURL == "" {
		http.Error(w, "Missing Forwarded URL", http.StatusBadRequest)
		return
	}
	_, err := url.Parse(forwardedURL)
	if err != nil {
		http.Error(w, "Invalid forward URL: "+err.Error(), http.StatusBadRequest)
		return
	}

	a.backend.ServeHTTP(w, req)
}

func (a *AuthProxy) extractRequestIP(req *http.Request) string {
	xff := req.Header.Get("X-Forwarded-For")
	entriesWithBlanks := strings.Split(xff, ",")

	entries := entriesWithBlanks[:0]
	for _, ip := range entriesWithBlanks {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			entries = append(entries, ip)
		}
	}

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		if !a.trustedRouters[entry] {
			return entry
		}
	}
	return ""
}
