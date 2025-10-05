package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePipeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/pipelines/" {
			t.Errorf("expected path /pipelines/, got %s", r.URL.Path)
		}

		response := PipelineRead{
			ID:       "pipeline-123",
			Name:     "test-pipeline",
			State:    PipelineStateDraft,
			TeamID:   "team-123",
			TeamName: "Test Team",
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	pipeline := &PipelineCreate{
		Name:   "test-pipeline",
		TeamID: "team-123",
		State:  PipelineStateDraft,
	}

	result, err := client.CreatePipeline(context.Background(), pipeline)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "pipeline-123" {
		t.Errorf("expected ID pipeline-123, got %s", result.ID)
	}

	if result.Name != "test-pipeline" {
		t.Errorf("expected Name test-pipeline, got %s", result.Name)
	}

	if result.State != PipelineStateDraft {
		t.Errorf("expected State draft, got %s", result.State)
	}
}

func TestGetPipeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		response := PipelineRead{
			ID:       "pipeline-123",
			Name:     "test-pipeline",
			State:    PipelineStateLive,
			TeamID:   "team-123",
			TeamName: "Test Team",
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	result, err := client.GetPipeline(context.Background(), "pipeline-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
		return
	}

	if result.ID != "pipeline-123" {
		t.Errorf("expected ID pipeline-123, got %s", result.ID)
	}

	if result.State != PipelineStateLive {
		t.Errorf("expected State live, got %s", result.State)
	}
}

func TestGetPipeline_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	result, err := client.GetPipeline(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Errorf("expected nil result for not found, got %v", result)
	}
}

func TestUpdatePipeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH request, got %s", r.Method)
		}

		newName := "updated-pipeline"
		response := PipelineRead{
			ID:       "pipeline-123",
			Name:     newName,
			State:    PipelineStateLive,
			TeamID:   "team-123",
			TeamName: "Test Team",
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	newName := "updated-pipeline"
	update := &PipelineUpdate{
		Name: &newName,
	}

	result, err := client.UpdatePipeline(context.Background(), "pipeline-123", update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != newName {
		t.Errorf("expected Name %s, got %s", newName, result.Name)
	}
}

func TestDeletePipeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeletePipeline(context.Background(), "pipeline-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeletePipeline_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeletePipeline(context.Background(), "nonexistent")
	// Should not return error for not found
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
