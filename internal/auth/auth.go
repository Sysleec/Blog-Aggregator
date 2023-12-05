package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey extracts an API key from
// the headers of an HTTP request
// Example:
// Authorization: ApiKey {insert apikey here}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication into found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("incorrect auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("incorrect first path of auth header")
	}
	return vals[1], nil
}
