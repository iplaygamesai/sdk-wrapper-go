package flows

import (
	"context"

	"github.com/iplaygamesai/api-client-go"
)

// GamesFlow provides high-level operations for games
type GamesFlow struct {
	api *apiclient.APIClient
}

// NewGamesFlow creates a new games flow
func NewGamesFlow(api *apiclient.APIClient) *GamesFlow {
	return &GamesFlow{api: api}
}

// ListParams contains parameters for listing games
type ListParams struct {
	Search     string
	ProducerID int
	Provider   string
	Type       string
	PerPage    int
}

// Game represents a game
type Game struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Producer      string `json:"producer"`
	Type          string `json:"type"`
	ImageURL      string `json:"image_url,omitempty"`
	DemoAvailable bool   `json:"demo_available,omitempty"`
}

// GamesListResponse represents a games list response
type GamesListResponse struct {
	Success bool             `json:"success"`
	Games   []Game           `json:"games"`
	Meta    PaginationMeta   `json:"meta"`
	Error   string           `json:"error,omitempty"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

// List lists available games
func (f *GamesFlow) List(ctx context.Context, params ListParams) GamesListResponse {
	req := f.api.GamesAPI.ListGames(ctx)

	if params.Search != "" {
		req = req.Search(params.Search)
	}
	if params.ProducerID > 0 {
		req = req.ProducerId(int32(params.ProducerID))
	}
	if params.Provider != "" {
		req = req.Provider(params.Provider)
	}
	if params.Type != "" {
		req = req.Type_(params.Type)
	}
	if params.PerPage > 0 {
		req = req.PerPage(params.PerPage)
	}

	resp, _, err := req.Execute()
	if err != nil {
		return GamesListResponse{
			Success: false,
			Error:   err.Error(),
			Games:   []Game{},
			Meta:    PaginationMeta{CurrentPage: 1, LastPage: 1, PerPage: 100, Total: 0},
		}
	}

	games := make([]Game, 0)
	if resp.Data != nil {
		for _, g := range resp.Data {
			games = append(games, Game{
				ID:            int(g.GetId()),
				Title:         g.GetTitle(),
				Producer:      g.GetProducer(),
				Type:          g.GetType(),
				ImageURL:      g.GetImageUrl(),
				DemoAvailable: g.GetDemoAvailable(),
			})
		}
	}

	meta := PaginationMeta{
		CurrentPage: 1,
		LastPage:    1,
		PerPage:     100,
		Total:       len(games),
	}
	if resp.Meta != nil {
		meta.CurrentPage = int(resp.Meta.GetCurrentPage())
		meta.LastPage = int(resp.Meta.GetLastPage())
		meta.PerPage = int(resp.Meta.GetPerPage())
		meta.Total = int(resp.Meta.GetTotal())
	}

	return GamesListResponse{
		Success: true,
		Games:   games,
		Meta:    meta,
	}
}

// Get retrieves a single game by ID
func (f *GamesFlow) Get(ctx context.Context, gameID int) ApiResponse {
	_, _, err := f.api.GamesAPI.GetApiV1GamesId(ctx, int32(gameID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"id": gameID},
	}
}

// ByProducer gets games by producer
func (f *GamesFlow) ByProducer(ctx context.Context, producerID int, params ListParams) GamesListResponse {
	params.ProducerID = producerID
	return f.List(ctx, params)
}

// ByCategory gets games by category/type
func (f *GamesFlow) ByCategory(ctx context.Context, gameType string, params ListParams) GamesListResponse {
	params.Type = gameType
	return f.List(ctx, params)
}

// Search searches games by title
func (f *GamesFlow) Search(ctx context.Context, query string, params ListParams) GamesListResponse {
	params.Search = query
	return f.List(ctx, params)
}
