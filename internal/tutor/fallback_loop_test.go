package tutor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// fakeToolCallingChatModel is a minimal in-process model.ToolCallingChatModel
// for runFallbackToolLoop's unit tests -- faster and more precise for
// round-by-round control flow (round cap, dedup, malformed retries) than
// routing every case through an httptest.Server mock. Stream/WithTools
// are never called by runFallbackToolLoop (it only calls Generate
// directly, no WithTools binding), so they're minimal stubs.
type fakeToolCallingChatModel struct {
	mu      sync.Mutex
	calls   int
	replyFn func(round int) (string, error)
}

func (f *fakeToolCallingChatModel) Generate(_ context.Context, _ []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	f.mu.Lock()
	round := f.calls
	f.calls++
	f.mu.Unlock()

	content, err := f.replyFn(round)
	if err != nil {
		return nil, err
	}
	return schema.AssistantMessage(content, nil), nil
}

func (f *fakeToolCallingChatModel) Stream(context.Context, []*schema.Message, ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return nil, fmt.Errorf("fakeToolCallingChatModel: Stream not implemented")
}

func (f *fakeToolCallingChatModel) WithTools([]*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	return f, nil
}

func (f *fakeToolCallingChatModel) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

// sequencedFakeModel replies in order, clamping at the last reply for
// any extra call -- same convention as tutor_test.go's sequencedOllama.
func sequencedFakeModel(replies ...string) *fakeToolCallingChatModel {
	return &fakeToolCallingChatModel{replyFn: func(round int) (string, error) {
		if round >= len(replies) {
			round = len(replies) - 1
		}
		return replies[round], nil
	}}
}

type echoToolInput struct {
	Text string `json:"text"`
}
type echoToolOutput struct {
	Text string `json:"text"`
}

func newEchoTool(t *testing.T) tool.InvokableTool {
	t.Helper()
	callCount := 0
	tl, err := utils.InferTool("echo", "Echoes back its input text.", func(_ context.Context, in echoToolInput) (echoToolOutput, error) {
		callCount++
		return echoToolOutput{Text: in.Text}, nil
	})
	if err != nil {
		t.Fatalf("InferTool: %v", err)
	}
	return tl
}

// newCountingEchoTool is newEchoTool but exposes how many times it was
// actually invoked, for tests asserting a guard prevented re-invocation.
func newCountingEchoTool(t *testing.T) (tool.InvokableTool, *int) {
	t.Helper()
	count := 0
	tl, err := utils.InferTool("echo", "Echoes back its input text.", func(_ context.Context, in echoToolInput) (echoToolOutput, error) {
		count++
		return echoToolOutput{Text: in.Text}, nil
	})
	if err != nil {
		t.Fatalf("InferTool: %v", err)
	}
	return tl, &count
}

func newPanicTool(t *testing.T) tool.InvokableTool {
	t.Helper()
	tl, err := utils.InferTool("boom", "Always panics.", func(_ context.Context, _ echoToolInput) (echoToolOutput, error) {
		panic("boom: simulated tool panic")
	})
	if err != nil {
		t.Fatalf("InferTool: %v", err)
	}
	return tl
}

// --- findTool ---

func TestFindTool_FindsMatchingToolByName(t *testing.T) {
	echo := newEchoTool(t)
	got := findTool(context.Background(), []tool.BaseTool{echo}, "echo")
	if got == nil {
		t.Fatal("findTool returned nil, want the echo tool")
	}
}

func TestFindTool_ReturnsNilWhenNotFound(t *testing.T) {
	echo := newEchoTool(t)
	got := findTool(context.Background(), []tool.BaseTool{echo}, "does_not_exist")
	if got != nil {
		t.Error("findTool returned a non-nil tool for an unknown name")
	}
}

// --- safeInvokeTool ---

func TestSafeInvokeTool_ReturnsNormalResult(t *testing.T) {
	echo := newEchoTool(t)
	out, err := safeInvokeTool(context.Background(), echo, `{"text": "hello"}`)
	if err != nil {
		t.Fatalf("safeInvokeTool: %v", err)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("result = %q, want it to contain %q", out, "hello")
	}
}

