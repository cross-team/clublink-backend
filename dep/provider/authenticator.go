package provider

import (
	"time"

	"github.com/cross-team/clublink-backend/app/usecase/authenticator"
	"github.com/short-d/app/fw/crypto"
	"github.com/short-d/app/fw/timer"
)

// TokenValidDuration represents the duration of a valid token.
type TokenValidDuration time.Duration

// NewAuthenticator creates Authenticator with TokenValidDuration to uniquely identify duration during dependency injection.
func NewAuthenticator(tokenizer crypto.Tokenizer, timer timer.Timer, duration TokenValidDuration) authenticator.Authenticator {
	return authenticator.NewAuthenticator(tokenizer, timer, time.Duration(duration))
}
