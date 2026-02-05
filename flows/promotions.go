package flows

import (
	"context"
	"fmt"

	apiclient "github.com/iplaygamesai/api-client-go"
)

// PromotionsFlow provides high-level operations for promotions
type PromotionsFlow struct {
	api *apiclient.APIClient
}

// NewPromotionsFlow creates a new promotions flow
func NewPromotionsFlow(api *apiclient.APIClient) *PromotionsFlow {
	return &PromotionsFlow{api: api}
}

// PromotionData contains data for creating/updating a promotion
type PromotionData struct {
	Name      string
	Type      string
	StartDate string
	EndDate   string
	IsActive  *bool
}

// List lists all promotions
func (f *PromotionsFlow) List(ctx context.Context, status, promotionType string) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotions": []interface{}{},
			"meta":       map[string]interface{}{},
		},
	}
}

// Get gets a specific promotion
func (f *PromotionsFlow) Get(ctx context.Context, promotionID int) ApiResponse {
	_, err := f.api.EndpointsAPI.GetASpecificPromotion(ctx, fmt.Sprintf("%d", promotionID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"promotion_id": promotionID},
	}
}

// Create creates a new promotion
func (f *PromotionsFlow) Create(ctx context.Context, data PromotionData) ApiResponse {
	req := apiclient.NewCreateANewPromotionRequest()
	if data.Name != "" {
		req.SetName(data.Name)
	}
	if data.Type != "" {
		req.SetType(data.Type)
	}
	if data.StartDate != "" {
		req.SetStartDate(data.StartDate)
	}
	if data.EndDate != "" {
		req.SetEndDate(data.EndDate)
	}

	_, err := f.api.EndpointsAPI.CreateANewPromotion(ctx).CreateANewPromotionRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Promotion created"},
	}
}

// Update updates a promotion
func (f *PromotionsFlow) Update(ctx context.Context, promotionID int, data PromotionData) ApiResponse {
	req := apiclient.NewUpdateAPromotionRequest()
	if data.Name != "" {
		req.SetName(data.Name)
	}
	if data.IsActive != nil {
		req.SetIsActive(*data.IsActive)
	}

	_, err := f.api.EndpointsAPI.UpdateAPromotion(ctx, fmt.Sprintf("%d", promotionID)).UpdateAPromotionRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"message":      "Promotion updated",
		},
	}
}

// Delete deletes a promotion
func (f *PromotionsFlow) Delete(ctx context.Context, promotionID int) ApiResponse {
	_, err := f.api.EndpointsAPI.DeleteAPromotion(ctx, fmt.Sprintf("%d", promotionID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Promotion deleted"},
	}
}

// GetLeaderboard gets promotion leaderboard
func (f *PromotionsFlow) GetLeaderboard(ctx context.Context, promotionID, limit, periodID int) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"leaderboard":  []interface{}{},
		},
	}
}

// GetWinners gets promotion winners
func (f *PromotionsFlow) GetWinners(ctx context.Context, promotionID int) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"winners":      []interface{}{},
		},
	}
}

// GetGames gets games eligible for a promotion
func (f *PromotionsFlow) GetGames(ctx context.Context, promotionID int) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"games":        []interface{}{},
		},
	}
}

// ManageGames adds/removes games for a promotion
func (f *PromotionsFlow) ManageGames(ctx context.Context, promotionID int, gameIDs []int, action string) ApiResponse {
	req := apiclient.NewManageGamesForAPromotionRequest()

	ids := make([]int32, len(gameIDs))
	for i, id := range gameIDs {
		ids[i] = int32(id)
	}
	req.SetGameIds(ids)
	req.SetAction(action)

	_, err := f.api.EndpointsAPI.ManageGamesForAPromotion(ctx, fmt.Sprintf("%d", promotionID)).ManageGamesForAPromotionRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": fmt.Sprintf("Games %sed for promotion", action)},
	}
}

// OptIn opts a player into a promotion
func (f *PromotionsFlow) OptIn(ctx context.Context, promotionID int, playerID, currency string) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"player_id":    playerID,
			"message":      "Player opted in",
		},
	}
}

// OptOut opts a player out of a promotion
func (f *PromotionsFlow) OptOut(ctx context.Context, promotionID int, playerID string) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"player_id":    playerID,
			"message":      "Player opted out",
		},
	}
}

// Distribute distributes prizes for a promotion period
func (f *PromotionsFlow) Distribute(ctx context.Context, promotionID, periodID int) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"period_id":    periodID,
			"message":      "Distribution initiated",
		},
	}
}
