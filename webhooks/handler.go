// Package webhooks provides webhook handling functionality
package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// Webhook type constants
const (
	TypeAuthenticate = "authenticate"
	TypeBalanceCheck = "balance_check"
	TypeBet          = "bet"
	TypeWin          = "win"
	TypeRollback     = "rollback"
	TypeReward       = "reward"
)

// Payload represents a parsed webhook payload
type Payload struct {
	Type      string `json:"type"`
	PlayerID  string `json:"player_id"`
	Currency  string `json:"currency"`
	Timestamp string `json:"timestamp"`

	// Game fields
	GameID   *int   `json:"game_id,omitempty"`
	GameType string `json:"game_type,omitempty"`

	// Transaction fields
	TransactionID *int   `json:"transaction_id,omitempty"`
	Amount        *int   `json:"amount,omitempty"` // In cents
	SessionID     string `json:"session_id,omitempty"`
	RoundID       string `json:"round_id,omitempty"`

	// Reward fields
	RewardType  string `json:"reward_type,omitempty"`
	RewardTitle string `json:"reward_title,omitempty"`

	// Freespin fields
	IsFreespin           bool     `json:"is_freespin"`
	FreespinID           string   `json:"freespin_id,omitempty"`
	FreespinTotal        *int     `json:"freespin_total,omitempty"`
	FreespinsRemaining   *int     `json:"freespins_remaining,omitempty"`
	FreespinRoundNumber  *int     `json:"freespin_round_number,omitempty"`
	FreespinTotalWinnings *float64 `json:"freespin_total_winnings,omitempty"`

	// Raw data
	Raw map[string]interface{} `json:"-"`
}

// IsBet checks if this is a bet transaction
func (p *Payload) IsBet() bool {
	return p.Type == TypeBet
}

// IsWin checks if this is a win transaction
func (p *Payload) IsWin() bool {
	return p.Type == TypeWin
}

// IsRollback checks if this is a rollback transaction
func (p *Payload) IsRollback() bool {
	return p.Type == TypeRollback
}

// IsReward checks if this is a reward
func (p *Payload) IsReward() bool {
	return p.Type == TypeReward
}

// IsAuthenticate checks if this is an authentication request
func (p *Payload) IsAuthenticate() bool {
	return p.Type == TypeAuthenticate
}

// IsBalanceCheck checks if this is a balance check
func (p *Payload) IsBalanceCheck() bool {
	return p.Type == TypeBalanceCheck
}

// GetAmountInDollars gets amount in dollars (converts from cents)
func (p *Payload) GetAmountInDollars() *float64 {
	if p.Amount == nil {
		return nil
	}
	dollars := float64(*p.Amount) / 100
	return &dollars
}

// Get gets a value from the raw data
func (p *Payload) Get(key string) interface{} {
	return p.Raw[key]
}

// Handler handles webhook verification and parsing
type Handler struct {
	secret string
}

// NewHandler creates a new webhook handler
func NewHandler(secret string) *Handler {
	return &Handler{secret: secret}
}

// Verify verifies webhook signature
func (h *Handler) Verify(payload, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// Parse parses webhook payload
func (h *Handler) Parse(payload string) (*Payload, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		return nil, errors.New("invalid JSON payload")
	}

	p := &Payload{
		Type:      getString(raw, "type"),
		PlayerID:  getString(raw, "player_id"),
		Currency:  getString(raw, "currency"),
		Timestamp: getString(raw, "timestamp"),
		GameType:  getString(raw, "game_type"),
		SessionID: getString(raw, "session_id"),
		RoundID:   getString(raw, "round_id"),
		RewardType:  getString(raw, "reward_type"),
		RewardTitle: getString(raw, "reward_title"),
		Raw:       raw,
	}

	if v, ok := raw["game_id"].(float64); ok {
		i := int(v)
		p.GameID = &i
	}
	if v, ok := raw["transaction_id"].(float64); ok {
		i := int(v)
		p.TransactionID = &i
	}
	if v, ok := raw["amount"].(float64); ok {
		i := int(v)
		p.Amount = &i
	}

	// Freespin fields
	if v, ok := raw["is_freespin_round"].(bool); ok {
		p.IsFreespin = v
	} else if v, ok := raw["is_freespin"].(bool); ok {
		p.IsFreespin = v
	}

	if v := getString(raw, "freespin_id"); v != "" {
		p.FreespinID = v
	} else if v := getString(raw, "bonus_id"); v != "" {
		p.FreespinID = v
	}

	if v, ok := raw["freespin_total"].(float64); ok {
		i := int(v)
		p.FreespinTotal = &i
	}
	if v, ok := raw["freespins_remaining"].(float64); ok {
		i := int(v)
		p.FreespinsRemaining = &i
	} else if v, ok := raw["freespin_left"].(float64); ok {
		i := int(v)
		p.FreespinsRemaining = &i
	}
	if v, ok := raw["freespin_round_number"].(float64); ok {
		i := int(v)
		p.FreespinRoundNumber = &i
	}
	if v, ok := raw["freespin_total_winnings"].(float64); ok {
		p.FreespinTotalWinnings = &v
	}

	return p, nil
}

// VerifyAndParse verifies and parses webhook in one step
func (h *Handler) VerifyAndParse(payload, signature string) (*Payload, error) {
	if !h.Verify(payload, signature) {
		return nil, errors.New("invalid webhook signature")
	}
	return h.Parse(payload)
}

// SuccessResponse creates a success response
func (h *Handler) SuccessResponse(balance float64, extra map[string]interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"status":  "success",
		"balance": int(balance * 100),
	}
	for k, v := range extra {
		resp[k] = v
	}
	return resp
}

// ErrorResponse creates an error response
func (h *Handler) ErrorResponse(code, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":        "error",
		"error_code":    code,
		"error_message": message,
	}
}

// PlayerNotFoundResponse creates a player not found error response
func (h *Handler) PlayerNotFoundResponse() map[string]interface{} {
	return h.ErrorResponse("PLAYER_NOT_FOUND", "Player not found")
}

// InsufficientFundsResponse creates an insufficient funds error response
func (h *Handler) InsufficientFundsResponse(balance float64) map[string]interface{} {
	resp := h.ErrorResponse("INSUFFICIENT_FUNDS", "Insufficient funds")
	resp["balance"] = int(balance * 100)
	return resp
}

// AlreadyProcessedResponse creates a transaction already processed response
func (h *Handler) AlreadyProcessedResponse(balance float64) map[string]interface{} {
	return h.SuccessResponse(balance, map[string]interface{}{"already_processed": true})
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
