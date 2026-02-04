package flows

// ApiResponse represents a generic API response
type ApiResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Raw     interface{}            `json:"raw,omitempty"`
}
