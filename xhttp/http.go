package xhttp

import (
	"net/http"
	"time"
)

func NewXhttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}
