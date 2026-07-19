package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Terminal geometry. The cell size is what makes the output line up:
// every glyph is placed at an exact grid position and stretched to
// exactly one cell wide, so the SVG matches the terminal's own layout
// rather than whatever the viewer's font would do on its own.
const (
	cellW     = 8.4
	lineH     = 17.6
	fontSize  = 14.0
	padX      = 12.0
	padY      = 12.0
	defaultFg = "#f2ebdd"
	defaultBg = "#101014"
)

// style is the SGR state carried across a line.
type style struct {
	fg, bg                       string
	bold, dim, italic, underline bool
	reverse                      bool
}

func newStyle() style { return style{fg: defaultFg, bg: defaultBg} }

type cell struct {
	r  rune
	st style
}

// base16 is the classic 8+8 ANSI set, used for the basic color codes
// (30-37, 90-97) that a few tools still emit.
var base16 = []string{
	"#1c1c22", "#f03c3c", "#4fbf67", "#e8a93c", "#4f8fbf", "#c678dd", "#2fa6a6", "#c7c7c7",
	"#686868", "#ff6b6b", "#6fdf87", "#ffd25c", "#6fafdf", "#e69bff", "#4fc6c6", "#ffffff",
}

// xterm256 resolves a 256-color index: the 16 base colors, then the
// 6x6x6 cube, then the 24-step grayscale ramp.
func xterm256(n int) string {
	switch {
	case n < 16:
		return base16[n]
	case n < 232:
		n -= 16
		levels := []int{0, 95, 135, 175, 215, 255}
		return fmt.Sprintf("#%02x%02x%02x", levels[n/36], levels[(n/6)%6], levels[n%6])
	default:
		v := 8 + (n-232)*10
		return fmt.Sprintf("#%02x%02x%02x", v, v, v)
	}
}

// applySGR folds one escape sequence's parameters into the running
// style. Unknown parameters are ignored rather than treated as an
// error: a capture is a recording of whatever the terminal emitted,
// and refusing to render it because of one unrecognized code would be
// worse than rendering it plainly.
func applySGR(s style, params string) style {
	if params == "" {
		return newStyle()
	}
	ps := strings.Split(params, ";")
	for i := 0; i < len(ps); i++ {
		n, err := strconv.Atoi(ps[i])
		if err != nil {
			continue
		}
		switch {
		case n == 0:
			s = newStyle()
		case n == 1:
			s.bold = true
		case n == 2:
			s.dim = true
		case n == 3:
			s.italic = true
		case n == 4:
			s.underline = true
		case n == 7:
			s.reverse = true
		case n == 22:
			s.bold, s.dim = false, false
		case n == 23:
			s.italic = false
		case n == 24:
			s.underline = false
		case n == 27:
			s.reverse = false
		case n >= 30 && n <= 37:
			s.fg = base16[n-30]
		case n >= 90 && n <= 97:
			s.fg = base16[n-90+8]
		case n >= 40 && n <= 47:
			s.bg = base16[n-40]
		case n >= 100 && n <= 107:
			s.bg = base16[n-100+8]
		case n == 39:
			s.fg = defaultFg
		case n == 49:
			s.bg = defaultBg
		case n == 38 || n == 48:
			col, consumed := extendedColor(ps[i+1:])
			if col != "" {
				if n == 38 {
					s.fg = col
				} else {
					s.bg = col
				}
			}
			i += consumed
		}
	}
	return s
}

// extendedColor reads the 5;N (256-color) or 2;R;G;B (truecolor) form
// that follows a 38 or 48, returning the color and how many parameters
// it consumed.
func extendedColor(rest []string) (string, int) {
	if len(rest) == 0 {
		return "", 0
	}
	switch rest[0] {
	case "5":
		if len(rest) < 2 {
			return "", 0
		}
		n, err := strconv.Atoi(rest[1])
		if err != nil {
			return "", 0
		}
		return xterm256(n), 2
	case "2":
		if len(rest) < 4 {
			return "", 0
		}
		var rgb [3]int
		for i := 0; i < 3; i++ {
			v, err := strconv.Atoi(rest[i+1])
			if err != nil {
				return "", 0
			}
			rgb[i] = v
		}
		return fmt.Sprintf("#%02x%02x%02x", rgb[0], rgb[1], rgb[2]), 4
	}
	return "", 0
}

