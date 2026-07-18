package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestElapsedMin_BasicArithmetic(t *testing.T) {
	if got := ElapsedMin(100, 700); got != 10 {
		t.Errorf("ElapsedMin(100, 700) = %v, want 10", got)
	}
}

func TestElapsedMin_SubMinuteRemainder(t *testing.T) {
	if got := ElapsedMin(0, 90); got != 1.5 {
		t.Errorf("ElapsedMin(0, 90) = %v, want 1.5", got)
	}
}

func TestElapsedMin_ZeroWhenNoTimeHasPassed(t *testing.T) {
	if got := ElapsedMin(500, 500); got != 0 {
		t.Errorf("ElapsedMin(500, 500) = %v, want 0", got)
	}
}

func TestElapsedMin_FractionalUptimeReadings(t *testing.T) {
	// /proc/uptime carries two decimals (e.g. "12345.67") -- the pure
	// function itself must handle that without any caller-side rounding.
	// Tolerance, not exact equality: float64 subtraction of two decimals
	// like these doesn't land on exactly 600.0 (classic binary floating
	// point representation error, not a bug in ElapsedMin).
	got := ElapsedMin(1000.12, 1600.12)
	if got < 9.999 || got > 10.001 {
		t.Errorf("ElapsedMin(1000.12, 1600.12) = %v, want ~10", got)
	}
}

// withFakeUptimeFile points procUptimePath at a fixture for the duration
// of the test, restoring the real path after -- lets containerUptime be
// tested deterministically on any host (including macOS, which has no
// /proc at all) instead of depending on this machine's real uptime.
func withFakeUptimeFile(t *testing.T, content string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "uptime")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write fixture uptime file: %v", err)
	}
	orig := procUptimePath
	procUptimePath = path
	t.Cleanup(func() { procUptimePath = orig })
}

func TestContainerUptime_ParsesFirstFieldOfProcUptime(t *testing.T) {
	withFakeUptimeFile(t, "12345.67 98765.43\n")

	got, ok := containerUptime()
	if !ok {
		t.Fatal("ok = false, want true for a well-formed /proc/uptime")
	}
	if got != 12345.67 {
		t.Errorf("containerUptime() = %v, want 12345.67", got)
	}
}

func TestContainerUptime_MissingFileReturnsNotOK(t *testing.T) {
	orig := procUptimePath
	procUptimePath = filepath.Join(t.TempDir(), "does-not-exist")
	t.Cleanup(func() { procUptimePath = orig })

	if _, ok := containerUptime(); ok {
		t.Error("ok = true for a missing file, want false (non-Linux hosts, e.g. macOS test runs)")
	}
}

func TestContainerUptime_UnparsableContentReturnsNotOK(t *testing.T) {
	withFakeUptimeFile(t, "not-a-number\n")

	if _, ok := containerUptime(); ok {
		t.Error("ok = true for unparsable content, want false")
	}
}

func TestContainerUptime_EmptyFileReturnsNotOK(t *testing.T) {
	withFakeUptimeFile(t, "")

	if _, ok := containerUptime(); ok {
		t.Error("ok = true for an empty file, want false")
	}
}
