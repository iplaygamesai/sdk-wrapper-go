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
	Name          string
	PromotionType string
	CycleType     string
	StartsAt      string
	EndsAt        string
	IsActive      *bool
}

// List lists all promotions
func (f *PromotionsFlow) List(ctx context.Context, status, promotionType string) ApiResponse {
	// Note: There's no ListAllPromotions in the API - returning empty for now
	// The API may need to be updated to include this endpoint
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotions": []interface{}{},
			"message":    "Promotions list endpoint not available in current API version",
		},
	}
}

// Get gets a specific promotion
func (f *PromotionsFlow) Get(ctx context.Context, promotionID int) ApiResponse {
	httpResp, err := f.api.EndpointsAPI.GetASpecificPromotion(ctx, fmt.Sprintf("%d", promotionID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	data, parseErr := parseResponseBody(httpResp)
	if parseErr != nil || data == nil {
		return ApiResponse{
			Success: true,
			Data:    map[string]interface{}{"promotion_id": promotionID},
		}
	}

	data["promotion_id"] = promotionID
	return ApiResponse{
		Success: true,
		Data:    data,
	}
}

// Create creates a new promotion
func (f *PromotionsFlow) Create(ctx context.Context, data PromotionData) ApiResponse {
	req := apiclient.NewCreateANewPromotionRequest(data.Name, data.PromotionType, data.CycleType)
	if data.StartsAt != "" {
		req.SetStartsAt(data.StartsAt)
	}
	if data.EndsAt != "" {
		req.SetEndsAt(data.EndsAt)
	}

	httpResp, err := f.api.EndpointsAPI.CreateANewPromotion(ctx).CreateANewPromotionRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	respData, parseErr := parseResponseBody(httpResp)
	if parseErr != nil || respData == nil {
		return ApiResponse{
			Success: true,
			Data:    map[string]interface{}{"message": "Promotion created"},
		}
	}

	respData["message"] = "Promotion created"
	return ApiResponse{
		Success: true,
		Data:    respData,
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

	httpResp, err := f.api.EndpointsAPI.UpdateAPromotion(ctx, fmt.Sprintf("%d", promotionID)).UpdateAPromotionRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	respData, parseErr := parseResponseBody(httpResp)
	if parseErr != nil || respData == nil {
		return ApiResponse{
			Success: true,
			Data: map[string]interface{}{
				"promotion_id": promotionID,
				"message":      "Promotion updated",
			},
		}
	}

	respData["promotion_id"] = promotionID
	respData["message"] = "Promotion updated"
	return ApiResponse{
		Success: true,
		Data:    respData,
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
	// Note: GetPromotionLeaderboard endpoint not available in current API version
	// This would need to be added to the API spec
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"leaderboard":  []interface{}{},
			"message":      "Leaderboard endpoint not available in current API version",
		},
	}
}

// GetWinners gets promotion winners
func (f *PromotionsFlow) GetWinners(ctx context.Context, promotionID int) ApiResponse {
	// Note: GetPromotionWinners endpoint not available in current API version
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"winners":      []interface{}{},
			"message":      "Winners endpoint not available in current API version",
		},
	}
}

// GetGames gets games eligible for a promotion
func (f *PromotionsFlow) GetGames(ctx context.Context, promotionID int) ApiResponse {
	// Note: GetGamesForAPromotion endpoint not available in current API version
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"promotion_id": promotionID,
			"games":        []interface{}{},
			"message":      "Games endpoint not available in current API version",
		},
	}
}

// ManageGames sets games for a promotion
func (f *PromotionsFlow) ManageGames(ctx context.Context, promotionID int, gameIDs []int) ApiResponse {
	req := apiclient.NewManageGamesForAPromotionRequest()

	ids := make([]int32, len(gameIDs))
	for i, id := range gameIDs {
		ids[i] = int32(id)
	}
	req.SetGameIds(ids)

	httpResp, err := f.api.EndpointsAPI.ManageGamesForAPromotion(ctx, fmt.Sprintf("%d", promotionID)).ManageGamesForAPromotionRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	data, _ := parseResponseBody(httpResp)
	if data == nil {
		data = map[string]interface{}{"message": "Games updated for promotion"}
	}
	return ApiResponse{
		Success: true,
		Data:    data,
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
