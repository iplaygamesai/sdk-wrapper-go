package iplaygames

import "errors"

var (
	// ErrAPIKeyRequired is returned when no API key is provided
	ErrAPIKeyRequired = errors.New("api key is required")

	// ErrWebhookSecretRequired is returned when webhook secret is not configured
	ErrWebhookSecretRequired = errors.New("webhook secret not configured")
)
