package tutor

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"
)

// kittyChunkSize caps how much base64 payload goes in a single escape
// sequence, per the Kitty graphics protocol's own documented chunked
// transmission scheme — avoids one huge escape sequence for a several-
// hundred-KB frame, which risks hitting terminal/pty buffer limits.
// Chunk boundaries don't need to land on base64 4-byte boundaries: the
// terminal concatenates every chunk's payload back into one string
// before decoding it as a whole.
const kittyChunkSize = 4096

// kittyAvailable reports whether the outer terminal (through however
// many passthrough layers) is Kitty, based on KITTY_WINDOW_ID being
// forwarded from the host via docker run -e (see
// internal/orchestrator/run.go's exerciseRunArgs/sandboxRunArgs).
// Ghostty also speaks the Kitty graphics protocol but doesn't set this
// var the same way — out of scope for now; unrecognized terminals just
// get the ANSI ball (internal/tutor/discoball.go), which is always
// correct, just less pretty.
func kittyAvailable() bool {
	return os.Getenv("KITTY_WINDOW_ID") != ""
}

// wrapForTmuxPassthrough wraps seq for tmux's allow-passthrough option
// (docker/tmux.conf) — tmux does not forward arbitrary escape sequences
// to the outer terminal by default (it owns its own screen-buffer
// model), so an application inside a tmux pane that wants to reach the
// real terminal directly (here: the Kitty graphics protocol, which tmux
// has no native understanding of) must wrap its sequence as a DCS
// "tmux;" passthrough command with every ESC byte in the original
// sequence doubled — this is tmux's own documented passthrough framing,
// not something specific to Kitty.
//
// A no-op when $TMUX isn't set (not running inside tmux at all), so
// this is always safe to call unconditionally rather than needing the
// caller to check first. In this codebase $TMUX is in practice always
// set for the tutor process — entrypoint.sh launches it as a command
// inside one of the container's own tmux panes — but the check keeps
// this function correct on its own terms rather than relying on that.
func wrapForTmuxPassthrough(seq string) string {
	if os.Getenv("TMUX") == "" {
		return seq
	}
	doubled := strings.ReplaceAll(seq, "\033", "\033\033")
	return "\033Ptmux;" + doubled + "\033\\"
}

// kittyTransmit returns the escape sequence(s) needed to upload pngData
// to Kitty under imageID, chunked at kittyChunkSize and each already
// wrapped via wrapForTmuxPassthrough. Sent as PNG bytes directly
// (f=100) — Kitty decodes PNG natively, so there's no image processing
// on this path at all, unlike the ANSI renderer's downsampling. Doesn't
// display anything; see kittyShow. Caller is responsible for sending
// every returned sequence, in order.
func kittyTransmit(imageID uint32, pngData []byte) []string {
	b64 := base64.StdEncoding.EncodeToString(pngData)

	var seqs []string
	for i := 0; i < len(b64); i += kittyChunkSize {
		end := min(i+kittyChunkSize, len(b64))
		chunk := b64[i:end]
		more := 1
		if end == len(b64) {
			more = 0
		}

		var ctrl string
		if i == 0 {
			ctrl = fmt.Sprintf("a=t,f=100,t=d,i=%d,m=%d,q=2", imageID, more)
		} else {
			ctrl = fmt.Sprintf("m=%d,q=2", more)
		}
		seqs = append(seqs, wrapForTmuxPassthrough(fmt.Sprintf("\033_G%s;%s\033\\", ctrl, chunk)))
	}
	return seqs
}

// kittyShow returns the sequence to display an already-transmitted
// image (via kittyTransmit) at the cursor's current position, scaled to
// occupy exactly cols x rows terminal cells — independent of the source
// image's actual resolution, which is what makes a small on-screen size
// with full source quality possible (Kitty does real image resampling
// for the display-time scaling, not a character-cell mosaic).
//
// C=1 tells Kitty not to move the cursor after placing the image —
// Kitty's own default post-placement cursor behavior isn't something
// that could be pinned down without a live terminal to test against, so
// this sidesteps the question entirely: the caller manages cursor
// position explicitly (see KittyLiveCheck's save/restore-cursor
// redraw-in-place technique, which avoids needing to know how many
// terminal rows the image actually rendered as).
func kittyShow(imageID uint32, cols, rows int) string {
	return wrapForTmuxPassthrough(fmt.Sprintf("\033_Ga=p,i=%d,c=%d,r=%d,C=1,q=2\033\\", imageID, cols, rows))
}

