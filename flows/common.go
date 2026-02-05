package flows

import (
	"encoding/json"
	"io"
	"net/http"
)

// ApiResponse represents a generic API response
type ApiResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Raw     interface{}            `json:"raw,omitempty"`
}

// parseResponseBody parses the HTTP response body into a map
func parseResponseBody(resp *http.Response) (map[string]interface{}, error) {
	if resp == nil || resp.Body == nil {
		return nil, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// If the response has a "data" field, extract it as the main content
	if data, ok := result["data"]; ok {
		if dataMap, ok := data.(map[string]interface{}); ok {
			return dataMap, nil
		}
		// If data is an array or other type, wrap it
		return map[string]interface{}{"data": data}, nil
	}

	return result, nil
}
