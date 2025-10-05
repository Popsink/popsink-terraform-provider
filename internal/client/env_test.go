package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateEnv(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/envs/" {
			t.Errorf("expected path /envs/, got %s", r.URL.Path)
		}

		response := EnvRead{
			ID:           "env-123",
			Name:         "test-env",
			UseRetention: true,
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	env := &EnvCreate{
		Name:         "test-env",
		UseRetention: true,
	}

	result, err := client.CreateEnv(context.Background(), env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "env-123" {
		t.Errorf("expected ID env-123, got %s", result.ID)
	}

	if result.Name != "test-env" {
		t.Errorf("expected Name test-env, got %s", result.Name)
	}

	if !result.UseRetention {
		t.Error("expected UseRetention to be true")
	}
}

func TestGetEnv(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		response := EnvRead{
			ID:           "env-123",
			Name:         "test-env",
			UseRetention: false,
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	result, err := client.GetEnv(context.Background(), "env-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
		return
	}

	if result.ID != "env-123" {
		t.Errorf("expected ID env-123, got %s", result.ID)
	}
}

func TestGetEnv_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	result, err := client.GetEnv(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Errorf("expected nil result for not found, got %v", result)
	}
}

func TestUpdateEnv(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH request, got %s", r.Method)
		}

		newName := "updated-env"
		response := EnvRead{
			ID:           "env-123",
			Name:         newName,
			UseRetention: true,
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	newName := "updated-env"
	update := &EnvUpdate{
		Name: &newName,
	}

	result, err := client.UpdateEnv(context.Background(), "env-123", update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != newName {
		t.Errorf("expected Name %s, got %s", newName, result.Name)
	}
}

func TestDeleteEnv(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteEnv(context.Background(), "env-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteEnv_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	err := client.DeleteEnv(context.Background(), "nonexistent")
	// Should not return error for 404
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
