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
}

func (t *Text) Frame() (canvas.Frame, error) {
	if len(t.Content) == 0 {
		return canvas.Frame{}, ErrEmptyText
	}
	gs := uniseg.NewGraphemes(t.Content)
	cells := make([]canvas.Cell, 0, len(t.Content))
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
	return canvas.NewFrame(len(cells), 1, cells)
}
