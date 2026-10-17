package compose

import (
	"errors"
	"testing"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
)

func TestTextFramePreservesGraphemes(t *testing.T) {
	text := &Text{Content: "A🙂e\u0301", Fg: color.FgRed, Bg: color.BgBlue}
	frame, err := text.Frame()
	if err != nil {
		t.Fatalf("Frame() error = %v", err)
	}
	if frame.Width() != 4 || frame.Height() != 1 {
		t.Fatalf("Frame() size = %dx%d, want 4x1", frame.Width(), frame.Height())
	}
	row, err := frame.RowAt(0)
	if err != nil {
		t.Fatalf("RowAt(0) error = %v", err)
	}
	for _, tt := range []struct {
		x, width int
		glyph    string
	}{
		{x: 0, width: 1, glyph: "A"},
		{x: 1, width: 2, glyph: "🙂"},
		{x: 3, width: 1, glyph: "e\u0301"},
	} {
		cell := row[tt.x]
		if cell.Glyph() != tt.glyph || cell.Width() != tt.width ||
			cell.Foreground() != color.FgRed || cell.Background() != color.BgBlue {
			t.Errorf("cell %d = %#v, want glyph %q, width %d, red on blue", tt.x, cell, tt.glyph, tt.width)
		}
	}
	if row[2] != (canvas.Cell{}) {
		t.Errorf("cell 2 = %#v, want continuation block", row[2])
	}
}

func TestTextRejectsEmptyContent(t *testing.T) {
	_, err := (&Text{}).Frame()
	if !errors.Is(err, ErrEmptyText) {
		t.Fatalf("Frame() error = %v, want %v", err, ErrEmptyText)
	}
}
