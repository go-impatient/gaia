package request

import (
	"net/http"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: time.Second * 3,
	}
	userAgent string
)

// SetHttpClient ...
func SetHttpClient(client *http.Client) {
	httpClient = client
}

// SetUserAgent ...
func SetUserAgent(ua string) {
	userAgent = ua
}
