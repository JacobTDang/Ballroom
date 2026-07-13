package tutor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// ansiPattern/stripAnsi strip terminal escape sequences (cursor
// movement, color) out of captured stdout so tests can assert on the
// visible text content — shared by discoball_test.go, thinkingdisplay_test.go,
// and this file's own Run-level integration tests.
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

func stripAnsi(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

// tutorChatRequest/sequencedOllama are the package's Ollama mock server
// for Go-native tests — used here and by agent_test.go. Originally kept
// independent of chat_test.go's now-deleted mockOllama (which tested
// tutor/chat.sh, the bash implementation this package replaced) so this
// file wouldn't need to change when that one went away.
type tutorChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type tutorChatRequest struct {
	Model    string             `json:"model"`
	Messages []tutorChatMessage `json:"messages"`
	Stream   bool               `json:"stream"`
}

// sequencedOllama serves one reply per request in order, repeating the
// last reply if more requests arrive than replies were given — lets a
// test simulate a multi-turn conversation (e.g. the comprehension
// check's reply, then a real answer) without needing a real model.
type sequencedOllama struct {
	*httptest.Server
	replies []string

	mu       sync.Mutex
	requests []tutorChatRequest
}

func newSequencedOllama(t *testing.T, replies ...string) *sequencedOllama {
	t.Helper()
	m := &sequencedOllama{replies: replies}
	m.Server = httptest.NewServer(http.HandlerFunc(m.handle))
	t.Cleanup(m.Server.Close)
	return m
}

func (m *sequencedOllama) handle(w http.ResponseWriter, r *http.Request) {
	var req tutorChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	idx := len(m.requests)
	m.requests = append(m.requests, req)
	m.mu.Unlock()

	reply := m.replies[len(m.replies)-1]
	if idx < len(m.replies) {
		reply = m.replies[idx]
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": map[string]string{"role": "assistant", "content": reply},
		"done":    true,
	})
}

func (m *sequencedOllama) allRequests() []tutorChatRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]tutorChatRequest, len(m.requests))
	copy(out, m.requests)
	return out
}

func testConfig(ollamaHost string) Config {
	return Config{
		OllamaHost:      ollamaHost,
		Model:           "test-model",
		Mode:            exercise.TutorModeFullAssist,
		MaxContextBytes: 8000,
	}
}

func TestRun_PrintsAssistantReply(t *testing.T) {
	mock := newSequencedOllama(t, "the answer is 42")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check for a single-turn test

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("question\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !strings.Contains(stdout.String(), "the answer is 42") {
		t.Errorf("stdout = %q, want it to contain the assistant reply", stdout.String())
	}
}

