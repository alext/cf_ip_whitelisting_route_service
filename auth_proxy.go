package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	CF_FORWARDED_URL_HEADER = "X-CF-Forwarded-Url"
)

type AuthProxy struct {
	backend http.Handler
}

func NewAuthProxy() http.Handler {
	return &AuthProxy{
		backend: buildBackendProxy(),
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
	//TODO actual auth.
	if false {
		w.Header().Set("WWW-Authenticate", `Basic realm="auth"`)
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
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
