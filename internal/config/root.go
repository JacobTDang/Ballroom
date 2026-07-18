package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// rootCacheEnvVar overrides where the last known-good checkout root is
// cached (see ResolveRoot) — used by tests to stay hermetic (isolated
// from the developer's real machine-wide cache), and available as an
// escape hatch if the OS default user-config dir isn't writable.
const rootCacheEnvVar = "PRACTICE_ROOT_CACHE_FILE"

// looksLikeCheckout reports whether dir contains docker/Dockerfile — the
// cheapest reliable signal that dir is an actual ballroom checkout (which
// docker build needs as its build context and -f path base, and which
// ExercisesDir/TestsDir/DataDir need as their parent), not just whatever
// directory the ballroom binary happened to be launched from.
func looksLikeCheckout(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "docker", "Dockerfile"))
	return err == nil
}

// rootCacheFile resolves the cache file's path: the override env var if
// set, otherwise a fixed path under the OS user-config dir.
func rootCacheFile() (string, error) {
	if p := os.Getenv(rootCacheEnvVar); p != "" {
		return p, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(dir, "ballroom", "root"), nil
}

// loadCachedRoot reads the last known-good checkout root cached at path.
// Returns ("", false) if unset, unreadable, or empty — best-effort, never
// an error the caller needs to handle specially.
func loadCachedRoot(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	root := strings.TrimSpace(string(data))
	return root, root != ""
}

// cacheRoot best-effort persists root at path for future runs launched
// from outside the checkout (e.g. `ballroom` on PATH after `go install`).
// Failures here (e.g. a read-only config dir) are silently ignored — the
// caller already resolved fine this run, so there's nothing to fail.
func cacheRoot(path, root string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(path, []byte(root), 0o644)
}

// ResolveRoot resolves the ballroom checkout root to use, given a
// candidate (usually PRACTICE_ROOT or the current working directory):
// candidate itself if it looks like a real ballroom checkout, else the
// last root that did (cached from a previous successful call). Without
// this, `ballroom` launched from PATH outside the checkout (its usual
// location once installed via `go install`) either fails with a raw,
// confusing error straight from docker, or — for Load's
// ExercisesDir/TestsDir/DataDir — silently resolves to empty/wrong
// directories (an empty picker, "couldn't load your progress") instead
// of either finding the real checkout or explaining clearly what's
// wrong.
//
// Shared by config.Load (exercises/tests/data dir resolution) and
// internal/orchestrator's docker-build root resolution (previously two
// separate copies of this exact logic — see that package's root.go,
// which now just delegates here).
func ResolveRoot(candidate string) (string, error) {
	if looksLikeCheckout(candidate) {
		if path, err := rootCacheFile(); err == nil {
			cacheRoot(path, candidate)
		}
		return candidate, nil
	}

	if path, err := rootCacheFile(); err == nil {
		if cached, ok := loadCachedRoot(path); ok && looksLikeCheckout(cached) {
			return cached, nil
		}
	}

	return "", fmt.Errorf("%s doesn't look like the ballroom repo (missing docker/Dockerfile) — run ballroom from the repo checkout at least once, or set PRACTICE_ROOT", candidate)
}
