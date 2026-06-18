package core

import (
	"errors"
	"testing"

	"github.com/NLipatov/tuigo/color"
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

func TestNewCellWithWidthAcceptsKnownWidth(t *testing.T) {
	fg, bg := testColors(t)

	tests := []struct {
		glyph string
		width int
	}{
		{glyph: "Ж", width: 1},
		{glyph: "界", width: 2},
	}

	for _, tt := range tests {
		cell, err := NewCellWithWidth(tt.glyph, tt.width, fg, bg)
		if err != nil {
			t.Fatalf("NewCellWithWidth() error = %v", err)
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
	}
}

func TestNewCellWithWidthRejectsInvalidInput(t *testing.T) {
	fg, bg := testColors(t)

	tests := []struct {
		name  string
		glyph string
		width int
		want  error
	}{
		{
			name:  "empty glyph",
			glyph: "",
			width: 1,
			want:  ErrEmptyCellGlyph,
		},
		{
			name:  "zero width",
			glyph: "x",
			width: 0,
			want:  ErrUnsupportedCellWidth,
		},
		{
			name:  "unsupported width",
			glyph: "x",
			width: 3,
			want:  ErrUnsupportedCellWidth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCellWithWidth(tt.glyph, tt.width, fg, bg)
			if !errors.Is(err, tt.want) {
				t.Fatalf("NewCellWithWidth(%q, %d) error = %v, want %v", tt.glyph, tt.width, err, tt.want)
			}
		})
	}
}

func testColors(t *testing.T) (color.Color, color.Color) {
	t.Helper()

	return color.FgRed, color.BgBlack
}
