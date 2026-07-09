package preflight

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
	if c.Output != `{"models":[]}` {
		t.Errorf("Output = %q, want the raw response body", c.Output)
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

func TestListModels_ReturnsAllPulledModelNames(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(sampleTagsJSON))
	}))
	defer srv.Close()

	names, err := ListModels(srv.URL)
	if err != nil {
		t.Fatalf("ListModels: %v", err)
	}
	want := []string{"qwen2.5-coder:7b", "qwen2.5-coder:1.5b"}
	if len(names) != len(want) {
		t.Fatalf("ListModels = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("ListModels[%d] = %q, want %q", i, names[i], want[i])
		}
	}
}

func TestListModels_EmptyModelsListReturnsEmptySlice(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"models": []}`))
	}))
	defer srv.Close()

	names, err := ListModels(srv.URL)
	if err != nil {
		t.Fatalf("ListModels: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("ListModels = %v, want empty", names)
	}
}

func TestListModels_UnreachableHostReturnsError(t *testing.T) {
	_, err := ListModels("http://127.0.0.1:1")
	if err == nil {
		t.Fatal("expected an error for an unreachable host, got nil")
	}
}

func TestListModels_MalformedJSONReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	_, err := ListModels(srv.URL)
	if err == nil {
		t.Fatal("expected an error for malformed JSON, got nil")
	}
}

func TestListModels_NonOKStatusReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := ListModels(srv.URL)
	if err == nil {
		t.Fatal("expected an error for a non-OK status, got nil")
	}
}

// streamNDJSONServer replies to POST /api/pull with the given lines,
// flushing after each one so the client genuinely reads them as separate
// streaming chunks rather than one buffered response — that's the whole
// point of PullModel over a plain ListModels-style single-body call.
func streamNDJSONServer(t *testing.T, lines []string) (*httptest.Server, *string) {
	t.Helper()
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("httptest ResponseWriter doesn't support flushing")
		}
		for _, line := range lines {
			w.Write([]byte(line + "\n"))
			flusher.Flush()
		}
	}))
	return srv, &gotBody
}

// drainLines collects every line PullModel sends until lineCh closes,
// with a generous timeout so a hung implementation fails the test instead
// of the test suite itself.
func drainLines(t *testing.T, lineCh <-chan string) []string {
	t.Helper()
	var lines []string
	for {
		select {
		case line, ok := <-lineCh:
			if !ok {
				return lines
			}
			lines = append(lines, line)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for lineCh")
		}
	}
}

func TestPullModel_StreamsFormattedStatusLines(t *testing.T) {
	srv, _ := streamNDJSONServer(t, []string{
		`{"status":"pulling manifest"}`,
		`{"status":"verifying sha256 digest"}`,
		`{"status":"success"}`,
	})
	defer srv.Close()

	lineCh, errCh := PullModel(srv.URL, "qwen2.5-coder:7b")
	lines := drainLines(t, lineCh)
	if err := <-errCh; err != nil {
		t.Fatalf("PullModel: %v", err)
	}

	want := []string{"pulling manifest", "verifying sha256 digest", "success"}
	if len(lines) != len(want) {
		t.Fatalf("lines = %v, want %v", lines, want)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("lines[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestPullModel_IncludesPercentageWhenTotalKnown(t *testing.T) {
	srv, _ := streamNDJSONServer(t, []string{
		`{"status":"downloading sha256:abc","total":1000,"completed":250}`,
	})
	defer srv.Close()

	lineCh, errCh := PullModel(srv.URL, "qwen2.5-coder:7b")
	lines := drainLines(t, lineCh)
	if err := <-errCh; err != nil {
		t.Fatalf("PullModel: %v", err)
	}

	if len(lines) != 1 {
		t.Fatalf("lines = %v, want 1 entry", lines)
	}
	if !strings.Contains(lines[0], "25%") {
		t.Errorf("lines[0] = %q, want it to include a 25%% progress figure", lines[0])
	}
}

func TestPullModel_ErrorFieldInStreamSurfacesAsError(t *testing.T) {
	srv, _ := streamNDJSONServer(t, []string{
		`{"status":"pulling manifest"}`,
		`{"error":"model not found"}`,
	})
	defer srv.Close()

	lineCh, errCh := PullModel(srv.URL, "does-not-exist:latest")
	drainLines(t, lineCh)
	err := <-errCh
	if err == nil {
		t.Fatal("expected a non-nil error when the stream reports an error status")
	}
	if !strings.Contains(err.Error(), "model not found") {
		t.Errorf("error = %v, want it to include the stream's error message", err)
	}
}

func TestPullModel_NonOKStatusReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	lineCh, errCh := PullModel(srv.URL, "qwen2.5-coder:7b")
	drainLines(t, lineCh)
	if err := <-errCh; err == nil {
		t.Fatal("expected an error for a non-OK status, got nil")
	}
}

func TestPullModel_UnreachableHostReturnsError(t *testing.T) {
	lineCh, errCh := PullModel("http://127.0.0.1:1", "qwen2.5-coder:7b")
	drainLines(t, lineCh)
	if err := <-errCh; err == nil {
		t.Fatal("expected an error for an unreachable host, got nil")
	}
}

func TestPullModel_SendsModelNameInRequestBody(t *testing.T) {
	srv, gotBody := streamNDJSONServer(t, []string{`{"status":"success"}`})
	defer srv.Close()

	lineCh, errCh := PullModel(srv.URL, "qwen2.5-coder:7b")
	drainLines(t, lineCh)
	<-errCh

	var payload struct {
		Model string `json:"model"`
	}
	if err := json.Unmarshal([]byte(*gotBody), &payload); err != nil {
		t.Fatalf("request body isn't valid JSON: %v (%q)", err, *gotBody)
	}
	if payload.Model != "qwen2.5-coder:7b" {
		t.Errorf("request model = %q, want %q", payload.Model, "qwen2.5-coder:7b")
	}
}
