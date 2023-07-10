package net

import (
	"errors"
	"net"
	"net/http"
)

// GetIP returns request real ip.
func GetIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("no valid ip found")
}
