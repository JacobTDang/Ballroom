package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// rootCacheEnvVar overrides where the last known-good checkout root is
// cached (see dockerBuildRoot) — used by tests to stay hermetic (isolated
// from the developer's real machine-wide cache), and available as an
// escape hatch if the OS default user-config dir isn't writable.
const rootCacheEnvVar = "PRACTICE_ROOT_CACHE_FILE"

// looksLikeCheckout reports whether dir contains docker/Dockerfile — the
// cheapest reliable signal that dir is an actual ballroom checkout (which
// docker build needs as its build context and -f path base), not just
// whatever directory the ballroom binary happened to be launched from.
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
// Failures here (e.g. a read-only config dir) are silently ignored — Root
// itself already resolved fine this run, so there's nothing to fail.
func cacheRoot(path, root string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(path, []byte(root), 0o644)
}

// dockerBuildRoot resolves the directory to use as docker build's -f path
// base and build context: cfgRoot itself if it looks like a real ballroom
// checkout, else the last root that did (cached from a previous
// successful run). Without this, `ballroom` launched from PATH outside
// the checkout (its usual location once installed via `go install`)
// fails with a raw, confusing error straight from docker — "lstat
// docker: no such file or directory" — instead of either finding the
// real checkout or explaining clearly what's wrong.
func dockerBuildRoot(cfgRoot string) (string, error) {
	if looksLikeCheckout(cfgRoot) {
		if path, err := rootCacheFile(); err == nil {
			cacheRoot(path, cfgRoot)
		}
		return cfgRoot, nil
	}

	if path, err := rootCacheFile(); err == nil {
		if cached, ok := loadCachedRoot(path); ok && looksLikeCheckout(cached) {
			return cached, nil
		}
	}

	return "", fmt.Errorf("%s doesn't look like the ballroom repo (missing docker/Dockerfile) — run ballroom from the repo checkout at least once, or set PRACTICE_ROOT", cfgRoot)
}
