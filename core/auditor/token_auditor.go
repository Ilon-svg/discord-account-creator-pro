// Package auditor handles real-time token health checks post-registration.
package auditor

import (
	"errors"
	"net/http"
	"time"
)

var (
	ErrTokenLocked     = errors.New("discord_token: account flagged or phone verification required")
	ErrTokenUnauthorized = errors.New("discord_token: invalid authentication session")
)

type TokenAuditor struct {
	Client  *http.Client
	Timeout time.Duration
}

func NewTokenAuditor() *TokenAuditor {
	return &TokenAuditor{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (ta *TokenAuditor) AuditSession(token string) (bool, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ta.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusUnauthorized:
		return false, ErrTokenUnauthorized
	case http.StatusForbidden:
		return false, ErrTokenLocked
	default:
		return false, errors.New("discord_token: unexpected verification gateway response")
	}
}

