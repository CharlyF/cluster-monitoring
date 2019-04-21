package util

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	defaultExpire = 5 * time.Minute
	defaultPurge  = 30 * time.Second

	// encapsulate the cache module for easy refactoring

	// NoExpiration maps to go-cache corresponding value
	NoExpiration = cache.NoExpiration
)

// Cache provides an in-memory key:value store similar to memcached
var Cache = cache.New(defaultExpire, defaultPurge)
