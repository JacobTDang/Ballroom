package tutor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestKittyAvailable_TrueWhenWindowIDSet(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "1")
	if !kittyAvailable() {
		t.Error("expected kittyAvailable() to be true when KITTY_WINDOW_ID is set")
	}
}

func TestKittyAvailable_FalseWhenWindowIDUnset(t *testing.T) {
	t.Setenv("KITTY_WINDOW_ID", "")
	if kittyAvailable() {
		t.Error("expected kittyAvailable() to be false when KITTY_WINDOW_ID is unset")
	}
}

func TestWrapForTmuxPassthrough_NoOpOutsideTmux(t *testing.T) {
	t.Setenv("TMUX", "")
	seq := "\033_Gaction=t\033\\"
	got := wrapForTmuxPassthrough(seq)
	if got != seq {
		t.Errorf("wrapForTmuxPassthrough outside tmux = %q, want unchanged %q", got, seq)
	}
}

func TestWrapForTmuxPassthrough_WrapsAndDoublesEscapesInsideTmux(t *testing.T) {
	t.Setenv("TMUX", "/tmp/tmux-1000/default,1234,0")
	seq := "\033_Gaction=t\033\\"
	got := wrapForTmuxPassthrough(seq)

	want := "\033Ptmux;" + "\033\033_Gaction=t\033\033\\" + "\033\\"
	if got != want {
		t.Errorf("wrapForTmuxPassthrough inside tmux =\n%q\nwant\n%q", got, want)
	}
}

func TestKittyTransmit_SmallPayloadIsOneChunk(t *testing.T) {
	t.Setenv("TMUX", "")
	data := []byte("not a real png, just testing chunking")

	seqs := kittyTransmit(7, data)
	if len(seqs) != 1 {
		t.Fatalf("len(seqs) = %d, want 1 for a payload under the chunk size", len(seqs))
	}
	seq := seqs[0]
	if !strings.HasPrefix(seq, "\033_G") || !strings.HasSuffix(seq, "\033\\") {
		t.Fatalf("seq = %q, want APC framing \\033_G...\\033\\\\", seq)
	}
	if !strings.Contains(seq, "a=t") || !strings.Contains(seq, "f=100") || !strings.Contains(seq, "i=7") {
		t.Errorf("seq = %q, want transmit action, PNG format, and image id 7", seq)
	}
	if !strings.Contains(seq, "m=0") {
		t.Errorf("seq = %q, want m=0 (final/only chunk) for a single-chunk payload", seq)
	}
	wantB64 := base64.StdEncoding.EncodeToString(data)
	if !strings.Contains(seq, wantB64) {
		t.Errorf("seq does not contain the expected base64 payload %q:\n%s", wantB64, seq)
	}
}

func TestKittyTransmit_LargePayloadChunksWithMoreFlag(t *testing.T) {
	t.Setenv("TMUX", "")
	data := make([]byte, kittyChunkSize*2) // base64 inflation alone forces >1 chunk
	for i := range data {
		data[i] = byte(i)
	}

	seqs := kittyTransmit(3, data)
	if len(seqs) < 2 {
		t.Fatalf("len(seqs) = %d, want >1 for a payload well over the chunk size", len(seqs))
	}
	for i, seq := range seqs {
		last := i == len(seqs)-1
		wantFlag := "m=1"
		if last {
			wantFlag = "m=0"
		}
		if !strings.Contains(seq, wantFlag) {
			t.Errorf("chunk %d/%d: seq = %q, want %s", i, len(seqs), seq, wantFlag)
		}
		if !last && strings.Contains(seq, "m=0") {
			t.Errorf("chunk %d/%d: non-final chunk incorrectly contains m=0: %q", i, len(seqs), seq)
		}
	}
	// Only the first chunk needs the full control data — repeating it
	// wastes bytes and isn't required by the protocol.
	if !strings.Contains(seqs[0], "f=100") {
		t.Errorf("first chunk = %q, want format/action control data", seqs[0])
	}
	if strings.Contains(seqs[1], "f=100") {
		t.Errorf("continuation chunk = %q, want only m= and q= control data", seqs[1])
	}
}