func TestSafeInvokeTool_RecoversPanicIntoError(t *testing.T) {
	boom := newPanicTool(t)
	_, err := safeInvokeTool(context.Background(), boom, `{"text": "x"}`)
	if err == nil {
		t.Fatal("expected an error recovered from the tool's panic, got nil")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("error = %v, want it to mention the panic value", err)
	}
}

// --- pushActivity ---

func TestPushActivity_SendsCurrentCalls(t *testing.T) {
	feed := &activityFeed{}
	feed.started("id1", "echo", "")
	ch := make(chan []activityCall, 1)

	pushActivity(feed, ch)

	select {
	case calls := <-ch:
		if len(calls) != 1 || calls[0].callID != "id1" {
			t.Errorf("pushed calls = %+v, want one call with callID id1", calls)
		}
	default:
		t.Fatal("nothing was pushed onto the channel")
	}
}

func TestPushActivity_DropsInsteadOfBlockingWhenChannelFull(t *testing.T) {
	feed := &activityFeed{}
	feed.started("id1", "echo", "")
	ch := make(chan []activityCall, 1)
	ch <- []activityCall{} // fill the buffer

	done := make(chan struct{})
	go func() {
		pushActivity(feed, ch)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("pushActivity blocked instead of dropping when the channel was full")
	}
}

// --- renderToolCatalog ---

func TestRenderToolCatalog_IncludesNameAndDescription(t *testing.T) {
	echo := newEchoTool(t)
	got, err := renderToolCatalog(context.Background(), []tool.BaseTool{echo})
	if err != nil {
		t.Fatalf("renderToolCatalog: %v", err)
	}
	if !strings.Contains(got, "echo") {
		t.Errorf("catalog = %q, want it to mention the tool name", got)
	}
	if !strings.Contains(got, "Echoes back its input text.") {
		t.Errorf("catalog = %q, want it to mention the tool description", got)
	}
}

func TestRenderToolCatalog_ToolWithArgumentsIncludesSchema(t *testing.T) {
	echo := newEchoTool(t)
	got, err := renderToolCatalog(context.Background(), []tool.BaseTool{echo})
	if err != nil {
		t.Fatalf("renderToolCatalog: %v", err)
	}
	if !strings.Contains(got, "text") {
		t.Errorf("catalog = %q, want it to mention the \"text\" argument from echoToolInput's schema", got)
	}
}

// --- runFallbackToolLoop ---

func TestRunFallbackToolLoop_ImmediateFinalAnswerReturnsIt(t *testing.T) {
	cm := sequencedFakeModel("The answer is 42.")
	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)

	reply, err := runFallbackToolLoop(context.Background(), cm, nil, []*schema.Message{schema.UserMessage("hi")}, feed, ch)
	if err != nil {
		t.Fatalf("runFallbackToolLoop: %v", err)
	}
	if reply.Content != "The answer is 42." {
		t.Errorf("reply = %q, want the model's plain final answer", reply.Content)
	}
	if cm.callCount() != 1 {
		t.Errorf("callCount = %d, want 1 -- no tool call means no second round", cm.callCount())
	}
}

func TestRunFallbackToolLoop_ExecutesToolThenReturnsFinalAnswer(t *testing.T) {
	echo := newEchoTool(t)
	cm := sequencedFakeModel(
		`{"name": "echo", "arguments": {"text": "hello"}}`,
		"Done, the tool said hello.",
	)
	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)

	reply, err := runFallbackToolLoop(context.Background(), cm, []tool.BaseTool{echo}, []*schema.Message{schema.UserMessage("hi")}, feed, ch)
	if err != nil {
		t.Fatalf("runFallbackToolLoop: %v", err)
	}
	if reply.Content != "Done, the tool said hello." {
		t.Errorf("reply = %q, want the final answer", reply.Content)
	}
	if cm.callCount() != 2 {
		t.Errorf("callCount = %d, want 2 (tool-call round + final-answer round)", cm.callCount())
	}
}

