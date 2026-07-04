package application

import (
	"crypto/sha256"
	"sync"
	"time"
)

// TokenStore is an in-memory store for revoked refresh tokens.
// Tokens are stored by their SHA-256 hash and auto-expire after TTL.
type TokenStore struct {
	mu      sync.RWMutex
	revoked map[string]time.Time
	ttl     time.Duration
}

func NewTokenStore(ttl time.Duration) *TokenStore {
	return &TokenStore{
		revoked: make(map[string]time.Time),
		ttl:     ttl,
	}
}

func (s *TokenStore) Revoke(token string) {
	h := hashToken(token)
	s.mu.Lock()
	s.revoked[h] = time.Now().Add(s.ttl)
	s.mu.Unlock()
}

func (s *TokenStore) IsRevoked(token string) bool {
	h := hashToken(token)
	s.mu.RLock()
	expiry, ok := s.revoked[h]
	s.mu.RUnlock()
	if !ok {
		return false
	}
	if time.Now().After(expiry) {
		s.mu.Lock()
		delete(s.revoked, h)
		s.mu.Unlock()
		return false
	}
	return true
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return string(h[:])
}
