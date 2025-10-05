package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PipelineState represents the state of a pipeline
type PipelineState string

const (
	PipelineStateDraft    PipelineState = "draft"
	PipelineStatePaused   PipelineState = "paused"
	PipelineStateLive     PipelineState = "live"
	PipelineStateError    PipelineState = "error"
	PipelineStateBuilding PipelineState = "building"
)

// PipelineConfiguration represents the configuration structure for a pipeline
type PipelineConfiguration struct {
	SourceName   string         `json:"source_name"`
	SourceType   *string        `json:"source_type,omitempty"`
	SourceConfig map[string]any `json:"source_config"`
	TargetName   string         `json:"target_name"`
	TargetType   *string        `json:"target_type,omitempty"`
	TargetConfig map[string]any `json:"target_config"`
	SMTName      string         `json:"smt_name"`
	SMTConfig    []any          `json:"smt_config"`
	DraftStep    string         `json:"draft_step"`
}

// PipelineCreate represents the request to create a pipeline
type PipelineCreate struct {
	Name              string                 `json:"name"`
	TeamID            string                 `json:"team_id"`
	State             PipelineState          `json:"state"`
	JSONConfiguration *PipelineConfiguration `json:"json_configuration"`
}

// PipelineUpdate represents the request to update a pipeline
type PipelineUpdate struct {
	Name              *string                `json:"name,omitempty"`
	TeamID            *string                `json:"team_id,omitempty"`
	State             *PipelineState         `json:"state,omitempty"`
	JSONConfiguration *PipelineConfiguration `json:"json_configuration,omitempty"`
}

// PipelineRead represents a pipeline response
type PipelineRead struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	State             PipelineState          `json:"state"`
	TeamID            string                 `json:"team_id"`
	TeamName          string                 `json:"team_name"`
	JSONConfiguration *PipelineConfiguration `json:"json_configuration"`
}

// CreatePipeline creates a new pipeline
func (c *Client) CreatePipeline(ctx context.Context, pipeline *PipelineCreate) (*PipelineRead, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, "/pipelines/", pipeline)
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

	var result PipelineRead
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetPipeline retrieves a pipeline by ID
func (c *Client) GetPipeline(ctx context.Context, pipelineID string) (*PipelineRead, error) {
	path := fmt.Sprintf("/pipelines/%s", pipelineID)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden {
		return nil, nil
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result PipelineRead
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdatePipeline updates an existing pipeline
func (c *Client) UpdatePipeline(ctx context.Context, pipelineID string, pipeline *PipelineUpdate) (*PipelineRead, error) {
	path := fmt.Sprintf("/pipelines/%s", pipelineID)
	resp, err := c.doRequest(ctx, http.MethodPatch, path, pipeline)
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

	var result PipelineRead
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeletePipeline deletes a pipeline by ID
func (c *Client) DeletePipeline(ctx context.Context, pipelineID string) error {
	path := fmt.Sprintf("/pipelines/%s", pipelineID)
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if err := checkResponse(resp); err != nil {
		return err
	}

	return nil
}