func TestRun_CompletesATurnThatMakesARealToolCall(t *testing.T) {
	// newToolCallOllama (toolcheck_test.go) simulates a real tool_calls
	// response for its first request, then a plain-text reply for its
	// second — driving Run() through an actual read_solution_file call
	// via the agent, not a synthetic call.
	mock := newToolCallOllama(t, "read_solution_file")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check for a single-turn test
	cfg.WorkDir = t.TempDir()

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("what does my code look like?\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	got := stripAnsi(stdout.String())
	if !strings.Contains(got, "pong received") {
		t.Fatalf("expected the final reply in stdout, got:\n%s", got)
	}
}

// fakeInputBoxAt substitutes newInputBoxFn so Run() gets a real,
// deterministically-sized inputBox instead of failing to find one (every
// test's stdin is not a real tty, so newInputBox would otherwise always
// return box == nil, and the activity display's real box-drawing path
// would never be exercised). Returns a restore func to defer.
func fakeInputBoxAt(rows, cols int) func() {
	orig := newInputBoxFn
	newInputBoxFn = func(w io.Writer) (*inputBox, error) {
		return newInputBoxAt(w, rows, cols)
	}
	return func() { newInputBoxFn = orig }
}

func TestRun_ShowsLiveToolCallActivityAndClearsItBeforeThePrints(t *testing.T) {
	defer fakeInputBoxAt(24, 80)()

	// newToolCallOllama (toolcheck_test.go) simulates a real tool_calls
	// response for its first request, then a plain-text reply for its
	// second -- driving Run() through an actual read_solution_file call
	// via the agent (and this test's real inputBox/newActivityOption
	// wiring), not a synthetic call.
	mock := newToolCallOllama(t, "read_solution_file")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check for a single-turn test
	cfg.WorkDir = t.TempDir()

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("what does my code look like?\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	out := stdout.String()
	toolLineIdx := strings.Index(out, "→ read_solution_file")
	if toolLineIdx == -1 {
		t.Fatalf("expected the live activity display to show the real tool being called, got:\n%s", out)
	}

	// The exact byte sequence clearActivity's single buffered write
	// produces for a 24x80 box (regionBottom=16): all 5 activity rows
	// positioned and cleared back to back, with nothing else interleaved.
	clearSeq := "\033[17;1H\033[2K\033[18;1H\033[2K\033[19;1H\033[2K\033[20;1H\033[2K\033[21;1H\033[2K"
	clearIdx := strings.Index(out, clearSeq)
	if clearIdx == -1 {
		t.Fatalf("expected the activity region to be cleared once the turn finished, got:\n%s", out)
	}

	replyIdx := strings.Index(out, "pong received")
	if replyIdx == -1 {
		t.Fatalf("expected the final reply in stdout, got:\n%s", out)
	}

	if !(toolLineIdx < clearIdx && clearIdx < replyIdx) {
		t.Errorf("expected order tool-call shown (%d) -> activity cleared (%d) -> reply printed (%d), got:\n%s", toolLineIdx, clearIdx, replyIdx, out)
	}
}

func TestLooksLikeLeakedToolCall_DetectsRawJSONShape(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    bool
	}{
		{"bare JSON", `{"name": "read_solution_file", "parameters": {}}`, true},
		{"prose then JSON", "To answer this, I need to check your code.\n\n" + `{"name": "read_cursor_position", "parameters": {}}`, true},
		{"hallucinated tool name still matches the shape", `{"name": "read_user_code", "parameters": {}}`, true},
		{"clean reply", "your code looks fine to me", false},
		{"empty", "", false},
	}
	for _, c := range cases {
		if got := looksLikeLeakedToolCall(c.content); got != c.want {
			t.Errorf("%s: looksLikeLeakedToolCall(%q) = %v, want %v", c.name, c.content, got, c.want)
		}
	}
}

func TestRun_RetriesWhenReplyLeaksFakeToolCallJSON(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_solution_file", "parameters": {}}`, "your code looks fine")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check for a single-turn test

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("what does my code look like?\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	got := stdout.String()
	if strings.Contains(got, `{"name"`) {
		t.Errorf("stdout still contains leaked tool-call JSON: %q", got)
	}
	if !strings.Contains(got, "your code looks fine") {
		t.Errorf("stdout = %q, want the retry's clean reply", got)
	}
	if n := len(mock.allRequests()); n != 2 {
		t.Errorf("requests = %d, want exactly 2 (original + one retry)", n)
	}
}

func TestRun_FallsBackToHonestMessageWhenRetryAlsoLeaks(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_solution_file", "parameters": {}}`, `{"name": "read_cursor_position", "parameters": {}}`)
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("where is my cursor?\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	got := stdout.String()
	if strings.Contains(got, `{"name"`) {
		t.Errorf("stdout still contains leaked tool-call JSON: %q", got)
	}
	if !strings.Contains(got, leakedToolCallFallbackReply) {
		t.Errorf("stdout = %q, want the honest fallback message", got)
	}
}

func TestRun_DoesNotRetryWhenReplyIsClean(t *testing.T) {
	mock := newSequencedOllama(t, "the answer is 42")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("question\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if n := len(mock.allRequests()); n != 1 {
		t.Errorf("requests = %d, want exactly 1 (no retry for a clean reply)", n)
	}
}

func TestRun_LeakedReplyNeverPersistedToHistory(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_solution_file", "parameters": {}}`, "your code looks fine", "second reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout, stderr strings.Builder
	input := "what does my code look like?\nanother question\n"
	if err := Run(context.Background(), cfg, strings.NewReader(input), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) != 3 {
		t.Fatalf("requests = %d, want exactly 3 (leaked original + retry + second turn)", len(reqs))
	}
	// The second turn's request carries history from the first turn --
	// confirm the leaked (never-shown) reply isn't in it, only the
	// clean retry reply.
	for _, m := range reqs[2].Messages {
		if strings.Contains(m.Content, `{"name"`) {
			t.Errorf("second turn's request carries leaked JSON in history: %+v", reqs[2].Messages)
		}
	}
}

func TestRun_RetainsConversationHistory(t *testing.T) {
	mock := newSequencedOllama(t, "assistant-reply-1", "assistant-reply-2")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check to isolate history behavior

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("first line\nsecond line\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requests (one per input line), got %d", len(reqs))
	}

	if len(reqs[0].Messages) != 2 {
		t.Errorf("first request: expected [system, user1] = 2 messages, got %d: %+v", len(reqs[0].Messages), reqs[0].Messages)
	}

	second := reqs[1].Messages
	if len(second) != 4 {
		t.Fatalf("second request: expected [system, user1, assistant1, user2] = 4 messages, got %d: %+v", len(second), second)
	}
	if second[1].Content != "first line" {
		t.Errorf("second request messages[1] (user1) = %q, want %q", second[1].Content, "first line")
	}
	if second[2].Role != "assistant" || second[2].Content != "assistant-reply-1" {
		t.Errorf("second request messages[2] (assistant1) = %+v, want role=assistant content=%q", second[2], "assistant-reply-1")
	}
	if second[3].Content != "second line" {
		t.Errorf("second request messages[3] (user2) = %q, want %q", second[3].Content, "second line")
	}
}

// TestRun_TurnFailureShowsUserFacingMessageNotJustStderr is a regression
// test for a real bug found live: when a turn's Generate call fails (a
// bad host, a rejected request, an upstream rate limit, anything), the
// old code only logged to stderr and silently continued the loop --
// nothing was ever printed to stdout, so the user's chat pane showed no
// acknowledgment at all that their message was received or that
// anything went wrong. A transient failure looked identical to "the
// tutor is completely unresponsive". Every turn must show the user
// something, even on failure.
func TestRun_TurnFailureShowsUserFacingMessageNotJustStderr(t *testing.T) {
	cfg := testConfig("http://127.0.0.1:1") // port 1 is reserved/unlisted, refuses immediately
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check for a single-turn test

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hello\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run should exit cleanly even when Ollama is unreachable, got error: %v", err)
	}
	if !strings.Contains(stdout.String(), turnFailedFallbackReply) {
		t.Errorf("stdout = %q, want it to include a user-facing message on failure, not just a stderr log", stdout.String())
	}
}

func TestRun_HandlesUnreachableHostGracefully(t *testing.T) {
	cfg := testConfig("http://127.0.0.1:1") // port 1 is reserved/unlisted, refuses immediately
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hello\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run should exit cleanly even when Ollama is unreachable, got error: %v", err)
	}
	if !strings.Contains(stderr.String(), "could not reach") {
		t.Errorf("stderr = %q, want a message about being unable to reach the host", stderr.String())
	}
}

