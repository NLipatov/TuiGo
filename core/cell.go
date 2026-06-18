package core

import (
	"errors"

	"github.com/NLipatov/tuigo/color"
	"github.com/rivo/uniseg"
)

var (
	ErrEmptyCellGlyph           = errors.New("cell glyph should not be empty")
	ErrMultipleGraphemeClusters = errors.New("cell glyph should contain exactly one grapheme cluster")
	ErrUnsupportedCellWidth     = errors.New("cell glyph width should be one or two columns")
)

type Cell struct {
	glyph  string
	fg, bg color.Color
	width  uint8
}

func NewCell(glyph string, fg, bg color.Color) (Cell, error) {
	cluster, rest, width, _ := uniseg.FirstGraphemeClusterInString(glyph, -1)
	if cluster == "" {
		return Cell{}, ErrEmptyCellGlyph
	}
	if rest != "" {
		return Cell{}, ErrMultipleGraphemeClusters
	}
	if width != 1 && width != 2 {
		return Cell{}, ErrUnsupportedCellWidth
	}
	return NewCellWithWidth(cluster, width, fg, bg)
}

// NewCellWithWidth creates a cell from a caller-validated single grapheme cluster.
// The caller must ensure glyph is exactly one grapheme cluster and width is its display width.
// Passing multiple grapheme clusters or an incorrect width can corrupt frame rendering.
// It is a faster but unsafe alternative to NewCell.
func NewCellWithWidth(glyph string, width int, fg, bg color.Color) (Cell, error) {
	if glyph == "" {
		return Cell{}, ErrEmptyCellGlyph
	}
	if width != 1 && width != 2 {
		return Cell{}, ErrUnsupportedCellWidth
	}
	return Cell{
		glyph: glyph,
		fg:    fg,
		bg:    bg,
		width: uint8(width),
	}, nil
}

func (c Cell) Foreground() color.Color {
	return c.fg
}

func (c Cell) Background() color.Color {
	return c.bg
}

func (c Cell) Width() int {
	return int(c.width)
}

func (c Cell) Glyph() string {
	return c.glyph
}
