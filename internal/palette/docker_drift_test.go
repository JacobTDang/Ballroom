package palette

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"
)

// The container's config (docker/tmux.conf, docker/clock.sh, the nvim
// colorscheme) can't import this package -- it isn't Go -- so it's the
// one color home this package's own doc comment says a drift test has
// to police from the outside, by grepping. These two tests are that
// drift test: every hex color under docker/ must be a real palette
// color, and nothing under docker/ may fall back to the terminal's
// approximated 256-color ramp instead of one.

// hexColorPattern matches a #RRGGBB literal the way every config file
// under docker/ spells one -- tmux.conf's #[fg=#RRGGBB] directives,
// clock.sh's printf'd colors, and the nvim colorscheme's hl() calls.
var hexColorPattern = regexp.MustCompile(`#[0-9a-fA-F]{6}\b`)

// indexedColorPattern matches tmux's colourNNN / colorNNN 256-color
// index syntax (0-255) -- the thing this project's palette exists to
// replace. Before this change, docker/tmux.conf's status-style carried
// exactly one of these (bg=colour234); this pattern is what makes that
// a test failure instead of a silent drift, and what stops a future
// edit from reintroducing one.
var indexedColorPattern = regexp.MustCompile(`\bcolou?r[0-9]{1,3}\b`)

// dockerDir resolves the repo's docker/ directory relative to this test
// file (the same runtime.Caller trick internal/tutor/nvimrpc_test.go
// uses for docker/nvim) -- go test's working directory is always the
// package directory, not the repo root, so a bare relative path would
// break depending on how the test binary gets invoked.
func dockerDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	abs, err := filepath.Abs(filepath.Join(filepath.Dir(thisFile), "..", "..", "docker"))
	if err != nil {
		t.Fatalf("resolve docker dir: %v", err)
	}
	if info, err := os.Stat(abs); err != nil || !info.IsDir() {
		t.Fatalf("docker dir not found at %s: %v", abs, err)
	}
	return abs
}

// walkDockerFiles calls fn with the contents of every regular file under
// docker/, relative path first (for readable failure messages).
func walkDockerFiles(t *testing.T, fn func(rel, content string)) {
	t.Helper()
	root := dockerDir(t)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		fn(filepath.ToSlash(rel), string(data))
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
}

// TestDockerConfigColorsAreInPalette is the drift test that keeps the
// container's chrome (tmux.conf, clock.sh, the nvim colorscheme) from
// ever disagreeing with this package: every #RRGGBB literal anywhere
// under docker/ must be a real palette color. Deliberately walks the
// whole docker/ tree rather than a hardcoded file list, so a new config
// file with its own invented color is caught exactly the same way an
// edit to an existing one would be.
func TestDockerConfigColorsAreInPalette(t *testing.T) {
	var checked int
	walkDockerFiles(t, func(rel, content string) {
		for _, hex := range hexColorPattern.FindAllString(content, -1) {
			checked++
			if !Contains(hex) {
				t.Errorf("%s: %s is not a palette color -- add it to palette.go or fix the typo", rel, hex)
			}
		}
	})
	if checked == 0 {
		t.Fatal("found no #RRGGBB colors under docker/ -- the walk or the pattern is broken")
	}
}

// TestDockerConfigHasNoIndexedColors bans tmux's colourNNN 256-color
// index syntax anywhere under docker/ -- the exact problem issue #260
// set out to fix (docker/tmux.conf's status-style used to carry
// bg=colour234, the only 256-color value in the codebase, approximating
// a color this package can already express exactly). A hex literal is
// still checked against the palette by TestDockerConfigColorsAreInPalette
// above; this test is what stops the indexed escape hatch from coming
// back instead of a real palette hex.
func TestDockerConfigHasNoIndexedColors(t *testing.T) {
	walkDockerFiles(t, func(rel, content string) {
		if m := indexedColorPattern.FindString(content); m != "" {
			t.Errorf("%s: found %q -- use the equivalent palette hex (see internal/palette), not a 256-color index", rel, m)
		}
	})
}
