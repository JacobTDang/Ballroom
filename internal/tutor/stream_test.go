package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// sseReply scripts one request's response from sseOpenRouter. chunks
// are streamed one SSE data event apiece for a stream=true request, or
// joined into a single completion body for a plain one. empty simulates
// a rate-limited free-tier reply: an SSE carrying no content chunks at
// all (stream) or a 200 with zero choices (plain). toolCall, when set,
// is emitted after the content -- as its own delta chunk on the stream
// path, or in the completion's tool_calls field on the plain path --
// reproducing a model that narrates text before calling a tool.
type sseReply struct {
	chunks   []string
	empty    bool
	toolCall *sseToolCall
}

type sseToolCall struct {
	name string
	args string
}

// sseOpenRouter is an OpenRouter-shaped mock that, unlike
// agent_test.go's single-completion mock, understands the request's own
// "stream" flag -- eino sends stream=true for Agent.Stream calls and
// stream=false for Generate, and streamWithLeakGuard's whole contract
// (progressive display, non-streaming leak retry) spans both kinds in
// one conversation.
type sseOpenRouter struct {
	*httptest.Server
	t *testing.T

	mu             sync.Mutex
	replies        []sseReply
	requests       int
	streamRequests int
}

func newSSEOpenRouter(t *testing.T, replies ...sseReply) *sseOpenRouter {
	t.Helper()
	m := &sseOpenRouter{t: t, replies: replies}
	m.Server = httptest.NewServer(http.HandlerFunc(m.handle))
	t.Cleanup(m.Server.Close)

	orig := openRouterBaseURL
	openRouterBaseURL = m.Server.URL
	t.Cleanup(func() { openRouterBaseURL = orig })
	return m
}

func (m *sseOpenRouter) counts() (requests, streamRequests int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.requests, m.streamRequests
}

func (m *sseOpenRouter) handle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Stream bool `json:"stream"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	m.requests++
	if req.Stream {
		m.streamRequests++
	}
	if len(m.replies) == 0 {
		m.mu.Unlock()
		m.t.Errorf("sseOpenRouter: unscripted request (stream=%v)", req.Stream)
		http.Error(w, "unscripted request", http.StatusInternalServerError)
		return
	}
	reply := m.replies[0]
	m.replies = m.replies[1:]
	m.mu.Unlock()

	if req.Stream {
		w.Header().Set("Content-Type", "text/event-stream")
		if !reply.empty {
			for i, c := range reply.chunks {
				delta := map[string]any{"content": c}
				if i == 0 {
					delta["role"] = "assistant"
				}
				chunk := map[string]any{
					"id": "chunk-test", "object": "chat.completion.chunk", "created": 1, "model": "test",
					"choices": []map[string]any{{"index": 0, "delta": delta}},
				}
				b, err := json.Marshal(chunk)
				if err != nil {
					m.t.Errorf("sseOpenRouter: marshal chunk: %v", err)
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", b)
			}
			if tc := reply.toolCall; tc != nil {
				chunk := map[string]any{
					"id": "chunk-test", "object": "chat.completion.chunk", "created": 1, "model": "test",
					"choices": []map[string]any{{"index": 0, "delta": map[string]any{
						"tool_calls": []map[string]any{{"index": 0, "id": "call-test", "type": "function", "function": map[string]string{"name": tc.name, "arguments": tc.args}}},
					}}},
				}
				b, err := json.Marshal(chunk)
				if err != nil {
					m.t.Errorf("sseOpenRouter: marshal tool-call chunk: %v", err)
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", b)
				fmt.Fprint(w, `data: {"id":"chunk-test","object":"chat.completion.chunk","created":1,"model":"test","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`+"\n\n")
			} else {
				fmt.Fprint(w, `data: {"id":"chunk-test","object":"chat.completion.chunk","created":1,"model":"test","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`+"\n\n")
			}
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if reply.empty {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "cmpl-test", "object": "chat.completion", "created": 1, "model": "test",
			"choices": []any{},
		})
		return
	}
	message := map[string]any{"role": "assistant", "content": strings.Join(reply.chunks, "")}
	finishReason := "stop"
	if tc := reply.toolCall; tc != nil {
		message["tool_calls"] = []map[string]any{{"id": "call-test", "type": "function", "function": map[string]string{"name": tc.name, "arguments": tc.args}}}
		finishReason = "tool_calls"
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id": "cmpl-test", "object": "chat.completion", "created": 1, "model": "test",
		"choices": []map[string]any{
			{"index": 0, "message": message, "finish_reason": finishReason},
		},
	})
}

