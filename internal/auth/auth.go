package auth

import (
	"errors"
	"net/http"
	"strings"
)

//Authorization: ApiKey {key}
func GetApiKey(headers http.Header) (string, error){
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no auth found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed api key")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed api key pattern")
	}

	return vals[1], nil
}