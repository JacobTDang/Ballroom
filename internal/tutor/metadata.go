package tutor

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// ClaimsToolSupport reports whether the model's PROVIDER METADATA
// claims tool-calling support: Ollama's /api/show capabilities list, or
// OpenRouter's /models supported_parameters. known is false when the
// metadata couldn't be fetched or doesn't list the model -- callers
// must treat that as "no signal", never as a denial.
//
// Metadata can rule a model OUT (a model that doesn't even claim tools
// will never pass the live probe) but never IN: verified live 2026-07-14
// against openrouter:poolside/laguna-xs-2.1:free, whose listing claims
// tools yet which failed 6/6 real probe calls. CheckToolCalling stays
// the ground truth; this is a cheap advisory pre-filter for the model
// picker.
func ClaimsToolSupport(ctx context.Context, ollamaHost, modelName, apiKey string) (claims, known bool) {
	if strings.HasPrefix(modelName, OpenRouterModelPrefix) {
		return openRouterClaimsTools(ctx, strings.TrimPrefix(modelName, OpenRouterModelPrefix), apiKey)
	}
	return ollamaClaimsTools(ctx, ollamaHost, modelName)
}

func ollamaClaimsTools(ctx context.Context, ollamaHost, modelName string) (claims, known bool) {
	body := strings.NewReader(`{"model": ` + string(mustJSON(modelName)) + `}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaHost+"/api/show", body)
	if err != nil {
		return false, false
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, false
	}
	var parsed struct {
		Capabilities []string `json:"capabilities"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return false, false
	}
	for _, c := range parsed.Capabilities {
		if c == "tools" {
			return true, true
		}
	}
	return false, true
}

func openRouterClaimsTools(ctx context.Context, slug, apiKey string) (claims, known bool) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, openRouterBaseURL+"/models", nil)
	if err != nil {
		return false, false
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, false
	}
	var parsed struct {
		Data []struct {
			ID                  string   `json:"id"`
			SupportedParameters []string `json:"supported_parameters"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return false, false
	}
	for _, m := range parsed.Data {
		if m.ID != slug {
			continue
		}
		for _, p := range m.SupportedParameters {
			if p == "tools" {
				return true, true
			}
		}
		return false, true
	}
	return false, false
}

// mustJSON marshals a string for safe embedding in a JSON body --
// model tags are caller-controlled input, not format-string material.
func mustJSON(s string) []byte {
	b, err := json.Marshal(s)
	if err != nil {
		// json.Marshal of a string cannot fail.
		panic(err)
	}
	return b
}