func TestKittyTransmit_WrapsEachChunkForTmuxPassthrough(t *testing.T) {
	t.Setenv("TMUX", "/tmp/tmux-1000/default,1234,0")
	seqs := kittyTransmit(1, []byte("x"))
	for _, seq := range seqs {
		if !strings.HasPrefix(seq, "\033Ptmux;") {
			t.Errorf("seq = %q, want tmux passthrough framing since $TMUX is set", seq)
		}
	}
}

func TestKittyShow_ReferencesImageIDAndCellSize(t *testing.T) {
	t.Setenv("TMUX", "")
	seq := kittyShow(42, 8, 4)
	for _, want := range []string{"a=p", "i=42", "c=8", "r=4", "C=1"} {
		if !strings.Contains(seq, want) {
			t.Errorf("kittyShow(42, 8, 4) = %q, want it to contain %q", seq, want)
		}
	}
}

func TestKittyDelete_ReferencesImageID(t *testing.T) {
	t.Setenv("TMUX", "")
	seq := kittyDelete(9)
	if !strings.Contains(seq, "a=d") || !strings.Contains(seq, "i=9") {
		t.Errorf("kittyDelete(9) = %q, want delete action and image id 9", seq)
	}
}

func TestKittyBallRenderer_InitTransmitsAllFrames(t *testing.T) {
	t.Setenv("TMUX", "")
	var buf bytes.Buffer
	kittyBallRenderer{}.init(&buf)

	got := buf.String()
	if strings.Count(got, "a=t") == 0 {
		t.Fatal("expected at least one transmit command from init, got none")
	}
	firstID := fmt.Sprintf("i=%d", kittyBallBaseImageID)
	lastID := fmt.Sprintf("i=%d", kittyBallBaseImageID+discoBallFrameCount-1)
	if !strings.Contains(got, firstID) {
		t.Errorf("expected the first frame's image id (%s) in the transmitted output", firstID)
	}
	if !strings.Contains(got, lastID) {
		t.Errorf("expected the last frame's image id (%s) in the transmitted output", lastID)
	}
}

func TestKittyBallRenderer_ShowFrameReferencesTheCorrectImageID(t *testing.T) {
	t.Setenv("TMUX", "")
	var buf bytes.Buffer
	kittyBallRenderer{}.showFrame(&buf, 3)

	want := fmt.Sprintf("i=%d", kittyBallBaseImageID+3)
	if !strings.Contains(buf.String(), want) {
		t.Errorf("showFrame(_, 3) = %q, want it to contain %q", buf.String(), want)
	}
}

func TestKittyBallRenderer_ShowFrameWrapsFrameIndex(t *testing.T) {
	t.Setenv("TMUX", "")
	var buf bytes.Buffer
	kittyBallRenderer{}.showFrame(&buf, discoBallFrameCount) // one past the last real frame

	want := fmt.Sprintf("i=%d", kittyBallBaseImageID) // wraps back to frame 0's id
	if !strings.Contains(buf.String(), want) {
		t.Errorf("showFrame(_, discoBallFrameCount) = %q, want it to wrap to %q", buf.String(), want)
	}
}

func TestKittyBallRenderer_ShowFrameAdvancesCursorByExactlyRows(t *testing.T) {
	// kittyShow uses C=1 (asks Kitty not to move the cursor after
	// placement) — showFrame must advance it manually via its own
	// newlines so its contract (cursor ends up exactly rows() lines
	// below where it started) holds. This was missing entirely until
	// this pass; redrawLocked's cursor-up math silently assumed it
	// without the display ever actually relying on it (masked first by
	// the fixed-row-count-that-happened-to-be-unused-by-DECSC design,
	// then by DECSC/DECRC itself not needing this contract at all).
	t.Setenv("TMUX", "")
	var buf bytes.Buffer
	r := kittyBallRenderer{}
	r.showFrame(&buf, 0)

	if got := strings.Count(buf.String(), "\n"); got != r.rows() {
		t.Errorf("showFrame wrote %d newlines, rows() reports %d — must match for redrawLocked's cursor-up math to stay correct", got, r.rows())
	}
}

func TestKittyBallRenderer_CloseDeletesEveryFrame(t *testing.T) {
	t.Setenv("TMUX", "")
	var buf bytes.Buffer
	kittyBallRenderer{}.close(&buf)

	got := strings.Count(buf.String(), "a=d")
	if got != discoBallFrameCount {
		t.Errorf("close() emitted %d delete commands, want %d", got, discoBallFrameCount)
	}
}
