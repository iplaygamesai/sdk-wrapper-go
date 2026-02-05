// Package iplaygames provides a high-level SDK for the IPlayGames API
package iplaygames

import (
	apiclient "github.com/iplaygamesai/api-client-go"
	"github.com/iplaygamesai/sdk-wrapper-go/flows"
	"github.com/iplaygamesai/sdk-wrapper-go/webhooks"
)

// ClientOptions contains configuration for the SDK client
type ClientOptions struct {
	APIKey        string
	BaseURL       string
	WebhookSecret string
	Timeout       int
	Debug         bool
}

// Client is the main entry point for the IPlayGames SDK
type Client struct {
	config        *apiclient.Configuration
	apiClient     *apiclient.APIClient
	webhookSecret string
	baseURL       string

	// Lazy-loaded flows
	gamesFlow           *flows.GamesFlow
	sessionsFlow        *flows.SessionsFlow
	multiSessionFlow    *flows.MultiSessionFlow
	jackpotFlow         *flows.JackpotFlow
	promotionsFlow      *flows.PromotionsFlow
	jackpotWidgetFlow   *flows.JackpotWidgetFlow
	promotionWidgetFlow *flows.PromotionWidgetFlow
	webhookHandler      *webhooks.Handler
}

// NewClient creates a new IPlayGames SDK client
func NewClient(opts ClientOptions) (*Client, error) {
	if opts.APIKey == "" {
		return nil, ErrAPIKeyRequired
	}

	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = "https://api.iplaygames.ai"
	}

	config := apiclient.NewConfiguration()
	config.Host = baseURL
	config.AddDefaultHeader("Authorization", "Bearer "+opts.APIKey)

	apiClient := apiclient.NewAPIClient(config)

	return &Client{
		config:        config,
		apiClient:     apiClient,
		webhookSecret: opts.WebhookSecret,
		baseURL:       baseURL,
	}, nil
}

// GetBaseURL returns the base URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// GetAPIClient returns the underlying API client
func (c *Client) GetAPIClient() *apiclient.APIClient {
	return c.apiClient
}

// Games returns the games flow
func (c *Client) Games() *flows.GamesFlow {
	if c.gamesFlow == nil {
		c.gamesFlow = flows.NewGamesFlow(c.apiClient)
	}
	return c.gamesFlow
}

// Sessions returns the sessions flow
func (c *Client) Sessions() *flows.SessionsFlow {
	if c.sessionsFlow == nil {
		c.sessionsFlow = flows.NewSessionsFlow(c.apiClient)
	}
	return c.sessionsFlow
}

// MultiSession returns the multi-session flow
func (c *Client) MultiSession() *flows.MultiSessionFlow {
	if c.multiSessionFlow == nil {
		c.multiSessionFlow = flows.NewMultiSessionFlow(c.apiClient)
	}
	return c.multiSessionFlow
}

// Jackpot returns the jackpot flow
func (c *Client) Jackpot() *flows.JackpotFlow {
	if c.jackpotFlow == nil {
		c.jackpotFlow = flows.NewJackpotFlow(c.apiClient)
	}
	return c.jackpotFlow
}

// Promotions returns the promotions flow
func (c *Client) Promotions() *flows.PromotionsFlow {
	if c.promotionsFlow == nil {
		c.promotionsFlow = flows.NewPromotionsFlow(c.apiClient)
	}
	return c.promotionsFlow
}

// JackpotWidget returns the jackpot widget flow
func (c *Client) JackpotWidget() *flows.JackpotWidgetFlow {
	if c.jackpotWidgetFlow == nil {
		c.jackpotWidgetFlow = flows.NewJackpotWidgetFlow(c.apiClient, c.baseURL)
	}
	return c.jackpotWidgetFlow
}

// PromotionWidget returns the promotion widget flow
func (c *Client) PromotionWidget() *flows.PromotionWidgetFlow {
	if c.promotionWidgetFlow == nil {
		c.promotionWidgetFlow = flows.NewPromotionWidgetFlow(c.apiClient, c.baseURL)
	}
	return c.promotionWidgetFlow
}

// Webhooks returns the webhook handler
func (c *Client) Webhooks() (*webhooks.Handler, error) {
	if c.webhookHandler == nil {
		if c.webhookSecret == "" {
			return nil, ErrWebhookSecretRequired
		}
		c.webhookHandler = webhooks.NewHandler(c.webhookSecret)
	}
	return c.webhookHandler, nil
}

// CreateWebhookHandler creates a webhook handler with a specific secret
func (c *Client) CreateWebhookHandler(secret string) *webhooks.Handler {
	return webhooks.NewHandler(secret)
}
