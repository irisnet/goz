package nets

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func GetHTTPClient(timeout time.Duration) *http.Client {
	if timeout != 0 {
		httpClient.Timeout = timeout
	}
	return httpClient
}
