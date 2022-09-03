package main

import (
	"net"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}

	// http.DefaultTransport
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			ForceAttemptHTTP2: true,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       10 * time.Second,
			ExpectContinueTimeout: 90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 1 * time.Second,
		},
	}
	_ = client

}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))

}
