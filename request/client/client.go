package client

import (
	"crypto/tls"
	"net/http"
	"time"
)

func NewClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