// newStreamTestAgent builds a real react.Agent pointed at the mock via
// the OpenRouter path -- the only provider streaming is enabled for.
// tools are bound when given, so scripted tool-call replies actually
// execute.
func newStreamTestAgent(t *testing.T, tools ...tool.BaseTool) *react.Agent {
	t.Helper()
	ctx := context.Background()
	cm, err := newChatModel(ctx, OpenRouterModelPrefix+"test/model", "", "sk-test-key")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}
	agent, err := newAgent(ctx, cm, tools)
	if err != nil {
		t.Fatalf("newAgent: %v", err)
	}
	return agent
}

// newPingTool returns a trivial invokable tool and a pointer to its
// call count -- the same shape toolcheck.go's probe uses.
func newPingTool(t *testing.T) (tool.BaseTool, *int) {
	t.Helper()
	calls := 0
	pingTool, err := utils.InferTool("ping", "Call this tool to respond.", func(_ context.Context, _ pingArgs) (string, error) {
		calls++
		return "pong", nil
	})
	if err != nil {
		t.Fatalf("InferTool: %v", err)
	}
	return pingTool, &calls
}

func TestStreamingEnabled_DefaultsAndOverride(t *testing.T) {
	cases := []struct {
		name  string
		model string
		env   string
		want  bool
	}{
		{"openrouter defaults on", OpenRouterModelPrefix + "some/model", "", true},
		{"ollama defaults off (drops tool_calls when streaming)", "llama3.1:8b", "", false},
		{"override off wins over openrouter default", OpenRouterModelPrefix + "some/model", "off", false},
		{"override on wins over ollama default", "llama3.1:8b", "on", true},
		{"unrecognized override falls back to the default", "llama3.1:8b", "banana", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("TUTOR_STREAM", c.env)
			if got := streamingEnabled(c.model); got != c.want {
				t.Errorf("streamingEnabled(%q) with TUTOR_STREAM=%q = %v, want %v", c.model, c.env, got, c.want)
			}
		})
	}
}

func TestStreamWithLeakGuard_ProgressiveDisplayAndFullFinalReply(t *testing.T) {
	newSSEOpenRouter(t, sseReply{chunks: []string{"Sure — ", "use a hash ", "map here."}})
	agent := newStreamTestAgent(t)

	var seen []string
	reply, err := streamWithLeakGuard(context.Background(), agent, []*schema.Message{schema.UserMessage("how?")}, func(text string) {
		seen = append(seen, text)
	})
	if err != nil {
		t.Fatalf("streamWithLeakGuard: %v", err)
	}

	const want = "Sure — use a hash map here."
	if reply.Content != want {
		t.Errorf("reply.Content = %q, want %q", reply.Content, want)
	}
	if len(seen) < 2 {
		t.Fatalf("onText called %d times (%q), want at least 2 progressive updates", len(seen), seen)
	}
	if last := seen[len(seen)-1]; last != want {
		t.Errorf("last onText text = %q, want the full reply %q", last, want)
	}
	for i := 1; i < len(seen); i++ {
		if !strings.HasPrefix(seen[i], seen[i-1]) {
			t.Errorf("onText text %q does not extend the previous %q -- display must only ever grow", seen[i], seen[i-1])
		}
	}
}