// TestRun_ErrorMessageIncludesRealUnderlyingDetail is a regression test
// for a real bug found live: a model picked without real tool-calling
// support made Ollama reject every request with 400 "does not support
// tools" -- but the generic "could not reach <host>" message swallowed
// that detail entirely and read exactly like a network/Docker
// connectivity problem, sending a live debugging session down the wrong
// path. The real error must be visible, not just the host.
func TestRun_ErrorMessageIncludesRealUnderlyingDetail(t *testing.T) {
	// Ollama's own real error responses are JSON with an "error" field
	// (see eino-contrib/ollama/api's checkError), not a plain-text body
	// -- matching that shape here so the client actually decodes the
	// message instead of failing on JSON unmarshal first.
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"does not support tools"}`))
	}))
	defer mock.Close()

	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hello\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run should exit cleanly even on a request error, got error: %v", err)
	}
	got := stderr.String()
	if !strings.Contains(got, "could not reach") {
		t.Errorf("stderr = %q, want the generic message preserved", got)
	}
	if !strings.Contains(got, "does not support tools") {
		t.Errorf("stderr = %q, want the real underlying error detail included, not just the generic host message", got)
	}
}

// TestRun_ComprehensionCheckErrorMessageIncludesRealUnderlyingDetail is
// the comprehension-check path's counterpart to
// TestRun_ErrorMessageIncludesRealUnderlyingDetail -- runComprehensionCheck
// has its own separate "could not reach" call site with the same bug.
func TestRun_ComprehensionCheckErrorMessageIncludesRealUnderlyingDetail(t *testing.T) {
	// Ollama's own real error responses are JSON with an "error" field
	// (see eino-contrib/ollama/api's checkError), not a plain-text body
	// -- matching that shape here so the client actually decodes the
	// message instead of failing on JSON unmarshal first.
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"does not support tools"}`))
	}))
	defer mock.Close()

	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist // wants the comprehension check

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hi\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run should exit cleanly even on a request error, got error: %v", err)
	}
	got := stderr.String()
	if !strings.Contains(got, "could not reach") {
		t.Errorf("stderr = %q, want the generic message preserved", got)
	}
	if !strings.Contains(got, "does not support tools") {
		t.Errorf("stderr = %q, want the real underlying error detail included, not just the generic host message", got)
	}
}

