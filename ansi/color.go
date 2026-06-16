package ansi

import "errors"

var (
	ErrInvalidEscapeSequence = errors.New("invalid Escape sequence")
)

type colorCode uint8

const (
	colorInvalid colorCode = iota
	colorFGBlack
	colorFGRed
	colorFGGreen
	colorFGYellow
	colorFGBlue
	colorFGPurple
	colorFGCyan
	colorFGWhite
	colorFGBoldBlack
	colorFGBoldRed
	colorFGBoldGreen
	colorFGBoldYellow
	colorFGBoldBlue
	colorFGBoldPurple
	colorFGBoldCyan
	colorFGBoldWhite
	colorFGUnderlineBlack
	colorFGUnderlineRed
	colorFGUnderlineGreen
	colorFGUnderlineYellow
	colorFGUnderlineBlue
	colorFGUnderlinePurple
	colorFGUnderlineCyan
	colorFGUnderlineWhite
	colorBGBlack
	colorBGRed
	colorBGGreen
	colorBGYellow
	colorBGBlue
	colorBGPurple
	colorBGCyan
	colorBGWhite
	colorFGHighIntensityBlack
	colorFGHighIntensityRed
	colorFGHighIntensityGreen
	colorFGHighIntensityYellow
	colorFGHighIntensityBlue
	colorFGHighIntensityPurple
	colorFGHighIntensityCyan
	colorFGHighIntensityWhite
	colorFGBoldHighIntensityBlack
	colorFGBoldHighIntensityRed
	colorFGBoldHighIntensityGreen
	colorFGBoldHighIntensityYellow
	colorFGBoldHighIntensityBlue
	colorFGBoldHighIntensityPurple
	colorFGBoldHighIntensityCyan
	colorFGBoldHighIntensityWhite
	colorBGHighIntensityBlack
	colorBGHighIntensityRed
	colorBGHighIntensityGreen
	colorBGHighIntensityYellow
	colorBGHighIntensityBlue
	colorBGHighIntensityPurple
	colorBGHighIntensityCyan
	colorBGHighIntensityWhite
)

