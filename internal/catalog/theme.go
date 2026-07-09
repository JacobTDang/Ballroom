package catalog

import (
	"os"
	"strings"
)

// Bold/saturated retro palette, truecolor (24-bit) ANSI so the exact hex
// values render instead of the nearest 256-color approximation.
const (
	ansiReset = "\x1b[0m"
	ansiBold  = "\x1b[1m"

	colorRed    = "\x1b[38;2;240;60;60m"  // #F03C3C — dominant
	colorOrange = "\x1b[38;2;240;134;46m" // #F0862E
	colorGold   = "\x1b[38;2;232;169;60m" // #E8A93C
	colorPink   = "\x1b[38;2;224;70;140m" // #E0468C
	colorPurple = "\x1b[38;2;155;95;176m" // #9B5FB0
	colorBlue   = "\x1b[38;2;60;125;196m" // #3C7DC4
	colorTeal   = "\x1b[38;2;47;166;166m" // #2FA6A6

	colorCream    = "\x1b[38;2;242;235;221m" // #F2EBDD — warm off-white
	colorPaleGray = "\x1b[38;2;217;211;196m" // #D9D3C4

	// No green in this palette — teal reads as the "good" cool color,
	// and red (the palette's dominant hue) is an unambiguous fail.
	colorPass = colorTeal
	colorFail = colorRed
	colorDim  = colorPaleGray
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
// block) — filled in completely. An earlier attempt used a two-tone font
// with a dim texture in each glyph's negative space, but that destroyed
// the letter silhouettes entirely — from normal viewing distance it read
// as colored static, not "BALLROOM". Texture now comes from the mosaic
// coloring alone (many colors within a solid, unambiguous letter shape),
// not from filling the background too.
var bannerArt = []string{
	`████████     ███    ██       ██       ████████   ███████   ███████  ██     ██ `,
	`██     ██   ██ ██   ██       ██       ██     ██ ██     ██ ██     ██ ███   ███ `,
	`██     ██  ██   ██  ██       ██       ██     ██ ██     ██ ██     ██ ████ ████ `,
	`████████  ██     ██ ██       ██       ████████  ██     ██ ██     ██ ██ ███ ██ `,
	`██     ██ █████████ ██       ██       ██   ██   ██     ██ ██     ██ ██     ██ `,
	`██     ██ ██     ██ ██       ██       ██    ██  ██     ██ ██     ██ ██     ██ `,
	`████████  ██     ██ ████████ ████████ ██     ██  ███████   ███████  ██     ██ `,
}

var bannerGradient = []string{colorRed, colorOrange, colorGold, colorPink, colorPurple, colorBlue, colorTeal}

// mosaicWidth is how many columns share a color before the mosaic shifts
// to the next one — small facets, like a disco ball's mirror tiles,
// rather than one smooth gradient band per row.
const mosaicWidth = 3

// bannerScale repeats each source column this many times — a nearest-
// neighbor pixel-art scale-up (the same technique retro game sprites
// use), so the banner reads as genuinely bigger, not just the same art
// with more padding around it.
const bannerScale = 2

// MosaicBanner renders the BALLROOM title art as a scattered multi-color
// mosaic (each small block of characters gets its own color from the
// palette, diagonally offset per row) instead of a smooth per-row
// gradient — closer to light scattering off a disco ball than a sunset
// gradient. Only the letter strokes are colored; true blank space stays
// blank, so the silhouette stays unambiguous no matter how many colors
// are cycling through it. phase shifts the mosaic pattern; incrementing
// it on a timer (see internal/tui's tick handling) animates a shimmer
// across the letters. Pass phase=0 for a static render.
func MosaicBanner(phase int) string {
	return MosaicBannerScaled(phase, bannerScale)
}

// MosaicBannerScaled is MosaicBanner with an explicit pixel-scale instead
// of the default — a smaller scale (e.g. 1) fits the banner into a
// narrower column, such as sitting beside the disco ball in the two-
// column dashboard layout, while still keeping the same animated mosaic
// coloring.
func MosaicBannerScaled(phase, scale int) string {
	var b strings.Builder
	b.WriteString("\n")
	for row, line := range bannerArt {
		b.WriteString("  ")
		col := 0
		for _, ch := range line {
			if ch == ' ' {
				b.WriteString(strings.Repeat(" ", scale))
				col += scale
				continue
			}
			idx := (row + col/mosaicWidth + phase) % len(bannerGradient)
			b.WriteString(styled(bannerGradient[idx], strings.Repeat("█", scale)))
			col += scale
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
	return "  " + styled(colorPink, "✦") + " " + styled(ansiBold+colorCream, "BALLROOM — INTERVIEW PREP") + " " + styled(colorTeal, "✦") + "\n"
}

func tagline() string {
	return "  " + styled(colorPink, "✦") + " " + styled(ansiBold+colorCream, "I N T E R V I E W   P R E P") + " " + styled(colorTeal, "✦") + "\n"
}

// Banner is the static (non-animated) full title art, used on the boot
// screen where it's shown only briefly.
func Banner() string {
	return MosaicBanner(0)
}