// TestRun_OpenRouterModelShowsOpenRouterNotEmptyHostInBannerAndErrors is a
// real bug found live (via a real OpenRouter session): the startup banner
// and both "could not reach" error sites print cfg.OllamaHost directly,
// which is meaningless -- empty, in practice -- for an
// OpenRouterModelPrefix-prefixed model, showing "connected to ." and
// "could not reach :" instead of naming the actual provider.
func TestRun_OpenRouterModelShowsOpenRouterNotEmptyHostInBannerAndErrors(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"rate limited"}}`))
	}))
	defer mock.Close()

	origBaseURL := openRouterBaseURL
	openRouterBaseURL = mock.URL
	defer func() { openRouterBaseURL = origBaseURL }()

	cfg := testConfig("") // OllamaHost deliberately empty/unused for this path
	cfg.Model = OpenRouterModelPrefix + "some/model"
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hello\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run should exit cleanly even on a request error, got error: %v", err)
	}

	banner := stdout.String()
	if !strings.Contains(banner, "connected to OpenRouter") {
		t.Errorf("stdout banner = %q, want it to say \"connected to OpenRouter\"", banner)
	}

	errOut := stderr.String()
	if !strings.Contains(errOut, "could not reach OpenRouter:") {
		t.Errorf("stderr = %q, want \"could not reach OpenRouter:\", not the empty/meaningless OllamaHost", errOut)
	}
	if !strings.Contains(errOut, "rate limited") {
		t.Errorf("stderr = %q, want the real underlying error detail included too", errOut)
	}
}

func TestRun_ComprehensionCheckIncludesUsersRealFirstMessage(t *testing.T) {
	// A real bug found live: an earlier version of runComprehensionCheck
	// deliberately excluded the user's actual first message from the
	// request, so literally any first message -- including a plain
	// "hi" -- got back the exact same canned restate-and-ask-questions
	// reply with no acknowledgment of what the user said.
	mock := newSequencedOllama(t, "hey! restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist // wants the comprehension check

	greeting := "hi"
	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader(greeting+"\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) == 0 {
		t.Fatal("expected at least 1 request (the comprehension check), got 0")
	}

	checkReq := reqs[0]
	found := false
	for _, m := range checkReq.Messages {
		if m.Content == greeting {
			found = true
		}
	}
	if !found {
		t.Errorf("comprehension check request never included the user's real first message %q: %+v", greeting, checkReq.Messages)
	}

	if !strings.Contains(stdout.String(), "hey! restated problem + questions") {
		t.Errorf("stdout = %q, want the comprehension check's reply printed", stdout.String())
	}
}

func TestRun_ComprehensionCheckRetriesWhenReplyLeaksFakeToolCallJSON(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_problem_statement", "parameters": {}}`, "clean restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hi\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	got := stdout.String()
	if strings.Contains(got, `{"name"`) {
		t.Errorf("stdout still contains leaked tool-call JSON: %q", got)
	}
	if !strings.Contains(got, "clean restated problem + questions") {
		t.Errorf("stdout = %q, want the retry's clean reply", got)
	}
	if n := len(mock.allRequests()); n != 2 {
		t.Errorf("requests = %d, want exactly 2 (original comprehension check + one retry)", n)
	}
}

func TestRun_ComprehensionCheckInjectsProblemStatementDirectly(t *testing.T) {
	mock := newSequencedOllama(t, "restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist
	cfg.WorkDir = t.TempDir()

	want := "# Two Sum\n\nReturn indices of the two numbers that add up to target."
	if err := os.WriteFile(filepath.Join(cfg.WorkDir, "problem.md"), []byte(want), 0o644); err != nil {
		t.Fatalf("write problem.md: %v", err)
	}

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("what's the problem?\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 request (the comprehension check), got %d", len(reqs))
	}

	found := false
	for _, m := range reqs[0].Messages {
		if strings.Contains(m.Content, want) {
			found = true
		}
	}
	if !found {
		t.Errorf("comprehension check request never included the injected problem statement %q: %+v", want, reqs[0].Messages)
	}
}

