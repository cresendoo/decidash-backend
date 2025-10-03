package utils

import "os"

func Hostname() string {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		return "localhost"
	}
	return hostname
}