// parse turns an ANSI capture into a grid of styled cells, padding
// every row to width so background runs stay rectangular.
func parse(input string, width int) [][]cell {
	var rows [][]cell
	cur := newStyle()
	for _, line := range strings.Split(strings.TrimRight(input, "\n"), "\n") {
		var row []cell
		rs := []rune(line)
		for i := 0; i < len(rs); {
			if rs[i] == 0x1b && i+1 < len(rs) && rs[i+1] == '[' {
				j := i + 2
				for j < len(rs) && !(rs[j] >= 0x40 && rs[j] <= 0x7e) {
					j++
				}
				if j < len(rs) && rs[j] == 'm' {
					cur = applySGR(cur, string(rs[i+2:j]))
				}
				i = j + 1
				continue
			}
			if rs[i] == '\r' {
				i++
				continue
			}
			row = append(row, cell{r: rs[i], st: cur})
			i++
		}
		for len(row) < width {
			row = append(row, cell{r: ' ', st: newStyle()})
		}
		if width > 0 && len(row) > width {
			row = row[:width]
		}
		rows = append(rows, row)
	}
	return rows
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// render emits the SVG. Two passes: background rectangles for runs
// whose background differs from the page, then one text element per
// run of same-styled non-space glyphs.
//
// Spaces are never rendered, only stepped over. A run that included
// them would need its textLength to cover them too, and any viewer
// whose font metrics differ even slightly from the assumed cell width
// then distributes that error across the whole run -- which shows up
// as visibly uneven letter spacing.
func render(rows [][]cell, cols int) string {
	w := padX*2 + float64(cols)*cellW
	h := padY*2 + float64(len(rows))*lineH

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n", w, h, w, h)
	fmt.Fprintf(&b, `<rect width="100%%" height="100%%" rx="8" fill="%s"/>`+"\n", defaultBg)

	for y, row := range rows {
		for x := 0; x < len(row); {
			bg := effectiveBg(row[x].st)
			run := x
			for run < len(row) && effectiveBg(row[run].st) == bg {
				run++
			}
			if bg != defaultBg {
				fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`+"\n",
					padX+float64(x)*cellW, padY+float64(y)*lineH, float64(run-x)*cellW, lineH, bg)
			}
			x = run
		}
	}

	fmt.Fprintf(&b, `<g font-family="SF Mono, Menlo, Consolas, monospace" font-size="%.0fpx" xml:space="preserve">`+"\n", fontSize)
	for y, row := range rows {
		for x := 0; x < len(row); {
			if row[x].r == ' ' {
				x++
				continue
			}
			st := row[x].st
			run := x
			var text []rune
			for run < len(row) && row[run].st == st && row[run].r != ' ' {
				text = append(text, row[run].r)
				run++
			}
			attrs := fmt.Sprintf(`x="%.1f" y="%.1f" textLength="%.1f" lengthAdjust="spacingAndGlyphs" fill="%s"`,
				padX+float64(x)*cellW, padY+float64(y)*lineH+fontSize-1, float64(run-x)*cellW, effectiveFg(st))
			if st.bold {
				attrs += ` font-weight="600"`
			}
			if st.dim {
				attrs += ` opacity="0.62"`
			}
			if st.italic {
				attrs += ` font-style="italic"`
			}
			if st.underline {
				attrs += ` text-decoration="underline"`
			}
			fmt.Fprintf(&b, "<text %s>%s</text>\n", attrs, escapeXML(string(text)))
			x = run
		}
	}
	b.WriteString("</g></svg>\n")
	return b.String()
}

func effectiveFg(s style) string {
	if s.reverse {
		return s.bg
	}
	return s.fg
}

func effectiveBg(s style) string {
	if s.reverse {
		return s.fg
	}
	return s.bg
}
