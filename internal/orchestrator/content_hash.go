package orchestrator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// contentRoots are the paths (relative to the checkout root) that make up
// the practice image's actual content: everything docker/Dockerfile COPYs
// in (docker/, cmd/, internal/) plus the two files that pin the Go module
// graph the in-container ballroom binary gets built against (go.mod,
// go.sum). Anything outside this list (docs/, exercises/, tests/, data/,
// README.md, ...) never reaches the image, so changing it must never
// trigger a rebuild.
var contentRoots = []string{"docker", "cmd", "internal", "go.mod", "go.sum"}

// contentHashCache memoizes contentHash's result for the last root it was
// asked about. cfg.Root never changes within one ballroom process
// invocation, and a single rebuild needs the hash twice (EnsureImage to
// decide a rebuild is needed, then BuildImage again to stamp the label) —
// this makes the second call free instead of re-walking cmd/ and
// internal/ a second time.
var (
	contentHashMu    sync.Mutex
	contentHashCache struct {
		root string
		hash string
		set  bool
	}
)

// contentHash returns a deterministic, hex-encoded sha256 over every file
// under root's contentRoots that ends up baked into the practice image —
// covering both each file's relative path (so a rename changes the hash,
// not just an edit) and its content. Test files (*_test.go) and any
// testdata/ directory are excluded: they never reach docker/Dockerfile's
// COPY instructions, so editing them must never trigger a multi-minute
// image rebuild. A missing content root (e.g. a checkout that predates
// one of them) is skipped rather than treated as an error.
func contentHash(root string) (string, error) {
	contentHashMu.Lock()
	if contentHashCache.set && contentHashCache.root == root {
		h := contentHashCache.hash
		contentHashMu.Unlock()
		return h, nil
	}
	contentHashMu.Unlock()

	hash, err := computeContentHash(root)
	if err != nil {
		return "", err
	}

	contentHashMu.Lock()
	contentHashCache.root = root
	contentHashCache.hash = hash
	contentHashCache.set = true
	contentHashMu.Unlock()

	return hash, nil
}

// computeContentHash does the actual walk-and-hash work for contentHash,
// uncached.
func computeContentHash(root string) (string, error) {
	var relPaths []string

	for _, rel := range contentRoots {
		base := filepath.Join(root, rel)
		info, err := os.Stat(base)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", fmt.Errorf("orchestrator: content hash: stat %s: %w", base, err)
		}

		if !info.IsDir() {
			relPaths = append(relPaths, filepath.ToSlash(rel))
			continue
		}

		err = filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == "testdata" {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(d.Name(), "_test.go") {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relPaths = append(relPaths, filepath.ToSlash(rel))
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("orchestrator: content hash: walk %s: %w", base, err)
		}
	}

	sort.Strings(relPaths)

	h := sha256.New()
	for _, rel := range relPaths {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := hashFile(h, rel, full); err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// hashFile feeds one file's relative path and content into h, each
// separated by a NUL byte so e.g. path "ab" + content "c" can never hash
// the same as path "a" + content "bc".
func hashFile(h io.Writer, rel, full string) error {
	f, err := os.Open(full)
	if err != nil {
		return fmt.Errorf("orchestrator: content hash: open %s: %w", full, err)
	}
	defer f.Close()

	if _, err := io.WriteString(h, rel); err != nil {
		return fmt.Errorf("orchestrator: content hash: hash %s: %w", full, err)
	}
	if _, err := h.Write([]byte{0}); err != nil {
		return fmt.Errorf("orchestrator: content hash: hash %s: %w", full, err)
	}
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("orchestrator: content hash: read %s: %w", full, err)
	}
	if _, err := h.Write([]byte{0}); err != nil {
		return fmt.Errorf("orchestrator: content hash: hash %s: %w", full, err)
	}
	return nil
}
