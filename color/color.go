package color

import "github.com/NLipatov/tuigo/internal/ansi"

type Color string

func (c Color) String() string {
	return string(c)
}

const (
	FgBlack  Color = Color(ansi.FG_BLACK)
	FgRed    Color = Color(ansi.FG_RED)
	FgGreen  Color = Color(ansi.FG_GREEN)
	FgYellow Color = Color(ansi.FG_YELLOW)
	FgBlue   Color = Color(ansi.FG_BLUE)
	FgPurple Color = Color(ansi.FG_PURPLE)
	FgCyan   Color = Color(ansi.FG_CYAN)
	FgWhite  Color = Color(ansi.FG_WHITE)

	FgBoldBlack  Color = Color(ansi.FG_BOLD_BLACK)
	FgBoldRed    Color = Color(ansi.FG_BOLD_RED)
	FgBoldGreen  Color = Color(ansi.FG_BOLD_GREEN)
	FgBoldYellow Color = Color(ansi.FG_BOLD_YELLOW)
	FgBoldBlue   Color = Color(ansi.FG_BOLD_BLUE)
	FgBoldPurple Color = Color(ansi.FG_BOLD_PURPLE)
	FgBoldCyan   Color = Color(ansi.FG_BOLD_CYAN)
	FgBoldWhite  Color = Color(ansi.FG_BOLD_WHITE)

	FgUnderlineBlack  Color = Color(ansi.FG_UNDERLINE_BLACK)
	FgUnderlineRed    Color = Color(ansi.FG_UNDERLINE_RED)
	FgUnderlineGreen  Color = Color(ansi.FG_UNDERLINE_GREEN)
	FgUnderlineYellow Color = Color(ansi.FG_UNDERLINE_YELLOW)
	FgUnderlineBlue   Color = Color(ansi.FG_UNDERLINE_BLUE)
	FgUnderlinePurple Color = Color(ansi.FG_UNDERLINE_PURPLE)
	FgUnderlineCyan   Color = Color(ansi.FG_UNDERLINE_CYAN)
	FgUnderlineWhite  Color = Color(ansi.FG_UNDERLINE_WHITE)

	FgHighIntensityBlack  Color = Color(ansi.FG_HIGH_INTENSITY_BLACK)
	FgHighIntensityRed    Color = Color(ansi.FG_HIGH_INTENSITY_RED)
	FgHighIntensityGreen  Color = Color(ansi.FG_HIGH_INTENSITY_GREEN)
	FgHighIntensityYellow Color = Color(ansi.FG_HIGH_INTENSITY_YELLOW)
	FgHighIntensityBlue   Color = Color(ansi.FG_HIGH_INTENSITY_BLUE)
	FgHighIntensityPurple Color = Color(ansi.FG_HIGH_INTENSITY_PURPLE)
	FgHighIntensityCyan   Color = Color(ansi.FG_HIGH_INTENSITY_CYAN)
	FgHighIntensityWhite  Color = Color(ansi.FG_HIGH_INTENSITY_WHITE)

	FgBoldHighIntensityBlack  Color = Color(ansi.FG_BOLD_HIGH_INTENSITY_BLACK)
	FgBoldHighIntensityRed    Color = Color(ansi.FG_BOLD_HIGH_INTENSITY_RED)
	FgBoldHighIntensityGreen  Color = Color(ansi.FG_BOLD_HIGH_INTENSITY_GREEN)
	FgBoldHighIntensityYellow Color = Color(ansi.FG_BOLD_HIGH_INTENSITY_YELLOW)
	FgBoldHighIntensityBlue   Color = Color(ansi.FG_BOLD_HIGH_INTENSITY_BLUE)
	FgBoldHighIntensityPurple Color = Color(ansi.FG_BOLD_HIGH_INTENSITY_PURPLE)
	FgBoldHighIntensityCyan   Color = Color(ansi.FG_BOLD_HIGH_INTENSITY_CYAN)
	FgBoldHighIntensityWhite  Color = Color(ansi.FG_BOLD_HIGH_INTENSITY_WHITE)

	BgBlack  Color = Color(ansi.BG_BLACK)
	BgRed    Color = Color(ansi.BG_RED)
	BgGreen  Color = Color(ansi.BG_GREEN)
	BgYellow Color = Color(ansi.BG_YELLOW)
	BgBlue   Color = Color(ansi.BG_BLUE)
	BgPurple Color = Color(ansi.BG_PURPLE)
	BgCyan   Color = Color(ansi.BG_CYAN)
	BgWhite  Color = Color(ansi.BG_WHITE)

	BgHighIntensityBlack  Color = Color(ansi.BG_HIGH_INTENSITY_BLACK)
	BgHighIntensityRed    Color = Color(ansi.BG_HIGH_INTENSITY_RED)
	BgHighIntensityGreen  Color = Color(ansi.BG_HIGH_INTENSITY_GREEN)
	BgHighIntensityYellow Color = Color(ansi.BG_HIGH_INTENSITY_YELLOW)
	BgHighIntensityBlue   Color = Color(ansi.BG_HIGH_INTENSITY_BLUE)
	BgHighIntensityPurple Color = Color(ansi.BG_HIGH_INTENSITY_PURPLE)
	BgHighIntensityCyan   Color = Color(ansi.BG_HIGH_INTENSITY_CYAN)
	BgHighIntensityWhite  Color = Color(ansi.BG_HIGH_INTENSITY_WHITE)
)
