package flows

import (
	"context"
	"fmt"

	apiclient "github.com/iplaygamesai/api-client-go"
)

// MultiSessionFlow provides high-level operations for multi-sessions
type MultiSessionFlow struct {
	api *apiclient.APIClient
}

// NewMultiSessionFlow creates a new multi-session flow
func NewMultiSessionFlow(api *apiclient.APIClient) *MultiSessionFlow {
	return &MultiSessionFlow{api: api}
}

// StartMultiSessionParams contains parameters for starting a multi-session
type StartMultiSessionParams struct {
	PlayerID    string
	Currency    string
	CountryCode string
	IPAddress   string
	GameIDs     []string
	Locale      string
	Device      string
}

// MultiSessionGame represents a game in a multi-session
type MultiSessionGame struct {
	Position  int    `json:"position"`
	GameName  string `json:"game_name"`
	GameImage string `json:"game_image,omitempty"`
}

// MultiSessionResponse represents a multi-session response
type MultiSessionResponse struct {
	Success        bool               `json:"success"`
	MultiSessionID string             `json:"multi_session_id"`
	SwipeURL       string             `json:"swipe_url"`
	TotalGames     int                `json:"total_games"`
	Games          []MultiSessionGame `json:"games"`
	ExpiresAt      string             `json:"expires_at,omitempty"`
	Error          string             `json:"error,omitempty"`
	Raw            interface{}        `json:"raw,omitempty"`
}

// Start starts a multi-session for a player
func (f *MultiSessionFlow) Start(ctx context.Context, params StartMultiSessionParams) MultiSessionResponse {
	req := apiclient.NewStartAMultiSessionRequest()
	req.SetPlayerId(params.PlayerID)
	req.SetCurrency(params.Currency)
	req.SetCountryCode(params.CountryCode)
	req.SetIpAddress(params.IPAddress)

	if len(params.GameIDs) > 0 {
		req.SetGameIds(params.GameIDs)
	}
	if params.Locale != "" {
		req.SetLocale(params.Locale)
	}
	if params.Device != "" {
		req.SetDevice(params.Device)
	}

	resp, _, err := f.api.MultiSessionsAPI.StartAMultiSession(ctx).StartAMultiSessionRequest(*req).Execute()
	if err != nil {
		return MultiSessionResponse{
			Success: false,
			Error:   err.Error(),
			Games:   []MultiSessionGame{},
		}
	}

	games := make([]MultiSessionGame, 0)
	if resp.Data != nil && resp.Data.Games != nil {
		for i, g := range resp.Data.Games {
			games = append(games, MultiSessionGame{
				Position:  int(g.GetPosition()),
				GameName:  g.GetGameName(),
				GameImage: g.GetGameImage(),
			})
			if games[i].Position == 0 {
				games[i].Position = i
			}
		}
	}

	return MultiSessionResponse{
		Success:        true,
		MultiSessionID: resp.Data.GetMultiSessionId(),
		SwipeURL:       resp.Data.GetSwipeUrl(),
		TotalGames:     int(resp.Data.GetTotalGames()),
		Games:          games,
		ExpiresAt:      resp.Data.GetExpiresAt(),
		Raw:            resp,
	}
}

// StartWithGames starts a multi-session with specific games
func (f *MultiSessionFlow) StartWithGames(ctx context.Context, gameIDs []string, params StartMultiSessionParams) MultiSessionResponse {
	params.GameIDs = gameIDs
	return f.Start(ctx, params)
}

// StartRandom starts a multi-session with random games
func (f *MultiSessionFlow) StartRandom(ctx context.Context, params StartMultiSessionParams) MultiSessionResponse {
	params.GameIDs = nil
	return f.Start(ctx, params)
}

// Status gets multi-session status
func (f *MultiSessionFlow) Status(ctx context.Context, token string) ApiResponse {
	resp, _, err := f.api.MultiSessionsAPI.GetMultiSessionStatus(ctx, token).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	games := make([]MultiSessionGame, 0)
	if resp.Data != nil && resp.Data.Games != nil {
		for i, g := range resp.Data.Games {
			games = append(games, MultiSessionGame{
				Position:  int(g.GetPosition()),
				GameName:  g.GetGameName(),
				GameImage: g.GetGameImage(),
			})
			if games[i].Position == 0 {
				games[i].Position = i
			}
		}
	}

	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"token":           token,
			"status":          resp.Data.GetStatus(),
			"total_games":     resp.Data.GetTotalGames(),
			"active_sessions": resp.Data.GetActiveSessions(),
			"current_index":   resp.Data.GetCurrentIndex(),
			"games":           games,
			"expires_at":      resp.Data.GetExpiresAt(),
		},
		Raw: resp,
	}
}

// End ends a multi-session
func (f *MultiSessionFlow) End(ctx context.Context, token string) ApiResponse {
	_, err := f.api.MultiSessionsAPI.EndMultiSession(ctx, token).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"token":   token,
			"message": "Multi-session ended successfully",
		},
	}
}

// IframeOptions contains options for generating iframe HTML
type IframeOptions struct {
	Width     string
	Height    string
	ID        string
	ClassName string
	Allow     string
}

// GetIframe generates an iframe HTML element
func (f *MultiSessionFlow) GetIframe(swipeURL string, opts IframeOptions) string {
	if opts.Width == "" {
		opts.Width = "100%"
	}
	if opts.Height == "" {
		opts.Height = "100%"
	}
	if opts.Allow == "" {
		opts.Allow = "fullscreen; autoplay; encrypted-media"
	}

	idAttr := ""
	if opts.ID != "" {
		idAttr = fmt.Sprintf(` id="%s"`, opts.ID)
	}
	classAttr := ""
	if opts.ClassName != "" {
		classAttr = fmt.Sprintf(` class="%s"`, opts.ClassName)
	}

	return fmt.Sprintf(`<iframe%s%s
    src="%s"
    width="%s"
    height="%s"
    frameborder="0"
    allow="%s"
    allowfullscreen>
</iframe>`, idAttr, classAttr, swipeURL, opts.Width, opts.Height, opts.Allow)
}
