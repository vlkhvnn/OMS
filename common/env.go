package common

import "syscall"

func GetString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}
	return fallback
}