var colorEscapeSequences = [...]ANSIEscapeSequence{
	colorFGBlack:                   FG_BLACK,
	colorFGRed:                     FG_RED,
	colorFGGreen:                   FG_GREEN,
	colorFGYellow:                  FG_YELLOW,
	colorFGBlue:                    FG_BLUE,
	colorFGPurple:                  FG_PURPLE,
	colorFGCyan:                    FG_CYAN,
	colorFGWhite:                   FG_WHITE,
	colorFGBoldBlack:               FG_BOLD_BLACK,
	colorFGBoldRed:                 FG_BOLD_RED,
	colorFGBoldGreen:               FG_BOLD_GREEN,
	colorFGBoldYellow:              FG_BOLD_YELLOW,
	colorFGBoldBlue:                FG_BOLD_BLUE,
	colorFGBoldPurple:              FG_BOLD_PURPLE,
	colorFGBoldCyan:                FG_BOLD_CYAN,
	colorFGBoldWhite:               FG_BOLD_WHITE,
	colorFGUnderlineBlack:          FG_UNDERLINE_BLACK,
	colorFGUnderlineRed:            FG_UNDERLINE_RED,
	colorFGUnderlineGreen:          FG_UNDERLINE_GREEN,
	colorFGUnderlineYellow:         FG_UNDERLINE_YELLOW,
	colorFGUnderlineBlue:           FG_UNDERLINE_BLUE,
	colorFGUnderlinePurple:         FG_UNDERLINE_PURPLE,
	colorFGUnderlineCyan:           FG_UNDERLINE_CYAN,
	colorFGUnderlineWhite:          FG_UNDERLINE_WHITE,
	colorBGBlack:                   BG_BLACK,
	colorBGRed:                     BG_RED,
	colorBGGreen:                   BG_GREEN,
	colorBGYellow:                  BG_YELLOW,
	colorBGBlue:                    BG_BLUE,
	colorBGPurple:                  BG_PURPLE,
	colorBGCyan:                    BG_CYAN,
	colorBGWhite:                   BG_WHITE,
	colorFGHighIntensityBlack:      FG_HIGH_INTENSITY_BLACK,
	colorFGHighIntensityRed:        FG_HIGH_INTENSITY_RED,
	colorFGHighIntensityGreen:      FG_HIGH_INTENSITY_GREEN,
	colorFGHighIntensityYellow:     FG_HIGH_INTENSITY_YELLOW,
	colorFGHighIntensityBlue:       FG_HIGH_INTENSITY_BLUE,
	colorFGHighIntensityPurple:     FG_HIGH_INTENSITY_PURPLE,
	colorFGHighIntensityCyan:       FG_HIGH_INTENSITY_CYAN,
	colorFGHighIntensityWhite:      FG_HIGH_INTENSITY_WHITE,
	colorFGBoldHighIntensityBlack:  FG_BOLD_HIGH_INTENSITY_BLACK,
	colorFGBoldHighIntensityRed:    FG_BOLD_HIGH_INTENSITY_RED,
	colorFGBoldHighIntensityGreen:  FG_BOLD_HIGH_INTENSITY_GREEN,
	colorFGBoldHighIntensityYellow: FG_BOLD_HIGH_INTENSITY_YELLOW,
	colorFGBoldHighIntensityBlue:   FG_BOLD_HIGH_INTENSITY_BLUE,
	colorFGBoldHighIntensityPurple: FG_BOLD_HIGH_INTENSITY_PURPLE,
	colorFGBoldHighIntensityCyan:   FG_BOLD_HIGH_INTENSITY_CYAN,
	colorFGBoldHighIntensityWhite:  FG_BOLD_HIGH_INTENSITY_WHITE,
	colorBGHighIntensityBlack:      BG_HIGH_INTENSITY_BLACK,
	colorBGHighIntensityRed:        BG_HIGH_INTENSITY_RED,
	colorBGHighIntensityGreen:      BG_HIGH_INTENSITY_GREEN,
	colorBGHighIntensityYellow:     BG_HIGH_INTENSITY_YELLOW,
	colorBGHighIntensityBlue:       BG_HIGH_INTENSITY_BLUE,
	colorBGHighIntensityPurple:     BG_HIGH_INTENSITY_PURPLE,
	colorBGHighIntensityCyan:       BG_HIGH_INTENSITY_CYAN,
	colorBGHighIntensityWhite:      BG_HIGH_INTENSITY_WHITE,
}

type Color struct {
	code colorCode
}

func NewColor(escapeSequence ANSIEscapeSequence) (Color, error) {
	code, ok := colorCodeForEscapeSequence(escapeSequence)
	if !ok {
		return Color{}, ErrInvalidEscapeSequence
	}
	return Color{
		code: code,
	}, nil
}

func (c Color) String() string {
	idx := int(c.code)
	if idx <= 0 || idx >= len(colorEscapeSequences) {
		return ""
	}
	return string(colorEscapeSequences[idx])
}

