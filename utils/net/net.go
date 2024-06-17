package net

import (
	"net"
	"time"
)

func IsHostReachable(url string) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
