package httputils

import "strings"

func ParseAllowedOrigins(s string) []string {
	return strings.Split(s, ",")
}
