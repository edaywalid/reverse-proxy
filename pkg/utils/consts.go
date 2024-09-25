package utils

import "time"

const (
	CacheExpiration = 5 * time.Minute
	CacheSize       = 100
	CleanupInterval = 10 * time.Second
)