func TestRun_ComprehensionCheckSkippedForSyntaxOnly(t *testing.T) {
	mock := newSequencedOllama(t, "direct reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("what's wrong with my code\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected exactly 1 request (no comprehension check for syntax-only), got %d", len(reqs))
	}
	if !strings.Contains(stdout.String(), "direct reply") {
		t.Errorf("stdout = %q, want the direct reply printed with no check", stdout.String())
	}
}

func TestRun_ComprehensionCheckHistoryPersistsBothTurns(t *testing.T) {
	mock := newSequencedOllama(t, "restated + questions", "real answer")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeHintsFirst

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("real question\nfollow up\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	// Only 2 requests total: the check consumes the first input line's
	// turn entirely (no separate normal request is also sent for it),
	// then the second input line is one normal turn.
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requests (check, then 1 real turn), got %d", len(reqs))
	}

	// Second request's history should include the check's exchange —
	// persisted using the real first question as the user turn, per
	// runComprehensionCheck's doc comment — followed by the ephemeral
	// hint-count note (hints-first mode, see turnMessages) and the
	// second line.
	second := reqs[1].Messages
	if len(second) != 5 {
		t.Fatalf("second request: expected [system, user1, assistant1, hint-note, user2] = 5 messages, got %d: %+v", len(second), second)
	}
	if second[1].Content != "real question" {
		t.Errorf("second request messages[1] = %q, want the real first question %q", second[1].Content, "real question")
	}
	if second[2].Content != "restated + questions" {
		t.Errorf("second request messages[2] = %q, want the check's reply", second[2].Content)
	}
	if second[3].Role != "system" || !strings.Contains(second[3].Content, "1st help request") {
		t.Errorf("second request messages[3] = %+v, want an ephemeral system note about the 1st help request", second[3])
	}
	if second[4].Content != "follow up" {
		t.Errorf("second request messages[4] = %q, want %q", second[4].Content, "follow up")
	}
}

func TestTurnMessages_NonHintsFirstModeHasNoNote(t *testing.T) {
	for _, mode := range []string{exercise.TutorModeSyntaxOnly, exercise.TutorModeFullAssist} {
		msgs := turnMessages(mode, 1, "hello")
		if len(msgs) != 1 {
			t.Errorf("mode %s: turnMessages returned %d messages, want 1 (just the user message)", mode, len(msgs))
			continue
		}
		if msgs[0].Content != "hello" {
			t.Errorf("mode %s: messages[0].Content = %q, want %q", mode, msgs[0].Content, "hello")
		}
	}
}

func TestTurnMessages_HintsFirstFirstRequestNotesFirstAsk(t *testing.T) {
	msgs := turnMessages(exercise.TutorModeHintsFirst, 1, "help")
	if len(msgs) != 2 {
		t.Fatalf("turnMessages returned %d messages, want 2 (note + user message)", len(msgs))
	}
	if msgs[0].Role != "system" || !strings.Contains(msgs[0].Content, "1st help request") {
		t.Errorf("messages[0] = %+v, want a system note mentioning the 1st help request", msgs[0])
	}
	if msgs[1].Content != "help" {
		t.Errorf("messages[1].Content = %q, want %q", msgs[1].Content, "help")
	}
}

func TestTurnMessages_HintsFirstRepeatRequestNotesCount(t *testing.T) {
	msgs := turnMessages(exercise.TutorModeHintsFirst, 3, "still stuck")
	if len(msgs) != 2 {
		t.Fatalf("turnMessages returned %d messages, want 2 (note + user message)", len(msgs))
	}
	if msgs[0].Role != "system" || !strings.Contains(msgs[0].Content, "#3") {
		t.Errorf("messages[0] = %+v, want a system note mentioning request #3", msgs[0])
	}
	if !strings.Contains(msgs[0].Content, "don't ask them to confirm") {
		t.Errorf("messages[0].Content = %q, want it to explicitly tell the model not to ask the user to confirm", msgs[0].Content)
	}
}

// --- decideHandoff (orchestrator/worker routing) ---

func TestDecideHandoff_YesRepliesMeanHandoff(t *testing.T) {
	mock := newSequencedOllama(t, "YES")
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "how do I solve this with a hash map")
	if err != nil {
		t.Fatalf("decideHandoff: %v", err)
	}
	if !handoff {
		t.Error("handoff = false, want true for a YES reply")
	}
}

func TestDecideHandoff_NoRepliesMeanNoHandoff(t *testing.T) {
	mock := newSequencedOllama(t, "no")
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "hi")
	if err != nil {
		t.Fatalf("decideHandoff: %v", err)
	}
	if handoff {
		t.Error("handoff = true, want false for a (lowercase) no reply")
	}
}

