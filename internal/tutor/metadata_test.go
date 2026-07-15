package tutor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClaimsToolSupport_OllamaModelWithToolsCapability(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/show" {
			http.NotFound(w, r)
			return
		}
		var req struct {
			Model string `json:"model"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		caps := []string{"completion"}
		if req.Model == "llama3.1:8b" {
			caps = append(caps, "tools")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"capabilities": caps})
	}))
	defer srv.Close()

	claims, known := ClaimsToolSupport(context.Background(), srv.URL, "llama3.1:8b", "")
	if !known || !claims {
		t.Errorf("ClaimsToolSupport(llama3.1:8b) = (%v, %v), want (true, true)", claims, known)
	}

	claims, known = ClaimsToolSupport(context.Background(), srv.URL, "gemma:2b", "")
	if !known || claims {
		t.Errorf("ClaimsToolSupport(gemma:2b) = (%v, %v), want (false, true)", claims, known)
	}
}

func TestClaimsToolSupport_OpenRouterSupportedParameters(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "good/model", "supported_parameters": []string{"temperature", "tools"}},
				{"id": "bad/model", "supported_parameters": []string{"temperature"}},
			},
		})
	}))
	defer srv.Close()
	orig := openRouterBaseURL
	openRouterBaseURL = srv.URL
	t.Cleanup(func() { openRouterBaseURL = orig })

	claims, known := ClaimsToolSupport(context.Background(), "", "openrouter:good/model", "key")
	if !known || !claims {
		t.Errorf("ClaimsToolSupport(good/model) = (%v, %v), want (true, true)", claims, known)
	}

	claims, known = ClaimsToolSupport(context.Background(), "", "openrouter:bad/model", "key")
	if !known || claims {
		t.Errorf("ClaimsToolSupport(bad/model) = (%v, %v), want (false, true)", claims, known)
	}

	// A model the listing doesn't know is unknown, not a denial.
	_, known = ClaimsToolSupport(context.Background(), "", "openrouter:ghost/model", "key")
	if known {
		t.Error("ClaimsToolSupport(ghost/model) known = true, want false for an unlisted model")
	}
}

func TestClaimsToolSupport_FetchFailureMeansUnknown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	if _, known := ClaimsToolSupport(context.Background(), srv.URL, "llama3.1:8b", ""); known {
		t.Error("known = true on a 500 response, want false -- metadata failures must degrade to no signal")
	}

	if _, known := ClaimsToolSupport(context.Background(), "http://127.0.0.1:1", "llama3.1:8b", ""); known {
		t.Error("known = true on a connection failure, want false")
	}
}
