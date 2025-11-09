package logger

import (
	"fmt"
	"strings"
)

const reset = "\x1b[0m"

type RGB struct{ R, G, B uint8 }

func Fg(s string, c RGB) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s%s", c.R, c.G, c.B, s, reset)
}
func Bg(s string, c RGB) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm%s%s", c.R, c.G, c.B, s, reset)
}
func FgBg(s string, fg, bg RGB) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%s%s",
		fg.R, fg.G, fg.B, bg.R, bg.G, bg.B, s, reset)
}

func FgLines(s string, c RGB) string {
	lines, trail := splitKeepTrail(s)
	for i, ln := range lines {
		lines[i] = Fg(ln, c)
	}
	return strings.Join(lines, "\n") + trail
}
func FgBgLines(s string, fg, bg RGB) string {
	lines, trail := splitKeepTrail(s)
	for i, ln := range lines {
		lines[i] = FgBg(ln, fg, bg)
	}
	return strings.Join(lines, "\n") + trail
}
func splitKeepTrail(s string) (lines []string, trailing string) {
	if strings.HasSuffix(s, "\n") {
		trailing = "\n"
		s = strings.TrimSuffix(s, "\n")
	}
	return strings.Split(s, "\n"), trailing
}

// ----- presets -----

var (
	Red   = RGB{255, 0, 0}
	Green = RGB{0, 255, 0}
	Blue  = RGB{0, 0, 255}

	SoftYellowBG = RGB{0xFF, 0xF8, 0xE1}
	SoftGreenBG  = RGB{0xEC, 0xFD, 0xF5}

	DimGray = RGB{120, 120, 120}
)

// Colorizer holds a name and the function that applies color.
// Use *Colorizer so nil means “no color”.
type Colorizer struct {
	Name string
	Fn   func(string) string
}

// Apply safely applies the colorizer if present.
func (c *Colorizer) Apply(s string) string {
	if c == nil || c.Fn == nil {
		return s
	}
	return c.Fn(s)
}

// Registry of reusable colorizers by name.
var Colorizers = map[string]Colorizer{
	"Red":         {Name: "Red", Fn: func(s string) string { return FgLines(s, Red) }},
	"Green":       {Name: "Green", Fn: func(s string) string { return FgLines(s, Green) }},
	"Blue":        {Name: "Blue", Fn: func(s string) string { return FgLines(s, Blue) }},
	"OnSoftYellow": {
		Name: "OnSoftYellow",
		Fn:   func(s string) string { return FgBgLines(s, RGB{0x43, 0x62, 0x12}, SoftYellowBG) },
	},
	"OnSoftGreen": {
		Name: "OnSoftGreen",
		Fn:   func(s string) string { return FgBgLines(s, RGB{0x16, 0x65, 0x34}, SoftGreenBG) },
	},
	"Dim": {Name: "Dim", Fn: func(s string) string { return FgLines(s, DimGray) }},
	"NoColor": {Name: "", Fn: nil},
}

// Convenience aliases you can import in call sites.
var (
	RedText      = Colorizers["Red"]
	GreenText    = Colorizers["Green"]
	BlueText     = Colorizers["Blue"]
	OnSoftYellow = Colorizers["OnSoftYellow"]
	OnSoftGreen  = Colorizers["OnSoftGreen"]
	DimText      = Colorizers["Dim"]
	NoColor      = Colorizers["NoColor"]
)

// for dynamic additions at runtime.
func RegisterColorizer(name string, fn func(string) string) Colorizer {
	c := Colorizer{Name: name, Fn: fn}
	Colorizers[name] = c
	return c
}
