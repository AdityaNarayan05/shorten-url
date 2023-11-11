package helpers

import (
	"os"
	"strings"
)

// EnforceHTTP adds the "http://" prefix to a URL if it doesn't already start with it.
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	} else {
		return url
	}
}

// RemoveDomainError checks if a URL contains the specified domain and removes common variations.
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}

	return true
}