func TestRunFallbackToolLoop_PushesActivityForExecutedTool(t *testing.T) {
	echo := newEchoTool(t)
	cm := sequencedFakeModel(
		`{"name": "echo", "arguments": {"text": "hello"}}`,
		"final",
	)
	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)

	if _, err := runFallbackToolLoop(context.Background(), cm, []tool.BaseTool{echo}, []*schema.Message{schema.UserMessage("hi")}, feed, ch); err != nil {
		t.Fatalf("runFallbackToolLoop: %v", err)
	}

	var lastCalls []activityCall
	for {
		select {
		case calls := <-ch:
			lastCalls = calls
			continue
		default:
		}
		break
	}
	if len(lastCalls) != 1 {
		t.Fatalf("got %d activity calls, want 1", len(lastCalls))
	}
	if lastCalls[0].name != "echo" || lastCalls[0].status != "done" {
		t.Errorf("activity call = %+v, want name=echo status=done", lastCalls[0])
	}
}

func TestRunFallbackToolLoop_UnknownToolNameRetriesWithCorrection(t *testing.T) {
	echo := newEchoTool(t)
	cm := sequencedFakeModel(
		`{"name": "does_not_exist", "arguments": {}}`,
		"Okay, here's my real answer.",
	)
	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)

	reply, err := runFallbackToolLoop(context.Background(), cm, []tool.BaseTool{echo}, []*schema.Message{schema.UserMessage("hi")}, feed, ch)
	if err != nil {
		t.Fatalf("runFallbackToolLoop: %v", err)
	}
	if reply.Content != "Okay, here's my real answer." {
		t.Errorf("reply = %q, want the eventual final answer", reply.Content)
	}
	if cm.callCount() != 2 {
		t.Errorf("callCount = %d, want 2 (failed attempt + retry)", cm.callCount())
	}
}

func TestRunFallbackToolLoop_MalformedCallRetriesWithCorrection(t *testing.T) {
	echo := newEchoTool(t)
	cm := sequencedFakeModel(
		`{"name": "echo", "arguments": {}`, // missing closing brace
		"Fixed, here's my real answer.",
	)
	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)

	reply, err := runFallbackToolLoop(context.Background(), cm, []tool.BaseTool{echo}, []*schema.Message{schema.UserMessage("hi")}, feed, ch)
	if err != nil {
		t.Fatalf("runFallbackToolLoop: %v", err)
	}
	if reply.Content != "Fixed, here's my real answer." {
		t.Errorf("reply = %q, want the eventual final answer", reply.Content)
	}
}

func TestRunFallbackToolLoop_RepeatedIdenticalCallIsGuardedNotReexecuted(t *testing.T) {
	echo, callCount := newCountingEchoTool(t)
	cm := sequencedFakeModel(
		`{"name": "echo", "arguments": {"text": "same"}}`,
		`{"name": "echo", "arguments": {"text": "same"}}`,
		"Okay, final answer now.",
	)
	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)

	reply, err := runFallbackToolLoop(context.Background(), cm, []tool.BaseTool{echo}, []*schema.Message{schema.UserMessage("hi")}, feed, ch)
	if err != nil {
		t.Fatalf("runFallbackToolLoop: %v", err)
	}
	if reply.Content != "Okay, final answer now." {
		t.Errorf("reply = %q, want the final answer", reply.Content)
	}
	if *callCount != 1 {
		t.Errorf("tool was invoked %d times, want 1 -- the identical repeat must be guarded, not re-executed", *callCount)
	}
	if cm.callCount() != 3 {
		t.Errorf("model callCount = %d, want 3 (call + guarded repeat + final)", cm.callCount())
	}
}

