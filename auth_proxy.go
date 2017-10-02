package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	CF_FORWARDED_URL_HEADER = "X-CF-Forwarded-Url"
)

type AuthProxy struct {
	ipset     *IPSet
	xffOffset int
	backend   http.Handler
}

func NewAuthProxy(ipset *IPSet, xffOffset int) *AuthProxy {
	return &AuthProxy{
		ipset:     ipset,
		xffOffset: xffOffset,
		backend:   buildBackendProxy(),
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
	reqIP := extractRequestIP(req, a.xffOffset)
	if !a.ipset.Contains(reqIP) {
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

func extractRequestIP(req *http.Request, xffOffset int) string {
	xff := req.Header.Get("X-Forwarded-For")
	entriesWithBlanks := strings.Split(xff, " ")

	entries := entriesWithBlanks[:0]
	for _, ip := range entriesWithBlanks {
		if ip != "" {
			entries = append(entries, ip)
		}
	}

	if len(entries) < xffOffset+1 {
		return ""
	}

	return entries[len(entries)-1-xffOffset]
}
