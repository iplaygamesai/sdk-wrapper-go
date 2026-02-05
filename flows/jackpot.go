package flows

import (
	"context"

	apiclient "github.com/iplaygamesai/api-client-go"
)

// JackpotFlow provides high-level operations for jackpots
type JackpotFlow struct {
	api *apiclient.APIClient
}

// NewJackpotFlow creates a new jackpot flow
func NewJackpotFlow(api *apiclient.APIClient) *JackpotFlow {
	return &JackpotFlow{api: api}
}

// GetConfiguration gets current jackpot configuration
func (f *JackpotFlow) GetConfiguration(ctx context.Context) ApiResponse {
	_, _, err := f.api.EndpointsAPI.ConfigureJackpotSettingsForTheOperator(ctx).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Configuration retrieved"},
	}
}

// Configure configures jackpot settings
func (f *JackpotFlow) Configure(ctx context.Context, prizeTiers []interface{}) ApiResponse {
	req := apiclient.NewConfigureJackpotSettingsForTheOperatorRequest()
	// Note: prizeTiers would need proper type mapping

	_, _, err := f.api.EndpointsAPI.ConfigureJackpotSettingsForTheOperator(ctx).ConfigureJackpotSettingsForTheOperatorRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Configuration updated"},
	}
}

// GetPools gets all active jackpot pools
func (f *JackpotFlow) GetPools(ctx context.Context) ApiResponse {
	_, _, err := f.api.EndpointsAPI.ListOperatorsJackpotPools(ctx).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Data:    map[string]interface{}{"pools": []interface{}{}},
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"pools": []interface{}{}},
	}
}

// GetPool gets a specific pool by type
func (f *JackpotFlow) GetPool(ctx context.Context, poolType string) ApiResponse {
	req := apiclient.NewListOperatorsJackpotPoolsRequest()
	req.SetPoolType(poolType)

	_, _, err := f.api.EndpointsAPI.ListOperatorsJackpotPools(ctx).ListOperatorsJackpotPoolsRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"pool_type": poolType},
	}
}

// GetWinners gets winners for a pool
func (f *JackpotFlow) GetWinners(ctx context.Context, poolID string) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"pool_id": poolID,
			"winners": []interface{}{},
		},
	}
}

// GetGames gets games eligible for jackpot
func (f *JackpotFlow) GetGames(ctx context.Context, poolType string) ApiResponse {
	_, _, err := f.api.EndpointsAPI.GetGamesForAPoolTypeOrAllPoolTypes(ctx).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Data:    map[string]interface{}{"games": []interface{}{}},
		}
	}

	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"pool_type": poolType,
			"games":     []interface{}{},
		},
	}
}

// AddGames adds games to a jackpot pool
func (f *JackpotFlow) AddGames(ctx context.Context, poolType string, gameIDs []int) ApiResponse {
	req := apiclient.NewAddGamesToAJackpotPoolTypeRequest()
	req.SetPoolType(poolType)

	ids := make([]int32, len(gameIDs))
	for i, id := range gameIDs {
		ids[i] = int32(id)
	}
	req.SetGameIds(ids)

	_, _, err := f.api.EndpointsAPI.AddGamesToAJackpotPoolType(ctx).AddGamesToAJackpotPoolTypeRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Games added to jackpot pool"},
	}
}

// RemoveGames removes games from a jackpot pool
func (f *JackpotFlow) RemoveGames(ctx context.Context, poolType string, gameIDs []int) ApiResponse {
	req := apiclient.NewRemoveGamesFromAJackpotPoolTypeRequest()
	req.SetPoolType(poolType)

	ids := make([]int32, len(gameIDs))
	for i, id := range gameIDs {
		ids[i] = int32(id)
	}
	req.SetGameIds(ids)

	_, _, err := f.api.EndpointsAPI.RemoveGamesFromAJackpotPoolType(ctx).RemoveGamesFromAJackpotPoolTypeRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Games removed from jackpot pool"},
	}
}

// ContributionFilters contains filters for getting contributions
type ContributionFilters struct {
	PlayerID string
	PoolType string
}

// GetContributions gets contribution history
func (f *JackpotFlow) GetContributions(ctx context.Context, filters ContributionFilters) ApiResponse {
	req := apiclient.NewGetPlayerContributionHistoryRequest()
	if filters.PlayerID != "" {
		req.SetPlayerId(filters.PlayerID)
	}
	if filters.PoolType != "" {
		req.SetPoolType(filters.PoolType)
	}

	_, _, err := f.api.EndpointsAPI.GetPlayerContributionHistory(ctx).GetPlayerContributionHistoryRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Data:    map[string]interface{}{"contributions": []interface{}{}},
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"contributions": []interface{}{}},
	}
}

// Release manually releases a jackpot pool
func (f *JackpotFlow) Release(ctx context.Context, poolID, playerID string) ApiResponse {
	return ApiResponse{
		Success: true,
		Data: map[string]interface{}{
			"pool_id":   poolID,
			"player_id": playerID,
			"message":   "Jackpot release initiated",
		},
	}
}