func TestStreamWithLeakGuard_LeakedToolCallNeverDisplaysAndRetriesNonStreaming(t *testing.T) {
	mock := newSSEOpenRouter(t,
		// The stream leaks a fake tool call as text, split mid-name to
		// prove detection works on accumulated text, not per-chunk.
		sseReply{chunks: []string{`{"name": "read_`, `solution_file", "parameters": {}}`}},
		// The corrective retry (a plain Generate) answers cleanly.
		sseReply{chunks: []string{"Here's a grounded answer."}},
	)
	agent := newStreamTestAgent(t)

	var seen []string
	reply, err := streamWithLeakGuard(context.Background(), agent, []*schema.Message{schema.UserMessage("read my file")}, func(text string) {
		seen = append(seen, text)
	})
	if err != nil {
		t.Fatalf("streamWithLeakGuard: %v", err)
	}

	if len(seen) != 0 {
		t.Errorf("onText was called with %q -- leaked tool-call JSON must never paint, not even partially", seen)
	}
	if reply.Content != "Here's a grounded answer." {
		t.Errorf("reply.Content = %q, want the non-streaming retry's clean answer", reply.Content)
	}
	requests, streamRequests := mock.counts()
	if streamRequests != 1 || requests != 2 {
		t.Errorf("mock saw %d requests (%d streamed), want exactly 2 total with 1 streamed -- the leak retry must be non-streaming", requests, streamRequests)
	}
}

func TestStreamWithLeakGuard_BraceContentHeldBackThenDisplayed(t *testing.T) {
	// Content that legitimately starts with "{" (e.g. a JSON snippet the
	// tutor is quoting) must not paint while it's still short enough to
	// be the start of a leaked tool call -- but once it's clearly not
	// one, it must display.
	newSSEOpenRouter(t, sseReply{chunks: []string{`{"id": 1}`, ` is the record shape you want.`}})
	agent := newStreamTestAgent(t)

	var seen []string
	reply, err := streamWithLeakGuard(context.Background(), agent, []*schema.Message{schema.UserMessage("shape?")}, func(text string) {
		seen = append(seen, text)
	})
	if err != nil {
		t.Fatalf("streamWithLeakGuard: %v", err)
	}

	const want = `{"id": 1} is the record shape you want.`
	if reply.Content != want {
		t.Errorf("reply.Content = %q, want %q", reply.Content, want)
	}
	if len(seen) == 0 {
		t.Fatal("onText never called -- held-back brace content must still display once it's clearly not a leak")
	}
	if first := seen[0]; utf8.RuneCountInString(first) < streamHoldBackRunes {
		t.Errorf("first onText text %q is %d runes -- brace-prefixed text under %d runes must be held back from display", first, utf8.RuneCountInString(first), streamHoldBackRunes)
	}
	if last := seen[len(seen)-1]; last != want {
		t.Errorf("last onText text = %q, want the full reply %q", last, want)
	}
}

func TestStreamWithLeakGuard_EmptyStreamRetriesWithBackoff(t *testing.T) {
	origBackoff := emptyChoicesRetryBackoff
	emptyChoicesRetryBackoff = time.Millisecond
	defer func() { emptyChoicesRetryBackoff = origBackoff }()

	mock := newSSEOpenRouter(t,
		sseReply{empty: true},
		sseReply{chunks: []string{"Recovered ", "after the empty reply."}},
	)
	agent := newStreamTestAgent(t)

	var seen []string
	reply, err := streamWithLeakGuard(context.Background(), agent, []*schema.Message{schema.UserMessage("hi")}, func(text string) {
		seen = append(seen, text)
	})
	if err != nil {
		t.Fatalf("streamWithLeakGuard: %v", err)
	}

	const want = "Recovered after the empty reply."
	if reply.Content != want {
		t.Errorf("reply.Content = %q, want %q", reply.Content, want)
	}
	if len(seen) == 0 || seen[len(seen)-1] != want {
		t.Errorf("onText texts = %q, want the retried stream's text displayed", seen)
	}
	if _, streamRequests := mock.counts(); streamRequests != 2 {
		t.Errorf("mock saw %d streamed requests, want 2 -- an empty stream must be retried", streamRequests)
	}
}

