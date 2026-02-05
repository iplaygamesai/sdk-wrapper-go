// Package tests provides integration tests for the IPlayGames Go SDK
//
// Usage: go test -v ./tests/
package tests

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	iplaygames "github.com/iplaygamesai/sdk-wrapper-go"
	"github.com/iplaygamesai/sdk-wrapper-go/flows"
	"github.com/iplaygamesai/sdk-wrapper-go/webhooks"
)

var (
	apiKey        = getEnvOrDefault("IPLAYGAMES_API_KEY", "YOUR_API_TOKEN")
	baseURL       = getEnvOrDefault("IPLAYGAMES_BASE_URL", "https://gamehub.test")
	webhookSecret = "test_secret_for_webhooks"
)

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func TestClientInitialization(t *testing.T) {
	client, err := iplaygames.NewClient(iplaygames.ClientOptions{
		APIKey:        apiKey,
		BaseURL:       baseURL,
		WebhookSecret: webhookSecret,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client.GetBaseURL() != baseURL {
		t.Errorf("Expected base URL %s, got %s", baseURL, client.GetBaseURL())
	}
}

func TestClientRequiresAPIKey(t *testing.T) {
	_, err := iplaygames.NewClient(iplaygames.ClientOptions{
		BaseURL: baseURL,
	})
	if err == nil {
		t.Error("Expected error when API key is missing")
	}
}

func TestWebhookVerifyValidSignature(t *testing.T) {
	handler := webhooks.NewHandler(webhookSecret)

	payload := `{"type":"bet","player_id":"player_456","currency":"USD","amount":1000}`

	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	if !handler.Verify(payload, signature) {
		t.Error("Valid signature should verify")
	}
}

func TestWebhookRejectInvalidSignature(t *testing.T) {
	handler := webhooks.NewHandler(webhookSecret)

	payload := `{"type":"bet"}`

	if handler.Verify(payload, "invalid_signature") {
		t.Error("Invalid signature should be rejected")
	}
}

func TestWebhookParsePayload(t *testing.T) {
	handler := webhooks.NewHandler(webhookSecret)

	payload := `{"type":"bet","player_id":"player_456","currency":"USD","amount":1000,"transaction_id":12345}`

	parsed, err := handler.Parse(payload)
	if err != nil {
		t.Fatalf("Failed to parse payload: %v", err)
	}

	if parsed.Type != "bet" {
		t.Errorf("Expected type 'bet', got '%s'", parsed.Type)
	}
	if parsed.PlayerID != "player_456" {
		t.Errorf("Expected player_id 'player_456', got '%s'", parsed.PlayerID)
	}
	if parsed.Amount == nil || *parsed.Amount != 1000 {
		t.Error("Expected amount 1000")
	}
}

func TestWebhookSuccessResponse(t *testing.T) {
	handler := webhooks.NewHandler(webhookSecret)

	response := handler.SuccessResponse(100.50, nil)

	if response["status"] != "success" {
		t.Error("Status should be success")
	}
	if response["balance"] != 10050 {
		t.Errorf("Balance should be 10050 cents, got %v", response["balance"])
	}
}

func TestWebhookErrorResponse(t *testing.T) {
	handler := webhooks.NewHandler(webhookSecret)

	response := handler.ErrorResponse("TEST_ERROR", "Test message")

	if response["status"] != "error" {
		t.Error("Status should be error")
	}
	if response["error_code"] != "TEST_ERROR" {
		t.Error("Error code should match")
	}
}

func TestWebhookPlayerNotFoundResponse(t *testing.T) {
	handler := webhooks.NewHandler(webhookSecret)

	response := handler.PlayerNotFoundResponse()

	if response["error_code"] != "PLAYER_NOT_FOUND" {
		t.Error("Should be player not found")
	}
}

func TestWebhookInsufficientFundsResponse(t *testing.T) {
	handler := webhooks.NewHandler(webhookSecret)

	response := handler.InsufficientFundsResponse(50.25)

	if response["error_code"] != "INSUFFICIENT_FUNDS" {
		t.Error("Should be insufficient funds")
	}
	if response["balance"] != 5025 {
		t.Errorf("Balance should be 5025 cents, got %v", response["balance"])
	}
}

func TestWebhookPayloadHelpers(t *testing.T) {
	handler := webhooks.NewHandler(webhookSecret)

	// Test bet payload
	betPayload := `{"type":"bet","player_id":"player_456","amount":1000}`
	bet, _ := handler.Parse(betPayload)
	if !bet.IsBet() {
		t.Error("Should be a bet")
	}
	if bet.IsWin() {
		t.Error("Should not be a win")
	}

	// Test win payload
	winPayload := `{"type":"win","player_id":"player_456","amount":2000}`
	win, _ := handler.Parse(winPayload)
	if !win.IsWin() {
		t.Error("Should be a win")
	}

	// Test amount conversion
	dollars := bet.GetAmountInDollars()
	if dollars == nil || *dollars != 10.0 {
		t.Errorf("Expected 10.0 dollars, got %v", dollars)
	}
}

func TestJackpotWidgetEmbedCode(t *testing.T) {
	client, _ := iplaygames.NewClient(iplaygames.ClientOptions{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})

	embedCode := client.JackpotWidget().GetEmbedCode("test_token", flows.EmbedOptions{
		Theme:     "dark",
		Container: "my-widget",
	})

	if !strings.Contains(embedCode, "test_token") {
		t.Error("Embed code should contain token")
	}
	if !strings.Contains(embedCode, "my-widget") {
		t.Error("Embed code should contain container ID")
	}
	if !strings.Contains(embedCode, "jackpot.js") {
		t.Error("Embed code should reference jackpot.js")
	}
}

func TestMultiSessionIframeGeneration(t *testing.T) {
	client, _ := iplaygames.NewClient(iplaygames.ClientOptions{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})

	iframe := client.MultiSession().GetIframe("https://example.com/swipe", flows.IframeOptions{
		Width:  "100%",
		Height: "600px",
		ID:     "game-swiper",
	})

	if !strings.Contains(iframe, "https://example.com/swipe") {
		t.Error("Iframe should contain URL")
	}
	if !strings.Contains(iframe, `id="game-swiper"`) {
		t.Error("Iframe should have ID")
	}
	if !strings.Contains(iframe, "allowfullscreen") {
		t.Error("Iframe should allow fullscreen")
	}
}

// Integration tests - require IPLAYGAMES_API_KEY and IPLAYGAMES_BASE_URL env vars
// These tests will be skipped if the API key is not set or is the default value

func skipIfNoAPIKey(t *testing.T) {
	if apiKey == "YOUR_API_TOKEN" || apiKey == "" {
		t.Skip("Skipping integration test: IPLAYGAMES_API_KEY not set")
	}
}

func TestGamesListIntegration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := iplaygames.NewClient(iplaygames.ClientOptions{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	response := client.Games().List(ctx, flows.ListParams{
		PerPage: 3,
	})

	if !response.Success {
		t.Fatalf("Games list failed: %s", response.Error)
	}

	t.Logf("Found %d games (Total: %d)", len(response.Games), response.Meta.Total)
	for _, game := range response.Games {
		t.Logf("  - %s (%s)", game.Title, game.Producer)
	}
}

func TestJackpotConfigurationIntegration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := iplaygames.NewClient(iplaygames.ClientOptions{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	response := client.Jackpot().GetConfiguration(ctx)

	if !response.Success {
		t.Fatalf("Jackpot configuration failed: %s", response.Error)
	}

	t.Logf("Jackpot configuration: %v", response.Data)
}

func TestJackpotPoolsIntegration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := iplaygames.NewClient(iplaygames.ClientOptions{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	response := client.Jackpot().GetPools(ctx)

	if !response.Success {
		t.Fatalf("Jackpot pools failed: %s", response.Error)
	}

	t.Logf("Jackpot pools: %v", response.Data)
}

func TestPromotionsListIntegration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := iplaygames.NewClient(iplaygames.ClientOptions{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	response := client.Promotions().List(ctx, "", "")

	if !response.Success {
		t.Fatalf("Promotions list failed: %s", response.Error)
	}

	t.Logf("Promotions: %v", response.Data)
}
