package sAuth

import (
	"sync"
	"time"
)

type TokenBlacklist struct {
	tokens map[string]time.Time
	mutex  sync.RWMutex
}

func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		tokens: make(map[string]time.Time),
	}
}

func (tb *TokenBlacklist) Add(token string, expiresAt time.Time) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.tokens[token] = expiresAt
}

func (tb *TokenBlacklist) IsRevoked(token string) bool {
	tb.mutex.RLock()
	defer tb.mutex.RUnlock()

	expiry, exists := tb.tokens[token]
	if !exists {
		return false
	}

	// Limpiar si está expirado
	if time.Now().After(expiry) {
		tb.mutex.Lock()
		delete(tb.tokens, token)
		tb.mutex.Unlock()
		return false
	}

	return true
}

func (tb *TokenBlacklist) Cleanup() {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	for token, expiry := range tb.tokens {
		if now.After(expiry) {
			delete(tb.tokens, token)
		}
	}
}
