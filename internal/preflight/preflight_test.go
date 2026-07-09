package preflight

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const sampleTagsJSON = `{
  "models": [
    {"name": "qwen2.5-coder:7b", "model": "qwen2.5-coder:7b", "size": 4683087561},
    {"name": "qwen2.5-coder:1.5b", "model": "qwen2.5-coder:1.5b", "size": 986000000}
  ]
}`

func TestModelPresent_FindsExactMatch(t *testing.T) {
	if !modelPresent([]byte(sampleTagsJSON), "qwen2.5-coder:7b") {
		t.Error("expected qwen2.5-coder:7b to be found in the sample response")
	}
}

func TestModelPresent_MissingModelReturnsFalse(t *testing.T) {
	if modelPresent([]byte(sampleTagsJSON), "llama3:70b") {
		t.Error("expected llama3:70b to NOT be found in the sample response")
	}
}

func TestModelPresent_EmptyModelsList(t *testing.T) {
	if modelPresent([]byte(`{"models": []}`), "qwen2.5-coder:7b") {
		t.Error("expected no match against an empty models list")
	}
}

func TestModelPresent_MalformedJSON(t *testing.T) {
	if modelPresent([]byte(`not json`), "qwen2.5-coder:7b") {
		t.Error("expected malformed JSON to be treated as no match, not a panic/crash")
	}
}

func TestCheckDocker_ReportsStructuredResult(t *testing.T) {
	// Doesn't assert OK/not-OK (depends on the test machine's Docker state)
	// — just that it returns a well-formed result instead of panicking.
	c := CheckDocker()
	if c.Name == "" {
		t.Error("expected a non-empty check name")
	}
	if c.Command != "docker info" {
		t.Errorf("Command = %q, want %q", c.Command, "docker info")
	}
}

func TestCheckImage_UnknownImageReportsNotOK(t *testing.T) {
	c := CheckImage("this-image-definitely-does-not-exist-12345")
	if c.OK {
		t.Error("expected OK=false for a nonexistent image")
	}
	if c.Detail == "" {
		t.Error("expected a non-empty detail message explaining how to fix it")
	}
	want := `docker image inspect this-image-definitely-does-not-exist-12345 --format "{{.Id}}"`
	if c.Command != want {
		t.Errorf("Command = %q, want %q", c.Command, want)
	}
}

func TestCheckOllama_UnreachableHostReportsNotOK(t *testing.T) {
	c := CheckOllama("http://127.0.0.1:1")
	if c.OK {
		t.Error("expected OK=false for an unreachable host")
	}
	want := "GET http://127.0.0.1:1/api/tags"
	if c.Command != want {
		t.Errorf("Command = %q, want %q", c.Command, want)
	}
}

func TestCheckOllama_ReachableHostReportsOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"models":[]}`))
	}))
	defer srv.Close()

	c := CheckOllama(srv.URL)
	if !c.OK {
		t.Errorf("expected OK=true for a reachable host, got Detail=%q", c.Detail)
	}
}

func TestCheckModel_ReportsOKWhenPresent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(sampleTagsJSON))
	}))
	defer srv.Close()

	c := CheckModel(srv.URL, "qwen2.5-coder:7b")
	if !c.OK {
		t.Errorf("expected OK=true, got Detail=%q", c.Detail)
	}
	want := "GET " + srv.URL + "/api/tags"
	if c.Command != want {
		t.Errorf("Command = %q, want %q", c.Command, want)
	}
}

func TestCheckModel_ReportsNotOKWhenMissing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(sampleTagsJSON))
	}))
	defer srv.Close()

	c := CheckModel(srv.URL, "llama3:70b")
	if c.OK {
		t.Error("expected OK=false for a model that isn't pulled")
	}
	if c.Detail == "" {
		t.Error("expected a non-empty detail message with the pull command")
	}
}