func TestDecideHandoff_NoWithTrailingTextStillMeansNoHandoff(t *testing.T) {
	mock := newSequencedOllama(t, "No, this is just a greeting.")
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "hi there")
	if err != nil {
		t.Fatalf("decideHandoff: %v", err)
	}
	if handoff {
		t.Error("handoff = true, want false -- reply starts with No even with trailing explanation")
	}
}

func TestDecideHandoff_AmbiguousReplyDefaultsToHandoff(t *testing.T) {
	mock := newSequencedOllama(t, "I'm not sure, could go either way")
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "something ambiguous")
	if err != nil {
		t.Fatalf("decideHandoff: %v", err)
	}
	if !handoff {
		t.Error("handoff = false, want true -- anything not clearly starting with No defaults to handoff (safer to over-use the specialist)")
	}
}

func TestDecideHandoff_RequestErrorDefaultsToHandoffAndReturnsError(t *testing.T) {
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", "http://127.0.0.1:1", "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "hi")
	if err == nil {
		t.Fatal("expected an error for an unreachable host")
	}
	if !handoff {
		t.Error("handoff = false, want true -- a routing failure must default to handoff, not silently leave the turn unanswered")
	}
}

// --- Run() orchestrator/worker routing (single mock server; the
// worker and orchestrator model NAMES differ, so which "role" made a
// given request is verified via sequencedOllama's own recorded
// req.Model, not separate servers) ---

func TestRun_NoRoutingWhenOrchestratorModelEmpty(t *testing.T) {
	// OrchestratorModel unset (the zero value, matching every config
	// this package built before routing existed) must never make an
	// extra routing-decision call -- exactly one request per turn, all
	// against cfg.Model, identical to pre-routing behavior.
	mock := newSequencedOllama(t, "the only reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hi\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected exactly 1 request (no routing decision), got %d: %+v", len(reqs), reqs)
	}
	if reqs[0].Model != "test-model" {
		t.Errorf("request model = %q, want %q", reqs[0].Model, "test-model")
	}
}

func TestRun_RoutesToOrchestratorWhenDecisionIsNo(t *testing.T) {
	mock := newSequencedOllama(t, "NO", "orchestrator answered")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hi\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(stripAnsi(stdout.String()), "orchestrator answered") {
		t.Errorf("stdout = %q, want the orchestrator's reply", stdout.String())
	}
	reqs := mock.allRequests()
	if len(reqs) != 2 {
		t.Fatalf("expected exactly 2 requests (routing decision + orchestrator answer), got %d: %+v", len(reqs), reqs)
	}
	for i, req := range reqs {
		if req.Model != "orchestrator-model" {
			t.Errorf("request[%d].Model = %q, want %q -- worker must never be touched when the decision is No", i, req.Model, "orchestrator-model")
		}
	}
}

func TestRun_RoutesToWorkerWhenDecisionIsYes(t *testing.T) {
	mock := newSequencedOllama(t, "YES", "worker answered")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("how do I solve this\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(stripAnsi(stdout.String()), "worker answered") {
		t.Errorf("stdout = %q, want the worker's reply", stdout.String())
	}
	reqs := mock.allRequests()
	if len(reqs) != 2 {
		t.Fatalf("expected exactly 2 requests (routing decision + worker answer), got %d: %+v", len(reqs), reqs)
	}
	if reqs[0].Model != "orchestrator-model" {
		t.Errorf("request[0].Model = %q, want %q -- the routing decision itself always goes to the orchestrator", reqs[0].Model, "orchestrator-model")
	}
	if reqs[1].Model != "worker-model" {
		t.Errorf("request[1].Model = %q, want %q -- a Yes decision must hand off to the worker", reqs[1].Model, "worker-model")
	}
}

func TestRun_ComprehensionCheckAlwaysUsesOrchestratorWhenRoutingEnabled(t *testing.T) {
	// The comprehension check is a fixed, generic, single-purpose task
	// (restate + ask questions) -- it never goes through decideHandoff,
	// regardless of what the user's first message looks like.
	mock := newSequencedOllama(t, "restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeFullAssist // wants the comprehension check

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hi\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected exactly 1 request (the comprehension check, no routing decision), got %d: %+v", len(reqs), reqs)
	}
	if reqs[0].Model != "orchestrator-model" {
		t.Errorf("request[0].Model = %q, want %q -- the comprehension check always uses the orchestrator", reqs[0].Model, "orchestrator-model")
	}
}
