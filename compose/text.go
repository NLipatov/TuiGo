package compose

import (
	"errors"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
	"github.com/rivo/uniseg"
)

var (
	ErrEmptyText = errors.New("text element should have non-empty content")
)

type Text struct {
	Content string
	Fg, Bg  color.Color
	cells   []canvas.Cell
}

func (t *Text) Frame() (canvas.Frame, error) {
	if len(t.Content) == 0 {
		return canvas.Frame{}, ErrEmptyText
	}
	gs := uniseg.NewGraphemes(t.Content)
	cells := t.cells[:0]
	if cap(cells) < len(t.Content) {
		cells = make([]canvas.Cell, 0, len(t.Content))
	}
	for gs.Next() {
		cell, err := canvas.NewCellWithWidth(gs.Str(), gs.Width(), t.Fg, t.Bg)
		if err != nil {
			return canvas.Frame{}, err
		}
		cells = append(cells, cell)
		for range gs.Width() - 1 {
			cells = append(cells, canvas.Cell{})
		}
	}
	t.cells = cells
	return canvas.NewFrame(len(cells), 1, cells)
}