func TestStreamWithLeakGuard_EmptyStreamsExhaustRetriesWithHonestError(t *testing.T) {
	origBackoff := emptyChoicesRetryBackoff
	emptyChoicesRetryBackoff = time.Millisecond
	defer func() { emptyChoicesRetryBackoff = origBackoff }()

	replies := make([]sseReply, emptyChoicesMaxRetries+1)
	for i := range replies {
		replies[i] = sseReply{empty: true}
	}
	newSSEOpenRouter(t, replies...)
	agent := newStreamTestAgent(t)

	_, err := streamWithLeakGuard(context.Background(), agent, []*schema.Message{schema.UserMessage("hi")}, func(string) {})
	if err == nil {
		t.Fatal("streamWithLeakGuard: want an error after every attempt streamed empty")
	}
	if !strings.Contains(err.Error(), "rate-limit") {
		t.Errorf("error %q, want it to explain the likely rate-limit cause", err)
	}
}

// TestStreamWithLeakGuard_NarrationBeforeToolCallStillExecutesTheTool
// pins the fix for a real bug found live (openrouter:poolside/laguna-xs
// answering "read my solution file" with just "I'll read your solution
// file to see what you're working on." -- and never reading anything):
// eino's default stream checker decides "no tool call" at the first
// content chunk, so a model that narrates text before emitting its
// tool_calls lost the call entirely. The windowed checker keeps reading
// past early narration.
func TestStreamWithLeakGuard_NarrationBeforeToolCallStillExecutesTheTool(t *testing.T) {
	mock := newSSEOpenRouter(t,
		// Round 1: narration first, then the tool call -- the shape that
		// defeated the first-chunk checker.
		sseReply{chunks: []string{"I'll read your ", "solution file now."}, toolCall: &sseToolCall{name: "ping", args: `{"reason": "check"}`}},
		// Round 2 (after the tool result): the grounded final answer.
		sseReply{chunks: []string{"Grounded answer ", "after the tool ran."}},
	)
	pingTool, calls := newPingTool(t)
	agent := newStreamTestAgent(t, pingTool)

	var seen []string
	reply, err := streamWithLeakGuard(context.Background(), agent, []*schema.Message{schema.UserMessage("read it")}, func(text string) {
		seen = append(seen, text)
	})
	if err != nil {
		t.Fatalf("streamWithLeakGuard: %v", err)
	}

	if *calls != 1 {
		t.Errorf("ping tool ran %d times, want exactly 1 -- narration before the tool_calls chunk must not swallow the call", *calls)
	}
	const want = "Grounded answer after the tool ran."
	if reply.Content != want {
		t.Errorf("reply.Content = %q, want the post-tool answer %q", reply.Content, want)
	}
	for _, text := range seen {
		if strings.Contains(text, "I'll read your") {
			t.Errorf("onText painted the pre-tool narration %q -- tool-round chunks must stay internal", text)
		}
	}
	if requests, streamRequests := mock.counts(); requests != 2 || streamRequests != 2 {
		t.Errorf("mock saw %d requests (%d streamed), want 2 streamed rounds", requests, streamRequests)
	}
}

// TestStreamWithLeakGuard_ToolCallPastTheDecisionWindowFallsBackToBlocking
// covers the windowed checker's own blind spot: narration longer than
// streamToolCallDecisionRunes before the tool call means the stream is
// routed to the caller as a final answer, tool never executed. The
// safety net must detect the unexecuted tool_calls on the reassembled
// reply, discard it, and redo the turn on the blocking Generate path
// (which handles tool calls correctly regardless of narration).
func TestStreamWithLeakGuard_ToolCallPastTheDecisionWindowFallsBackToBlocking(t *testing.T) {
	longNarration := strings.Repeat("Let me think about what to check here. ", 10) // well past the decision window
	mock := newSSEOpenRouter(t,
		// Round 1 (streamed): tool call hidden past the window.
		sseReply{chunks: []string{longNarration}, toolCall: &sseToolCall{name: "ping", args: `{"reason": "check"}`}},
		// The blocking redo: tool call answered properly...
		sseReply{chunks: []string{""}, toolCall: &sseToolCall{name: "ping", args: `{"reason": "check"}`}},
		// ...then the grounded final answer.
		sseReply{chunks: []string{"Grounded blocking answer."}},
	)
	pingTool, calls := newPingTool(t)
	agent := newStreamTestAgent(t, pingTool)

	reply, err := streamWithLeakGuard(context.Background(), agent, []*schema.Message{schema.UserMessage("read it")}, func(string) {})
	if err != nil {
		t.Fatalf("streamWithLeakGuard: %v", err)
	}

	if *calls != 1 {
		t.Errorf("ping tool ran %d times, want exactly 1 (on the blocking redo)", *calls)
	}
	if reply.Content != "Grounded blocking answer." {
		t.Errorf("reply.Content = %q, want the blocking redo's grounded answer", reply.Content)
	}
	if requests, streamRequests := mock.counts(); requests != 3 || streamRequests != 1 {
		t.Errorf("mock saw %d requests (%d streamed), want 3 total with only the defeated first round streamed", requests, streamRequests)
	}
}

