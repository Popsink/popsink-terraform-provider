package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// BrokerConfiguration represents the retention configuration for an environment
type BrokerConfiguration map[string]any

// EnvCreate represents the request structure for creating an environment
type EnvCreate struct {
	Name                   string               `json:"name"`
	UseRetention           bool                 `json:"use_retention"`
	RetentionConfiguration *BrokerConfiguration `json:"retention_configuration,omitempty"`
}

// EnvRead represents the response structure for reading an environment
type EnvRead struct {
	ID                     string               `json:"id"`
	Name                   string               `json:"name"`
	UseRetention           bool                 `json:"use_retention"`
	RetentionConfiguration *BrokerConfiguration `json:"retention_configuration,omitempty"`
}

// EnvUpdate represents the request structure for updating an environment
type EnvUpdate struct {
	Name                   *string              `json:"name,omitempty"`
	UseRetention           *bool                `json:"use_retention,omitempty"`
	RetentionConfiguration *BrokerConfiguration `json:"retention_configuration,omitempty"`
}

// CreateEnv creates a new environment
func (c *Client) CreateEnv(ctx context.Context, env *EnvCreate) (*EnvRead, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/envs/", env)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result EnvRead
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetEnv retrieves an environment by ID
func (c *Client) GetEnv(ctx context.Context, envID string) (*EnvRead, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/envs/"+envID, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Environment not found
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result EnvRead
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateEnv updates an existing environment
func (c *Client) UpdateEnv(ctx context.Context, envID string, env *EnvUpdate) (*EnvRead, error) {
	resp, err := c.doRequest(ctx, http.MethodPatch, "/envs/"+envID, env)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result EnvRead
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Note: DeleteEnv is not implemented as the API doesn't seem to provide a delete endpoint for environments
// Based on the OpenAPI specification, environments appear to be managed through other means