// kittyDelete returns the sequence to free an uploaded image's data —
// call once the display is done with an image id so it doesn't leak
// across many tutor turns over a session.
func kittyDelete(imageID uint32) string {
	return wrapForTmuxPassthrough(fmt.Sprintf("\033_Ga=d,i=%d,q=2\033\\", imageID))
}

// kittyLiveCheckImageID is fixed and outside the range newThinkingDisplay
// will ever use for a real turn (see thinkingdisplay.go once the full
// animation path exists) — recognizable in a Kitty debug dump, and safe
// to run without colliding with a real session's own uploaded frames.
// One consecutive id per real sprite frame.
//
// kittyLiveCheckCols is 15, not 16 — tuned down from an initial 16x8
// (a plain 2:1 cols:rows guess at typical monospace cell aspect) after
// live feedback that it came out very slightly wider than tall. Still a
// guess, just a closer one; cell pixel aspect isn't queryable from here
// and will vary by font/terminal anyway, so exact circularity can't be
// guaranteed for every setup.
const (
	kittyLiveCheckImageID   = 999901
	kittyLiveCheckCols      = 15
	kittyLiveCheckRows      = 8
	kittyLiveCheckTickDelay = 350 * time.Millisecond
)

// KittyLiveCheck transmits all 16 real disco ball frames via the Kitty
// graphics protocol, then loops through them in place — testing not
// just that a single image can be displayed, but the same
// redraw-in-place mechanics (cursor up, reprint, repeat) thinkingDisplay
// needs for the real animation, at closer to its real pacing than a
// short fixed number of ticks would show. Loops until interrupted
// (Ctrl-C); cleans up its uploaded images either way. See
// cmd/kitty-image-check.
//
// kittyLiveCheckFillerLines is printed before the animation starts —
// enough to push the pane past a typical terminal's bottom margin and
// force real scrolling before the loop begins. An earlier version of
// this check started animating near the top of a fresh, empty terminal
// window and never scrolled at all, which is exactly why it didn't
// catch a real bug: a since-reverted DECSC/DECRC (save/restore cursor)
// redraw technique that looked fine in that fresh-window test but
// duplicated content instead of overwriting once tried in a real,
// already-full tutor pane (see thinkingdisplay.go's redrawLocked doc
// comment). This check now always forces the same scrolled starting
// condition, so it actually exercises that scenario.
const kittyLiveCheckFillerLines = 60

