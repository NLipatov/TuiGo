package core

import (
	"errors"
	"testing"

	"github.com/NLipatov/tuigo/ansi"
)

func TestNewCellAcceptsSingleGraphemeCluster(t *testing.T) {
	fg, bg := testColors(t)

	tests := []struct {
		name  string
		glyph string
		width int
	}{
		{
			name:  "ascii",
			glyph: "x",
			width: 1,
		},
		{
			name:  "combining mark cluster",
			glyph: "a\u0301",
			width: 1,
		},
		{
			name:  "emoji",
			glyph: "🙂",
			width: 2,
		},
		{
			name:  "emoji modifier cluster",
			glyph: "👍🏽",
			width: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell, err := NewCell(tt.glyph, fg, bg)
			if err != nil {
				t.Fatalf("NewCell(%q) error = %v", tt.glyph, err)
			}
			if cell.Glyph() != tt.glyph {
				t.Fatalf("Glyph() = %q, want %q", cell.Glyph(), tt.glyph)
			}
			if cell.Width() != tt.width {
				t.Fatalf("Width() = %d, want %d", cell.Width(), tt.width)
			}
			if cell.Foreground() != fg {
				t.Fatalf("Foreground() = %#v, want %#v", cell.Foreground(), fg)
			}
			if cell.Background() != bg {
				t.Fatalf("Background() = %#v, want %#v", cell.Background(), bg)
			}
		})
	}
}

func TestNewCellRejectsInvalidGlyph(t *testing.T) {
	fg, bg := testColors(t)

	tests := []struct {
		name  string
		glyph string
		want  error
	}{
		{
			name:  "empty glyph",
			glyph: "",
			want:  ErrEmptyCellGlyph,
		},
		{
			name:  "multiple grapheme clusters",
			glyph: "ab",
			want:  ErrMultipleGraphemeClusters,
		},
		{
			name:  "zero-width cluster",
			glyph: "\u0301",
			want:  ErrUnsupportedCellWidth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCell(tt.glyph, fg, bg)
			if !errors.Is(err, tt.want) {
				t.Fatalf("NewCell(%q) error = %v, want %v", tt.glyph, err, tt.want)
			}
		})
	}
}

func testColors(t *testing.T) (ansi.Color, ansi.Color) {
	t.Helper()

	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}
	return fg, bg
}
