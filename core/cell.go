package core

import (
	"errors"

	"github.com/NLipatov/tuigo/ansi"
	"github.com/rivo/uniseg"
)

var (
	ErrEmptyCellGlyph           = errors.New("cell glyph should not be empty")
	ErrMultipleGraphemeClusters = errors.New("cell glyph should contain exactly one grapheme cluster")
	ErrUnsupportedCellWidth     = errors.New("cell glyph width should be one or two columns")
)

type Cell struct {
	fg, bg ansi.Color
	glyph  string
	width  uint8
}

func NewCell(glyph string, fg, bg ansi.Color) (Cell, error) {
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
	return Cell{
		glyph: cluster,
		fg:    fg,
		bg:    bg,
		width: uint8(width),
	}, nil
}

func (c Cell) Foreground() ansi.Color {
	return c.fg
}

func (c Cell) Background() ansi.Color {
	return c.bg
}

func (c Cell) Width() int {
	return int(c.width)
}

func (c Cell) Glyph() string {
	return c.glyph
}
