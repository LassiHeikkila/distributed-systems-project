package accountclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var addr string

const (
	userIDKey = "userID"
)

func SetAccountServiceAddr(a string) {
	addr = a
}

func ValidateUserToken(t string) (string, error) {
	p, _ := url.JoinPath(addr, "/internal/token/validate/", t)
	resp, err := http.Get(url.PathEscape(p))
	if err != nil {
		return "", fmt.Errorf("error performing GET request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("response was not 200 OK: %d", resp.StatusCode)
	}

	d := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var m map[string]any
	err = d.Decode(&m)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	userID, ok := m[userIDKey]
	if !ok {
		return "", fmt.Errorf("response did not contain user ID")
	}

	return userID.(string), nil
}
