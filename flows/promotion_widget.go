package flows

import (
	"context"
	"encoding/json"
	"fmt"

	apiclient "github.com/iplaygamesai/api-client-go"
)

// PromotionWidgetFlow provides high-level operations for promotion widgets
type PromotionWidgetFlow struct {
	api     *apiclient.APIClient
	baseURL string
}

// NewPromotionWidgetFlow creates a new promotion widget flow
func NewPromotionWidgetFlow(api *apiclient.APIClient, baseURL string) *PromotionWidgetFlow {
	return &PromotionWidgetFlow{api: api, baseURL: baseURL}
}

// RegisterDomain registers a domain for widget embedding
func (f *PromotionWidgetFlow) RegisterDomain(ctx context.Context, domain, name string) ApiResponse {
	req := apiclient.NewRegisterANewDomainRequest()
	req.SetDomain(domain)
	if name != "" {
		req.SetName(name)
	}

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
func (f *PromotionWidgetFlow) ListDomains(ctx context.Context) ApiResponse {
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

// CreateToken creates a widget token for promotions
func (f *PromotionWidgetFlow) CreateToken(ctx context.Context, domainToken, playerID, currency string) ApiResponse {
	req := apiclient.NewGenerateAWidgetTokenRequest()
	req.SetDomainToken(domainToken)
	if playerID != "" {
		req.SetPlayerId(playerID)
	}
	if currency != "" {
		req.SetCurrency(currency)
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
func (f *PromotionWidgetFlow) CreateAnonymousToken(ctx context.Context, domainToken string) ApiResponse {
	return f.CreateToken(ctx, domainToken, "", "")
}

// CreatePlayerToken creates a player-specific widget token
func (f *PromotionWidgetFlow) CreatePlayerToken(ctx context.Context, domainToken, playerID, currency string) ApiResponse {
	return f.CreateToken(ctx, domainToken, playerID, currency)
}

// ListTokens lists all widget tokens
func (f *PromotionWidgetFlow) ListTokens(ctx context.Context, domainID *int, active *bool) ApiResponse {
	req := f.api.WidgetManagementAPI.ListWidgetTokens(ctx)
	if domainID != nil {
		req = req.DomainId(fmt.Sprintf("%d", *domainID))
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

// RevokeToken revokes a widget token
func (f *PromotionWidgetFlow) RevokeToken(ctx context.Context, tokenID int) ApiResponse {
	_, err := f.api.WidgetManagementAPI.RevokeAToken(ctx, fmt.Sprintf("%d", tokenID)).Execute()
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

// PromotionEmbedOptions contains options for generating embed code
type PromotionEmbedOptions struct {
	Container    string `json:"container,omitempty"`
	Theme        string `json:"theme,omitempty"`
	PromotionIDs []int  `json:"promotion_ids,omitempty"`
	Token        string `json:"token"`
}

// GetEmbedCode generates the JavaScript snippet for embedding the widget
func (f *PromotionWidgetFlow) GetEmbedCode(token string, opts PromotionEmbedOptions) string {
	opts.Token = token
	if opts.Container == "" {
		opts.Container = "iplaygames-promotion-widget"
	}

	configJSON, _ := json.Marshal(opts)

	return fmt.Sprintf(`<div id="%s"></div>
<script src="%s/widgets/promotions.js"></script>
<script>
    IPlayGamesPromotionWidget.init(%s);
</script>`, opts.Container, f.baseURL, string(configJSON))
}
