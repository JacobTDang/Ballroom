package tui

import "strings"

// Retro-modern chrome: the app keeps its palette and its disco ball,
// and borrows structure from the terminals it's descended from --
// double-ruled frames, letterspaced headings, block-shaded meters. The
// goal is that it reads as deliberately old-fashioned rather than
// merely plain, without costing any of the readability the screens
// already had.

// heading renders a screen title the way the banner's tagline already
// did ("I N T E R V I E W   P R E P") -- uppercased, one space between
// letters, three between words so a two-word title doesn't read as one
// long one. Operates on runes: a byte loop would split multi-byte
// characters into garbage.
func heading(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	for i, word := range strings.Fields(strings.ToUpper(s)) {
		if i > 0 {
			b.WriteString("   ")
		}
		for j, r := range word {
			if j > 0 {
				b.WriteRune(' ')
			}
			b.WriteRune(r)
		}
	}
	return b.String()
}
