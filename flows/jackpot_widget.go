package flows

import (
	"context"
	"encoding/json"
	"fmt"

	apiclient "github.com/iplaygamesai/api-client-go"
)

// JackpotWidgetFlow provides high-level operations for jackpot widgets
type JackpotWidgetFlow struct {
	api     *apiclient.APIClient
	baseURL string
}

// NewJackpotWidgetFlow creates a new jackpot widget flow
func NewJackpotWidgetFlow(api *apiclient.APIClient, baseURL string) *JackpotWidgetFlow {
	return &JackpotWidgetFlow{api: api, baseURL: baseURL}
}

// RegisterDomain registers a domain for widget embedding
func (f *JackpotWidgetFlow) RegisterDomain(ctx context.Context, domain, name string) ApiResponse {
	req := apiclient.NewRegisterANewDomainRequest(domain)

	resp, _, err := f.api.WidgetManagementAPI.RegisterANewDomain(ctx).RegisterANewDomainRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"domain": resp.Data},
	}
}

// ListDomains lists registered domains
func (f *JackpotWidgetFlow) ListDomains(ctx context.Context) ApiResponse {
	resp, _, err := f.api.WidgetManagementAPI.ListRegisteredDomains(ctx).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Data:    map[string]interface{}{"domains": []interface{}{}},
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"domains": resp.Data},
	}
}

// GetDomain gets domain details
func (f *JackpotWidgetFlow) GetDomain(ctx context.Context, domainID int) ApiResponse {
	resp, _, err := f.api.WidgetManagementAPI.GetDomainDetails(ctx, int32(domainID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"domain": resp.Data},
	}
}

// UpdateDomain updates domain settings
func (f *JackpotWidgetFlow) UpdateDomain(ctx context.Context, domainID int, isActive *bool) ApiResponse {
	req := apiclient.NewUpdateDomainSettingsRequest()
	if isActive != nil {
		req.SetIsActive(*isActive)
	}

	resp, _, err := f.api.WidgetManagementAPI.UpdateDomainSettings(ctx, int32(domainID)).UpdateDomainSettingsRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"domain": resp.Data},
	}
}

// DeleteDomain deletes a domain
func (f *JackpotWidgetFlow) DeleteDomain(ctx context.Context, domainID int) ApiResponse {
	_, _, err := f.api.WidgetManagementAPI.RemoveADomain(ctx, int32(domainID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Domain deleted successfully"},
	}
}

// RegenerateDomainToken regenerates domain token
func (f *JackpotWidgetFlow) RegenerateDomainToken(ctx context.Context, domainID int) ApiResponse {
	resp, _, err := f.api.WidgetManagementAPI.RegenerateDomainToken(ctx, int32(domainID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"domain": resp.Data},
	}
}

// CreateTokenParams contains parameters for creating a widget token
type CreateTokenParams struct {
	DomainToken string
	PlayerID    string
	Currency    string
}

// CreateToken creates a widget token
func (f *JackpotWidgetFlow) CreateToken(ctx context.Context, params CreateTokenParams) ApiResponse {
	req := apiclient.NewGenerateAWidgetTokenRequest(params.DomainToken)
	if params.PlayerID != "" {
		req.SetPlayerId(params.PlayerID)
	}
	if params.Currency != "" {
		req.SetCurrency(params.Currency)
	}

	resp, _, err := f.api.WidgetManagementAPI.GenerateAWidgetToken(ctx).GenerateAWidgetTokenRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"token": resp.Data},
	}
}

// CreateAnonymousToken creates an anonymous widget token
func (f *JackpotWidgetFlow) CreateAnonymousToken(ctx context.Context, domainToken string) ApiResponse {
	return f.CreateToken(ctx, CreateTokenParams{DomainToken: domainToken})
}

// CreatePlayerToken creates a player-specific widget token
func (f *JackpotWidgetFlow) CreatePlayerToken(ctx context.Context, domainToken, playerID, currency string) ApiResponse {
	return f.CreateToken(ctx, CreateTokenParams{
		DomainToken: domainToken,
		PlayerID:    playerID,
		Currency:    currency,
	})
}

// ListTokens lists all widget tokens
func (f *JackpotWidgetFlow) ListTokens(ctx context.Context, domainID *int, active *bool) ApiResponse {
	req := f.api.WidgetManagementAPI.ListWidgetTokens(ctx)
	if domainID != nil {
		req = req.DomainId(int32(*domainID))
	}
	if active != nil {
		req = req.Active(*active)
	}

	resp, _, err := req.Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
			Data:    map[string]interface{}{"tokens": []interface{}{}},
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"tokens": resp.Data},
	}
}

// GetToken gets token details
func (f *JackpotWidgetFlow) GetToken(ctx context.Context, tokenID int) ApiResponse {
	resp, _, err := f.api.WidgetManagementAPI.GetTokenDetails(ctx, int32(tokenID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"token": resp.Data},
	}
}

// RevokeToken revokes a widget token
func (f *JackpotWidgetFlow) RevokeToken(ctx context.Context, tokenID int) ApiResponse {
	_, _, err := f.api.WidgetManagementAPI.RevokeAToken(ctx, int32(tokenID)).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"message": "Token revoked successfully"},
	}
}

// BulkRevokeTokens bulk revokes widget tokens
func (f *JackpotWidgetFlow) BulkRevokeTokens(ctx context.Context, tokenIDs []int) ApiResponse {
	ids := make([]string, len(tokenIDs))
	for i, id := range tokenIDs {
		ids[i] = fmt.Sprintf("%d", id)
	}
	req := apiclient.NewBulkRevokeTokensRequest(ids)

	resp, _, err := f.api.WidgetManagementAPI.BulkRevokeTokens(ctx).BulkRevokeTokensRequest(*req).Execute()
	if err != nil {
		return ApiResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ApiResponse{
		Success: true,
		Data:    map[string]interface{}{"result": resp},
	}
}

// EmbedOptions contains options for generating embed code
type EmbedOptions struct {
	Container string   `json:"container,omitempty"`
	Theme     string   `json:"theme,omitempty"`
	PoolTypes []string `json:"pool_types,omitempty"`
	Token     string   `json:"token"`
}

// GetEmbedCode generates the JavaScript snippet for embedding the widget
func (f *JackpotWidgetFlow) GetEmbedCode(token string, opts EmbedOptions) string {
	opts.Token = token
	if opts.Container == "" {
		opts.Container = "iplaygames-jackpot-widget"
	}

	configJSON, _ := json.Marshal(opts)

	return fmt.Sprintf(`<div id="%s"></div>
<script src="%s/widgets/jackpot.js"></script>
<script>
    IPlayGamesJackpotWidget.init(%s);
</script>`, opts.Container, f.baseURL, string(configJSON))
}
