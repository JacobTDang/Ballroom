package orchestrator

import (
	"fmt"
	"os"
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
// streaming docker build's normal output directly to the terminal — this
// runs before the TUI starts, so there's no screen to conflict with. It
// then removes any dangling images left over from a previous build of
// this image. Errors are non-fatal: if Docker itself isn't reachable,
// there's nothing to do here — the boot screen's own checks will
// surface that clearly instead.
func EnsureImage(cfg config.Config) error {
	if exec.Command("docker", "info").Run() != nil {
		return nil
	}

	out, err := exec.Command("docker", "image", "inspect", cfg.DockerImage, "--format", "{{.Id}}").Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		fmt.Printf("Practice image %q not found — building it now (this can take a minute or two)...\n\n", cfg.DockerImage)

		build := exec.Command("docker", "build", "-f", "docker/Dockerfile", "-t", cfg.DockerImage, ".")
		build.Dir = cfg.Root
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr
		if err := build.Run(); err != nil {
			return fmt.Errorf("build image: %w", err)
		}
		fmt.Println()
	}

	cleanupDanglingBallroomImages()
	return nil
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
