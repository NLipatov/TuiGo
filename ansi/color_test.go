package ansi

import (
	"errors"
	"testing"
)

func TestNewColorAcceptsColors(t *testing.T) {
	tests := []ANSIEscapeSequence{
		FG_BLACK,
		FG_RED,
		FG_GREEN,
		FG_YELLOW,
		FG_BLUE,
		FG_PURPLE,
		FG_CYAN,
		FG_WHITE,
		FG_BOLD_BLACK,
		FG_BOLD_RED,
		FG_BOLD_GREEN,
		FG_BOLD_YELLOW,
		FG_BOLD_BLUE,
		FG_BOLD_PURPLE,
		FG_BOLD_CYAN,
		FG_BOLD_WHITE,
		FG_UNDERLINE_BLACK,
		FG_UNDERLINE_RED,
		FG_UNDERLINE_GREEN,
		FG_UNDERLINE_YELLOW,
		FG_UNDERLINE_BLUE,
		FG_UNDERLINE_PURPLE,
		FG_UNDERLINE_CYAN,
		FG_UNDERLINE_WHITE,
		BG_BLACK,
		BG_RED,
		BG_GREEN,
		BG_YELLOW,
		BG_BLUE,
		BG_PURPLE,
		BG_CYAN,
		BG_WHITE,
		FG_HIGH_INTENSITY_BLACK,
		FG_HIGH_INTENSITY_RED,
		FG_HIGH_INTENSITY_GREEN,
		FG_HIGH_INTENSITY_YELLOW,
		FG_HIGH_INTENSITY_BLUE,
		FG_HIGH_INTENSITY_PURPLE,
		FG_HIGH_INTENSITY_CYAN,
		FG_HIGH_INTENSITY_WHITE,
		FG_BOLD_HIGH_INTENSITY_BLACK,
		FG_BOLD_HIGH_INTENSITY_RED,
		FG_BOLD_HIGH_INTENSITY_GREEN,
		FG_BOLD_HIGH_INTENSITY_YELLOW,
		FG_BOLD_HIGH_INTENSITY_BLUE,
		FG_BOLD_HIGH_INTENSITY_PURPLE,
		FG_BOLD_HIGH_INTENSITY_CYAN,
		FG_BOLD_HIGH_INTENSITY_WHITE,
		BG_HIGH_INTENSITY_BLACK,
		BG_HIGH_INTENSITY_RED,
		BG_HIGH_INTENSITY_GREEN,
		BG_HIGH_INTENSITY_YELLOW,
		BG_HIGH_INTENSITY_BLUE,
		BG_HIGH_INTENSITY_PURPLE,
		BG_HIGH_INTENSITY_CYAN,
		BG_HIGH_INTENSITY_WHITE,
	}

	for _, sequence := range tests {
		t.Run(string(sequence), func(t *testing.T) {
			color, err := NewColor(sequence)
			if err != nil {
				t.Fatalf("NewColor(%q) error = %v", sequence, err)
			}
			if got := color.String(); got != string(sequence) {
				t.Fatalf("String() = %q, want %q", got, sequence)
			}
		})
	}
}

func TestNewColorRejectsNonColors(t *testing.T) {
	tests := []ANSIEscapeSequence{
		"",
		RESET,
		CLEAR_SCREEN,
		ENTER_ALTERNATE_SCREEN,
		ENABLE_MOUSE_REPORTING,
	}

	for _, sequence := range tests {
		t.Run(string(sequence), func(t *testing.T) {
			_, err := NewColor(sequence)
			if !errors.Is(err, ErrInvalidEscapeSequence) {
				t.Fatalf("NewColor(%q) error = %v, want %v", sequence, err, ErrInvalidEscapeSequence)
			}
		})
	}
}

func TestZeroColorStringIsEmpty(t *testing.T) {
	if got := (Color{}).String(); got != "" {
		t.Fatalf("String() = %q, want empty string", got)
	}
}
