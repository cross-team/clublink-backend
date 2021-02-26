package provider

import (
	"github.com/cross-team/clublink-backend/app/usecase/requester"
	"github.com/short-d/app/fw/env"
)

// NewVerifier creates Verifier based on
// server environment.
func NewVerifier(
	deployment env.Deployment,
	service requester.ReCaptcha,
) requester.Verifier {
	if deployment.IsDevelopment() {
		return requester.NewVerifierFake()
	}
	return requester.NewReCaptchaVerifier(service)
}