//nolint:cyclop // mechanical ANSI color lookup table; avoids init-time map allocation.
func colorCodeForEscapeSequence(escapeSequence ANSIEscapeSequence) (colorCode, bool) {
	switch escapeSequence {
	case FG_BLACK:
		return colorFGBlack, true
	case FG_RED:
		return colorFGRed, true
	case FG_GREEN:
		return colorFGGreen, true
	case FG_YELLOW:
		return colorFGYellow, true
	case FG_BLUE:
		return colorFGBlue, true
	case FG_PURPLE:
		return colorFGPurple, true
	case FG_CYAN:
		return colorFGCyan, true
	case FG_WHITE:
		return colorFGWhite, true
	case FG_BOLD_BLACK:
		return colorFGBoldBlack, true
	case FG_BOLD_RED:
		return colorFGBoldRed, true
	case FG_BOLD_GREEN:
		return colorFGBoldGreen, true
	case FG_BOLD_YELLOW:
		return colorFGBoldYellow, true
	case FG_BOLD_BLUE:
		return colorFGBoldBlue, true
	case FG_BOLD_PURPLE:
		return colorFGBoldPurple, true
	case FG_BOLD_CYAN:
		return colorFGBoldCyan, true
	case FG_BOLD_WHITE:
		return colorFGBoldWhite, true
	case FG_UNDERLINE_BLACK:
		return colorFGUnderlineBlack, true
	case FG_UNDERLINE_RED:
		return colorFGUnderlineRed, true
	case FG_UNDERLINE_GREEN:
		return colorFGUnderlineGreen, true
	case FG_UNDERLINE_YELLOW:
		return colorFGUnderlineYellow, true
	case FG_UNDERLINE_BLUE:
		return colorFGUnderlineBlue, true
	case FG_UNDERLINE_PURPLE:
		return colorFGUnderlinePurple, true
	case FG_UNDERLINE_CYAN:
		return colorFGUnderlineCyan, true
	case FG_UNDERLINE_WHITE:
		return colorFGUnderlineWhite, true
	case BG_BLACK:
		return colorBGBlack, true
	case BG_RED:
		return colorBGRed, true
	case BG_GREEN:
		return colorBGGreen, true
	case BG_YELLOW:
		return colorBGYellow, true
	case BG_BLUE:
		return colorBGBlue, true
	case BG_PURPLE:
		return colorBGPurple, true
	case BG_CYAN:
		return colorBGCyan, true
	case BG_WHITE:
		return colorBGWhite, true
	case FG_HIGH_INTENSITY_BLACK:
		return colorFGHighIntensityBlack, true
	case FG_HIGH_INTENSITY_RED:
		return colorFGHighIntensityRed, true
	case FG_HIGH_INTENSITY_GREEN:
		return colorFGHighIntensityGreen, true
	case FG_HIGH_INTENSITY_YELLOW:
		return colorFGHighIntensityYellow, true
	case FG_HIGH_INTENSITY_BLUE:
		return colorFGHighIntensityBlue, true
	case FG_HIGH_INTENSITY_PURPLE:
		return colorFGHighIntensityPurple, true
	case FG_HIGH_INTENSITY_CYAN:
		return colorFGHighIntensityCyan, true
	case FG_HIGH_INTENSITY_WHITE:
		return colorFGHighIntensityWhite, true
	case FG_BOLD_HIGH_INTENSITY_BLACK:
		return colorFGBoldHighIntensityBlack, true
	case FG_BOLD_HIGH_INTENSITY_RED:
		return colorFGBoldHighIntensityRed, true
	case FG_BOLD_HIGH_INTENSITY_GREEN:
		return colorFGBoldHighIntensityGreen, true
	case FG_BOLD_HIGH_INTENSITY_YELLOW:
		return colorFGBoldHighIntensityYellow, true
	case FG_BOLD_HIGH_INTENSITY_BLUE:
		return colorFGBoldHighIntensityBlue, true
	case FG_BOLD_HIGH_INTENSITY_PURPLE:
		return colorFGBoldHighIntensityPurple, true
	case FG_BOLD_HIGH_INTENSITY_CYAN:
		return colorFGBoldHighIntensityCyan, true
	case FG_BOLD_HIGH_INTENSITY_WHITE:
		return colorFGBoldHighIntensityWhite, true
	case BG_HIGH_INTENSITY_BLACK:
		return colorBGHighIntensityBlack, true
	case BG_HIGH_INTENSITY_RED:
		return colorBGHighIntensityRed, true
	case BG_HIGH_INTENSITY_GREEN:
		return colorBGHighIntensityGreen, true
	case BG_HIGH_INTENSITY_YELLOW:
		return colorBGHighIntensityYellow, true
	case BG_HIGH_INTENSITY_BLUE:
		return colorBGHighIntensityBlue, true
	case BG_HIGH_INTENSITY_PURPLE:
		return colorBGHighIntensityPurple, true
	case BG_HIGH_INTENSITY_CYAN:
		return colorBGHighIntensityCyan, true
	case BG_HIGH_INTENSITY_WHITE:
		return colorBGHighIntensityWhite, true
	default:
		return colorInvalid, false
	}
}
