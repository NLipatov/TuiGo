package render

import (
	"github.com/NLipatov/tuigo/color"
	"github.com/NLipatov/tuigo/internal/ansi"
)

var colorEscapes = [...]ansi.ANSIEscapeSequence{
	color.FgBlack:  ansi.FG_BLACK,
	color.FgRed:    ansi.FG_RED,
	color.FgGreen:  ansi.FG_GREEN,
	color.FgYellow: ansi.FG_YELLOW,
	color.FgBlue:   ansi.FG_BLUE,
	color.FgPurple: ansi.FG_PURPLE,
	color.FgCyan:   ansi.FG_CYAN,
	color.FgWhite:  ansi.FG_WHITE,

	color.FgBoldBlack:  ansi.FG_BOLD_BLACK,
	color.FgBoldRed:    ansi.FG_BOLD_RED,
	color.FgBoldGreen:  ansi.FG_BOLD_GREEN,
	color.FgBoldYellow: ansi.FG_BOLD_YELLOW,
	color.FgBoldBlue:   ansi.FG_BOLD_BLUE,
	color.FgBoldPurple: ansi.FG_BOLD_PURPLE,
	color.FgBoldCyan:   ansi.FG_BOLD_CYAN,
	color.FgBoldWhite:  ansi.FG_BOLD_WHITE,

	color.FgUnderlineBlack:  ansi.FG_UNDERLINE_BLACK,
	color.FgUnderlineRed:    ansi.FG_UNDERLINE_RED,
	color.FgUnderlineGreen:  ansi.FG_UNDERLINE_GREEN,
	color.FgUnderlineYellow: ansi.FG_UNDERLINE_YELLOW,
	color.FgUnderlineBlue:   ansi.FG_UNDERLINE_BLUE,
	color.FgUnderlinePurple: ansi.FG_UNDERLINE_PURPLE,
	color.FgUnderlineCyan:   ansi.FG_UNDERLINE_CYAN,
	color.FgUnderlineWhite:  ansi.FG_UNDERLINE_WHITE,

	color.FgHighIntensityBlack:  ansi.FG_HIGH_INTENSITY_BLACK,
	color.FgHighIntensityRed:    ansi.FG_HIGH_INTENSITY_RED,
	color.FgHighIntensityGreen:  ansi.FG_HIGH_INTENSITY_GREEN,
	color.FgHighIntensityYellow: ansi.FG_HIGH_INTENSITY_YELLOW,
	color.FgHighIntensityBlue:   ansi.FG_HIGH_INTENSITY_BLUE,
	color.FgHighIntensityPurple: ansi.FG_HIGH_INTENSITY_PURPLE,
	color.FgHighIntensityCyan:   ansi.FG_HIGH_INTENSITY_CYAN,
	color.FgHighIntensityWhite:  ansi.FG_HIGH_INTENSITY_WHITE,

	color.FgBoldHighIntensityBlack:  ansi.FG_BOLD_HIGH_INTENSITY_BLACK,
	color.FgBoldHighIntensityRed:    ansi.FG_BOLD_HIGH_INTENSITY_RED,
	color.FgBoldHighIntensityGreen:  ansi.FG_BOLD_HIGH_INTENSITY_GREEN,
	color.FgBoldHighIntensityYellow: ansi.FG_BOLD_HIGH_INTENSITY_YELLOW,
	color.FgBoldHighIntensityBlue:   ansi.FG_BOLD_HIGH_INTENSITY_BLUE,
	color.FgBoldHighIntensityPurple: ansi.FG_BOLD_HIGH_INTENSITY_PURPLE,
	color.FgBoldHighIntensityCyan:   ansi.FG_BOLD_HIGH_INTENSITY_CYAN,
	color.FgBoldHighIntensityWhite:  ansi.FG_BOLD_HIGH_INTENSITY_WHITE,

	color.BgBlack:  ansi.BG_BLACK,
	color.BgRed:    ansi.BG_RED,
	color.BgGreen:  ansi.BG_GREEN,
	color.BgYellow: ansi.BG_YELLOW,
	color.BgBlue:   ansi.BG_BLUE,
	color.BgPurple: ansi.BG_PURPLE,
	color.BgCyan:   ansi.BG_CYAN,
	color.BgWhite:  ansi.BG_WHITE,

	color.BgHighIntensityBlack:  ansi.BG_HIGH_INTENSITY_BLACK,
	color.BgHighIntensityRed:    ansi.BG_HIGH_INTENSITY_RED,
	color.BgHighIntensityGreen:  ansi.BG_HIGH_INTENSITY_GREEN,
	color.BgHighIntensityYellow: ansi.BG_HIGH_INTENSITY_YELLOW,
	color.BgHighIntensityBlue:   ansi.BG_HIGH_INTENSITY_BLUE,
	color.BgHighIntensityPurple: ansi.BG_HIGH_INTENSITY_PURPLE,
	color.BgHighIntensityCyan:   ansi.BG_HIGH_INTENSITY_CYAN,
	color.BgHighIntensityWhite:  ansi.BG_HIGH_INTENSITY_WHITE,
}

func colorEscape(c color.Color) ansi.ANSIEscapeSequence {
	if int(c) >= len(colorEscapes) {
		return ""
	}
	return colorEscapes[c]
}