// TestTutorModel_StreamingTurnPaintsProvisionalTextThenFinalReply is
// the model-level wiring test: a streamed turn must show the partial
// reply in the viewport while the turn is still in flight, then settle
// into exactly the same final state a non-streamed turn produces --
// provisional text gone, styled reply in displayLines, raw reply in
// history.
func TestTutorModel_StreamingTurnPaintsProvisionalTextThenFinalReply(t *testing.T) {
	newSSEOpenRouter(t, sseReply{chunks: []string{"First half, ", "second half."}})

	cfg := testConfig("")
	cfg.Model = OpenRouterModelPrefix + "test/model"
	cfg.APIKey = "sk-test-key"
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m = detectAndApply(t, m)

	m.textarea.SetValue("stream me")
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if cmd == nil {
		t.Fatal("submit produced no command")
	}

	sawProvisional := false
	for i := 0; ; i++ {
		if i > 200 {
			t.Fatal("turn never completed")
		}
		msg := cmd()
		newM, cmd = m.Update(msg)
		m = newM.(tutorModel)
		if _, ok := msg.(streamTextMsg); ok {
			if !m.turnInFlight {
				t.Error("streamTextMsg arrived with turnInFlight=false -- provisional text is a mid-turn state")
			}
			if strings.Contains(m.viewport.View(), "First half,") {
				sawProvisional = true
			}
		}
		if _, ok := msg.(turnCompleteMsg); ok {
			break
		}
		if cmd == nil {
			t.Fatal("turn ended without a turnCompleteMsg")
		}
	}

	if !sawProvisional {
		t.Error("the streamed partial reply never appeared in the viewport while the turn was in flight")
	}
	if m.streamingText != "" {
		t.Errorf("streamingText = %q after the turn completed, want it cleared", m.streamingText)
	}
	view := m.viewport.View()
	if n := strings.Count(view, "second half."); n != 1 {
		t.Errorf("final view shows the reply %d times, want exactly once (provisional block replaced by the final append):\n%s", n, view)
	}
	last := m.history[len(m.history)-1]
	if last.Role != schema.Assistant || last.Content != "First half, second half." {
		t.Errorf("history tail = %v %q, want the raw assistant reply", last.Role, last.Content)
	}
}

// TestTutorModel_OllamaTurnNeverStreams pins the Ollama-side contract:
// streaming stays off (Ollama drops tool_calls on streamed requests --
// the reason Generate was chosen originally), so a turn against an
// Ollama model must make only stream=false requests and deliver no
// streamTextMsg.
func TestTutorModel_OllamaTurnNeverStreams(t *testing.T) {
	mock := newSequencedOllama(t, "plain reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m = detectAndApply(t, m)

	m.textarea.SetValue("hello")
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if cmd == nil {
		t.Fatal("submit produced no command")
	}
	for i := 0; ; i++ {
		if i > 100 {
			t.Fatal("turn never completed")
		}
		msg := cmd()
		if _, ok := msg.(streamTextMsg); ok {
			t.Error("got a streamTextMsg for an Ollama model -- streaming must stay off for Ollama")
		}
		newM, cmd = m.Update(msg)
		m = newM.(tutorModel)
		if _, ok := msg.(turnCompleteMsg); ok {
			break
		}
		if cmd == nil {
			t.Fatal("turn ended without a turnCompleteMsg")
		}
	}

	for _, req := range mock.allRequests() {
		if req.Stream {
			t.Error("an Ollama request was sent with stream=true -- Ollama drops tool_calls when streaming")
		}
	}
}