func TestRunFallbackToolLoop_RoundCapExhaustedReturnsFallbackReply(t *testing.T) {
	echo := newEchoTool(t)
	// A different tool call every round -- never triggers the repeat
	// guard, never gives a final answer, so the loop must exhaust
	// fallbackRoundCap on its own.
	cm := &fakeToolCallingChatModel{replyFn: func(round int) (string, error) {
		return fmt.Sprintf(`{"name": "echo", "arguments": {"text": "%d"}}`, round), nil
	}}
	feed := &activityFeed{}
	ch := make(chan []activityCall, 64)

	reply, err := runFallbackToolLoop(context.Background(), cm, []tool.BaseTool{echo}, []*schema.Message{schema.UserMessage("hi")}, feed, ch)
	if err != nil {
		t.Fatalf("runFallbackToolLoop: %v", err)
	}
	if reply.Content != fallbackLoopExhaustedReply {
		t.Errorf("reply = %q, want the exhausted-loop fallback message", reply.Content)
	}
	if cm.callCount() != fallbackRoundCap {
		t.Errorf("callCount = %d, want exactly fallbackRoundCap (%d)", cm.callCount(), fallbackRoundCap)
	}
}

func TestRunFallbackToolLoop_TransportErrorPropagates(t *testing.T) {
	wantErr := fmt.Errorf("simulated transport failure")
	cm := &fakeToolCallingChatModel{replyFn: func(int) (string, error) {
		return "", wantErr
	}}
	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)

	_, err := runFallbackToolLoop(context.Background(), cm, nil, []*schema.Message{schema.UserMessage("hi")}, feed, ch)
	if err == nil {
		t.Fatal("expected the transport error to propagate, got nil")
	}
}

func TestRunFallbackToolLoop_ToolPanicIsRecoveredAndFedBackAsError(t *testing.T) {
	boom := newPanicTool(t)
	cm := sequencedFakeModel(
		`{"name": "boom", "arguments": {}}`,
		"Recovered, here's my real answer.",
	)
	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)

	reply, err := runFallbackToolLoop(context.Background(), cm, []tool.BaseTool{boom}, []*schema.Message{schema.UserMessage("hi")}, feed, ch)
	if err != nil {
		t.Fatalf("runFallbackToolLoop returned an error instead of recovering the tool panic: %v", err)
	}
	if reply.Content != "Recovered, here's my real answer." {
		t.Errorf("reply = %q, want the loop to have continued past the panic", reply.Content)
	}
}

// TestRunFallbackToolLoop_RealToolsAndMockOllamaRoundTrip uses the real
// buildTools(cfg) + newChatModel against sequencedOllama (tutor_test.go's
// existing mock, unchanged -- it already never populates tool_calls,
// exactly matching a real non-tool-calling model), catching any drift
// between the fallback protocol and this package's actual production
// tools that a fake-tool unit test couldn't.
func TestRunFallbackToolLoop_RealToolsAndMockOllamaRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "solution.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write solution file: %v", err)
	}
	cfg := Config{WorkDir: dir, MaxContextBytes: 8000}
	tools, err := buildTools(cfg)
	if err != nil {
		t.Fatalf("buildTools: %v", err)
	}

	mock := newSequencedOllama(t,
		`{"name": "read_solution_file", "arguments": {}}`,
		"Your file looks fine.",
	)
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	feed := &activityFeed{}
	ch := make(chan []activityCall, 32)
	reply, err := runFallbackToolLoop(ctx, cm, tools, []*schema.Message{schema.UserMessage("is my code okay?")}, feed, ch)
	if err != nil {
		t.Fatalf("runFallbackToolLoop: %v", err)
	}
	if reply.Content != "Your file looks fine." {
		t.Errorf("reply = %q, want the final answer", reply.Content)
	}

	reqs := mock.allRequests()
	if len(reqs) != 2 {
		t.Fatalf("got %d requests to the mock, want 2", len(reqs))
	}
	second := reqs[1]
	found := false
	for _, m := range second.Messages {
		if strings.Contains(m.Content, "Tool result:") && strings.Contains(m.Content, "package main") {
			found = true
		}
	}
	if !found {
		t.Errorf("second request's messages = %+v, want one containing \"Tool result:\" and the real file content", second.Messages)
	}
}
