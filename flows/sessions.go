package flows

import (
	"context"
	"fmt"
	"time"

	apiclient "github.com/iplaygamesai/api-client-go"
)

// SessionsFlow provides high-level operations for sessions
type SessionsFlow struct {
	api *apiclient.APIClient
}

// NewSessionsFlow creates a new sessions flow
func NewSessionsFlow(api *apiclient.APIClient) *SessionsFlow {
	return &SessionsFlow{api: api}
}

// StartSessionParams contains parameters for starting a session
type StartSessionParams struct {
	GameID           int
	PlayerID         string
	Currency         string
	CountryCode      string
	IPAddress        string
	ReturnURL        string
	Locale           string
	Device           string
	Provider         string
	FreespinID       string
	FreespinCount    int
	FreespinBetAmount float64
	ExpireDays       int
}

// SessionResponse represents a session response
type SessionResponse struct {
	Success   bool        `json:"success"`
	SessionID string      `json:"session_id"`
	GameURL   string      `json:"game_url"`
	ExpiresAt string      `json:"expires_at,omitempty"`
	Error     string      `json:"error,omitempty"`
	Raw       interface{} `json:"raw,omitempty"`
}

// Start starts a new game session
func (f *SessionsFlow) Start(ctx context.Context, params StartSessionParams) SessionResponse {
	req := apiclient.NewStartAGameSessionRequest()
	req.SetGameId(int32(params.GameID))
	req.SetPlayerId(params.PlayerID)
	req.SetCurrency(params.Currency)
	req.SetCountryCode(params.CountryCode)
	req.SetIpAddress(params.IPAddress)

	if params.ReturnURL != "" {
		req.SetReturnUrl(params.ReturnURL)
	}
	if params.Locale != "" {
		req.SetLocale(params.Locale)
	}
	if params.Device != "" {
		req.SetDevice(params.Device)
	}
	if params.Provider != "" {
		req.SetProvider(params.Provider)
	}
	if params.FreespinID != "" {
		req.SetFreespinId(params.FreespinID)
	}
	if params.FreespinCount > 0 {
		req.SetFreespinCount(int32(params.FreespinCount))
	}
	if params.FreespinBetAmount > 0 {
		req.SetFreespinBetAmount(float32(params.FreespinBetAmount))
	}
	if params.ExpireDays > 0 {
		req.SetExpireDays(int32(params.ExpireDays))
	}

	resp, _, err := f.api.GameSessionsAPI.StartAGameSession(ctx).StartAGameSessionRequest(*req).Execute()
	if err != nil {
		return SessionResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return SessionResponse{
		Success:   true,
		SessionID: resp.Data.GetSessionId(),
		GameURL:   resp.Data.GetGameUrl(),
		ExpiresAt: resp.Data.GetExpiresAt(),
		Raw:       resp,
	}
}

// Status gets session status
func (f *SessionsFlow) Status(ctx context.Context, sessionID string) ApiResponse {
	resp, _, err := f.api.GameSessionsAPI.GetSessionStatus(ctx, sessionID).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"session_id": sessionID,
			"status":     resp.Data.GetStatus(),
			"player_id":  resp.Data.GetPlayerId(),
			"game_id":    resp.Data.GetGameId(),
			"game":       resp.Data.GetGame(),
			"created_at": resp.Data.GetCreatedAt(),
		},
		Raw: resp,
	}
}

// End ends a game session
func (f *SessionsFlow) End(ctx context.Context, sessionID string) ApiResponse {
	_, err := f.api.GameSessionsAPI.EndAGameSession(ctx, sessionID).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"session_id": sessionID,
			"message":    "Session ended successfully",
		},
	}
}

// StartDemo starts a demo session
func (f *SessionsFlow) StartDemo(ctx context.Context, gameID int, params StartSessionParams) SessionResponse {
	if params.PlayerID == "" {
		params.PlayerID = fmt.Sprintf("demo_%d", time.Now().UnixMilli())
	}
	if params.Currency == "" {
		params.Currency = "USD"
	}
	if params.CountryCode == "" {
		params.CountryCode = "US"
	}
	if params.IPAddress == "" {
		params.IPAddress = "127.0.0.1"
	}
	params.GameID = gameID

	return f.Start(ctx, params)
}
