package orchestrator

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/JacobTDang/Ballroom/internal/config"
)

// ballroomImageLabel is baked into docker/Dockerfile so cleanupDangling
// can find and remove this project's own stale images after a rebuild,
// without ever touching unrelated dangling images from other projects on
// the same machine.
const ballroomImageLabel = "com.ballroom.app"

// EnsureImage builds cfg.DockerImage if it doesn't already exist,
// printing docker build's output directly to stdout as it runs — used by
// the CLI's direct exercise/sandbox launch paths (`ballroom practice`,
// `ballroom sandbox`), which don't go through the boot screen's own
// live-streaming build UI. Errors are non-fatal: if Docker itself isn't
// reachable, there's nothing to do here — the boot screen's own checks
// surface that clearly instead.
func EnsureImage(cfg config.Config) error {
	if exec.Command("docker", "info").Run() != nil {
		return nil
	}

	out, err := exec.Command("docker", "image", "inspect", cfg.DockerImage, "--format", "{{.Id}}").Output()
	if err == nil && strings.TrimSpace(string(out)) != "" {
		cleanupDanglingBallroomImages()
		return nil
	}

	fmt.Printf("Practice image %q not found — building it now (this can take a minute or two)...\n\n", cfg.DockerImage)
	lineCh, errCh := BuildImage(cfg)
	for line := range lineCh {
		fmt.Println(line)
	}
	buildErr := <-errCh
	fmt.Println()
	return buildErr
}

// BuildImage runs `docker build` for cfg.DockerImage, streaming each
// output line on the returned channel (closed once the build's output
// ends) and sending exactly one final result (nil on success) on the
// error channel. On success, also cleans up stale ballroom images left
// over from a previous build (see cleanupDanglingBallroomImages) — the
// caller doesn't need to do this separately.
func BuildImage(cfg config.Config) (<-chan string, <-chan error) {
	lineCh := make(chan string, 200)
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		pr, pw := io.Pipe()
		cmd := exec.Command("docker", "build", "-f", "docker/Dockerfile", "-t", cfg.DockerImage, ".")
		cmd.Dir = cfg.Root
		cmd.Stdout = pw
		cmd.Stderr = pw

		scanDone := make(chan struct{})
		go func() {
			defer close(scanDone)
			defer close(lineCh)
			scanner := bufio.NewScanner(pr)
			scanner.Buffer(make([]byte, 64*1024), 1024*1024)
			for scanner.Scan() {
				lineCh <- scanner.Text()
			}
		}()

		if err := cmd.Start(); err != nil {
			pw.Close()
			<-scanDone
			errCh <- fmt.Errorf("start docker build: %w", err)
			return
		}

		waitErr := cmd.Wait()
		pw.Close()
		<-scanDone

		if waitErr != nil {
			errCh <- fmt.Errorf("docker build: %w", waitErr)
			return
		}

		cleanupDanglingBallroomImages()
		errCh <- nil
	}()

	return lineCh, errCh
}

// cleanupDanglingBallroomImages removes untagged images left behind by a
// previous build of this project's image (see ballroomImageLabel).
// Best-effort: failures here shouldn't block starting the app.
func cleanupDanglingBallroomImages() {
	out, err := exec.Command("docker", "images",
		"--filter", "dangling=true",
		"--filter", "label="+ballroomImageLabel,
		"-q",
	).Output()
	if err != nil {
		return
	}

	ids := parseImageIDs(string(out))
	if len(ids) == 0 {
		return
	}

	args := append([]string{"rmi"}, ids...)
	if exec.Command("docker", args...).Run() != nil {
		return
	}
	fmt.Printf("Cleaned up %d old ballroom image build(s)\n", len(ids))
}

// parseImageIDs splits `docker ... -q` output (one id per line) into a
// slice, dropping blank lines.
func parseImageIDs(out string) []string {
	return strings.Fields(out)
}
