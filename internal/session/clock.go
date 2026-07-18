package session

import (
	"os"
	"strconv"
	"strings"
)

// procUptimePath is a var, not a const, so tests can point containerUptime
// at a fixture file instead of the real /proc/uptime -- which doesn't
// exist at all on a macOS/Windows test host, and even on Linux reflects
// whatever this machine's real uptime happens to be, not something a
// test can control. See clock_test.go.
var procUptimePath = "/proc/uptime"

// containerUptime reads the first field of /proc/uptime: seconds since
// this container's kernel booted (see docker/entrypoint.sh and
// docker/clock.sh, which read the same file for the same reason). Two
// readings taken minutes apart, fed to ElapsedMin, measure real elapsed
// time even if the host laptop slept in between -- the Docker Desktop
// Linux VM a session runs in suspends along with the host, so uptime
// does not advance across a lid-close the way wall clock does (issue
// #229). ok is false when the file doesn't exist (non-Linux hosts, e.g.
// running `go test` on macOS) or its content doesn't parse, telling the
// caller to fall back to wall clock instead of trusting a zero value.
func containerUptime() (float64, bool) {
	data, err := os.ReadFile(procUptimePath)
	if err != nil {
		return 0, false
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, false
	}
	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, false
	}
	return uptime, true
}

// ElapsedMin is the pure arithmetic behind the uptime-based clock:
// minutes between two /proc/uptime readings. Kept separate from
// containerUptime so the arithmetic itself is unit-testable without a
// real (or fixture) /proc/uptime file. Callers are expected to pass
// startUptime <= nowUptime -- uptime only resets if the container itself
// restarts, which doesn't happen mid-session.
func ElapsedMin(startUptime, nowUptime float64) float64 {
	return (nowUptime - startUptime) / 60
}
