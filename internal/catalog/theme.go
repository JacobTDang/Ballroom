package catalog

import (
	"os"
	"strings"
)

// Retro/disco palette: teal/cyan through white to purple/pink, plus
// semantic colors for pass/fail/not-attempted. 256-color ANSI codes —
// broadly supported, no truecolor assumption needed.
const (
	ansiReset = "\x1b[0m"
	ansiBold  = "\x1b[1m"

	colorTeal1  = "\x1b[38;5;51m"  // bright cyan
	colorTeal2  = "\x1b[38;5;45m"  // turquoise
	colorWhite1 = "\x1b[38;5;87m"  // pale cyan-white
	colorWhite2 = "\x1b[38;5;189m" // pale lavender-white
	colorPurple = "\x1b[38;5;141m" // medium purple
	colorPink   = "\x1b[38;5;213m" // pink/magenta
	colorPink2  = "\x1b[38;5;207m" // deeper magenta, gradient's last stop

	colorPass = "\x1b[38;5;120m" // soft green
	colorFail = "\x1b[38;5;210m" // soft red/salmon — still reads as "fail" but stays in-palette
	colorDim  = "\x1b[38;5;244m" // gray, for secondary text (ids, "not attempted")
)

// colorEnabled is re-checked on every call (not cached at package init) so
// tests can toggle NO_COLOR with t.Setenv and see the effect immediately.
// Per https://no-color.org: only a non-empty NO_COLOR disables color.
func colorEnabled() bool {
	return os.Getenv("NO_COLOR") == ""
}

// styled wraps s in the given ANSI codes. Callers that also pad s to a
// fixed width MUST pad first, then style — styling first would make the
// escape-code bytes count toward the padding width and break column
// alignment.
func styled(codes, s string) string {
	if !colorEnabled() || s == "" {
		return s
	}
	return codes + s + ansiReset
}

// Solid block-letter art (figlet "banner3", # swapped for a solid Unicode
// block) — filled in completely, not an outline font.
var bannerArt = []string{
	`████████     ███    ██       ██       ████████   ███████   ███████  ██     ██ `,
	`██     ██   ██ ██   ██       ██       ██     ██ ██     ██ ██     ██ ███   ███ `,
	`██     ██  ██   ██  ██       ██       ██     ██ ██     ██ ██     ██ ████ ████ `,
	`████████  ██     ██ ██       ██       ████████  ██     ██ ██     ██ ██ ███ ██ `,
	`██     ██ █████████ ██       ██       ██   ██   ██     ██ ██     ██ ██     ██ `,
	`██     ██ ██     ██ ██       ██       ██    ██  ██     ██ ██     ██ ██     ██ `,
	`████████  ██     ██ ████████ ████████ ██     ██  ███████   ███████  ██     ██ `,
}

var bannerGradient = []string{colorTeal1, colorTeal2, colorWhite1, colorWhite2, colorPurple, colorPink, colorPink2}

// mosaicWidth is how many columns share a color before the mosaic shifts
// to the next one — small facets, like a disco ball's mirror tiles,
// rather than one smooth gradient band per row.
const mosaicWidth = 3

// MosaicBanner renders the BALLROOM title art as a scattered multi-color
// mosaic (each small block of characters gets its own color from the
// palette, diagonally offset per row) instead of a smooth per-row
// gradient — closer to light scattering off a disco ball than a sunset
// gradient. phase shifts the pattern; incrementing it on a timer (see
// internal/tui's tick handling) animates a shimmer across the letters.
// Pass phase=0 for a static render.
func MosaicBanner(phase int) string {
	var b strings.Builder
	b.WriteString("\n")
	for row, line := range bannerArt {
		b.WriteString("  ")
		col := 0
		for _, ch := range line {
			if ch == ' ' {
				b.WriteRune(ch)
				col++
				continue
			}
			idx := (row + col/mosaicWidth + phase) % len(bannerGradient)
			b.WriteString(styled(bannerGradient[idx], string(ch)))
			col++
		}
		b.WriteString("\n")
	}
	b.WriteString(tagline())
	return b.String()
}

// CompactBanner is a single-line wordmark for screens that need the
// vertical space back (the tree picker, stats) — full art stays reserved
// for the boot screen and the main menu's title moment.
func CompactBanner() string {
	return "  " + styled(colorPink, "✦") + " " + styled(ansiBold+colorWhite1, "BALLROOM — INTERVIEW PREP") + " " + styled(colorTeal1, "✦") + "\n"
}

func tagline() string {
	return "  " + styled(colorPink, "✦") + " " + styled(ansiBold+colorWhite1, "I N T E R V I E W   P R E P") + " " + styled(colorTeal1, "✦") + "\n"
}

// Banner is the static (non-animated) full title art, used on the boot
// screen where it's shown only briefly.
func Banner() string {
	return MosaicBanner(0)
}
