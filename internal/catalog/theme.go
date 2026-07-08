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

var bannerArt = []string{
	`______  ___   _      _     ______ _____  ________  ___ `,
	`| ___ \/ _ \ | |    | |    | ___ \  _  ||  _  |  \/  |`,
	`| |_/ / /_\ \| |    | |    | |_/ / | | || | | | .  . |`,
	`| ___ \  _  || |    | |    |    /| | | || | | | |\/| |`,
	`| |_/ / | | || |____| |____| |\ \\ \_/ /\ \_/ / |  | |`,
	`\____/\_| |_/\_____/\_____/\_| \_|\___/  \___/\_|  |_/`,
}

var bannerGradient = []string{colorTeal1, colorTeal2, colorWhite1, colorWhite2, colorPurple, colorPink}

// Banner renders the "BALLROOM" title art in a teal -> white -> purple/pink
// gradient, disco-ball style, with a tagline underneath.
func Banner() string {
	var b strings.Builder
	b.WriteString("\n")
	for i, line := range bannerArt {
		b.WriteString("  ")
		b.WriteString(styled(bannerGradient[i%len(bannerGradient)], line))
		b.WriteString("\n")
	}
	b.WriteString("  ")
	b.WriteString(styled(colorPink, "✦"))
	b.WriteString(" ")
	b.WriteString(styled(ansiBold+colorWhite1, "I N T E R V I E W   P R E P"))
	b.WriteString(" ")
	b.WriteString(styled(colorTeal1, "✦"))
	b.WriteString("\n")
	return b.String()
}