// This is the one part of the disco ball rewrite that could not be
// verified in the environment it was built in — no Kitty install, no
// display. Everything here is a best-effort implementation of Kitty's
// documented graphics protocol; this function exists specifically to
// make that first real check as small and cheap to debug as possible.
func KittyLiveCheck(w io.Writer) error {
	imageIDs := make([]uint32, discoBallFrameCount)
	for i := 0; i < discoBallFrameCount; i++ {
		imageIDs[i] = kittyLiveCheckImageID + uint32(i)
		data, err := discoBallFramePNG(i)
		if err != nil {
			return fmt.Errorf("tutor: kitty live check: %w", err)
		}
		for _, seq := range kittyTransmit(imageIDs[i], data) {
			if _, err := io.WriteString(w, seq); err != nil {
				return fmt.Errorf("tutor: kitty live check: %w", err)
			}
		}
	}
	defer func() {
		for _, id := range imageIDs {
			io.WriteString(w, kittyDelete(id))
		}
	}()

	for i := 0; i < kittyLiveCheckFillerLines; i++ {
		if _, err := fmt.Fprintf(w, "filler conversation line %d\n", i); err != nil {
			return fmt.Errorf("tutor: kitty live check: %w", err)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	defer signal.Stop(stop)

	ticker := time.NewTicker(kittyLiveCheckTickDelay)
	defer ticker.Stop()

	// Relative cursor-up-by-rows math — the same scroll-safe technique
	// thinkingdisplay.go's redrawLocked uses (and the same reasoning:
	// content and cursor scroll together at the bottom margin, so this
	// stays correct regardless of where scrolling has put the block on
	// the physical screen, unlike DECSC/DECRC).
	drawn := 0
	for tick := 0; ; tick++ {
		select {
		case <-stop:
			return nil
		case <-ticker.C:
			if drawn > 0 {
				if _, err := fmt.Fprintf(w, "\033[%dA", drawn); err != nil {
					return fmt.Errorf("tutor: kitty live check: %w", err)
				}
			}
			frame := tick % discoBallFrameCount
			if _, err := io.WriteString(w, kittyShow(imageIDs[frame], kittyLiveCheckCols, kittyLiveCheckRows)); err != nil {
				return fmt.Errorf("tutor: kitty live check: %w", err)
			}
			for i := 0; i < kittyLiveCheckRows; i++ {
				if _, err := io.WriteString(w, "\n"); err != nil {
					return fmt.Errorf("tutor: kitty live check: %w", err)
				}
			}
			drawn = kittyLiveCheckRows
		}
	}
}

// kittyBallBaseImageID is fixed, not generated per turn — safe because
// kittyBallRenderer.close always deletes these ids before the *next*
// turn's init re-transmits under the same ones (tutor.go's Run loop is
// strictly sequential: one turn's display fully finishes, via finish(),
// before the next is constructed). Distinct from kittyLiveCheckImageID
// purely so the two are never confused when reading a debug capture.
const kittyBallBaseImageID = 900001

// kittyBallCols/Rows started at 15:8 (tuned live against a real Kitty
// terminal via cmd/kitty-image-check — close to, not exactly, the
// terminal's actual cell aspect ratio; see kittyShow's doc comment on
// why an exact ratio isn't knowable from here), halved to 8:4 for a
// more compact "still thinking" indicator, then halved again to 4:2 —
// confirmed live still too big/prominent next to the tool-call list at
// 8:4. Keeps the same 2:1 cols:rows ratio throughout. A starting point
// for further live tuning, not a measured-correct value.
const (
	kittyBallCols = 4
	kittyBallRows = 2
)

// kittyBallRenderer implements discoBallRenderer (thinkingdisplay.go)
// using the real Kitty graphics protocol — genuine full-quality sprite
// images, not the ANSI mosaic's downsampled approximation. Only ever
// constructed when kittyAvailable() (checked once by
// newThinkingDisplay, not here).
type kittyBallRenderer struct{}

func (kittyBallRenderer) init(w io.Writer) {
	for i := 0; i < discoBallFrameCount; i++ {
		data, err := discoBallFramePNG(i)
		if err != nil {
			continue // fail soft: a missing frame just won't show — this
			// renderer only ever runs for someone who could also have
			// gotten the guaranteed-correct ANSI fallback instead, so
			// there's nothing to escalate to.
		}
		for _, seq := range kittyTransmit(kittyBallBaseImageID+uint32(i), data) {
			io.WriteString(w, seq)
		}
	}
}

func (kittyBallRenderer) showFrame(w io.Writer, frame int) {
	io.WriteString(w, kittyShow(kittyBallBaseImageID+uint32(frame%discoBallFrameCount), kittyBallCols, kittyBallRows))
	// kittyShow uses C=1 (asks Kitty not to move the cursor after
	// placement) — advance it ourselves so showFrame's contract holds
	// (cursor ends up exactly rows() lines below where it started, same
	// as the ANSI renderer's own \n-per-line prints). This was missing
	// here even before the DECSC/DECRC experiment — thinkingDisplay's
	// redraw never actually relied on it advancing correctly until now.
	for i := 0; i < kittyBallRows; i++ {
		io.WriteString(w, "\n")
	}
}

func (kittyBallRenderer) rows() int { return kittyBallRows }

func (kittyBallRenderer) close(w io.Writer) {
	for i := 0; i < discoBallFrameCount; i++ {
		io.WriteString(w, kittyDelete(kittyBallBaseImageID+uint32(i)))
	}
}
