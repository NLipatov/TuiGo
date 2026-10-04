package compose

import (
	"errors"
	"math"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
	"github.com/rivo/uniseg"
)

var (
	ErrInvalidDivDimensions = errors.New("invalid div dimensions")
	ErrEmptyText            = errors.New("text element should have non-empty content")
	ErrDivAreaOverflow      = errors.New("given div width and height multiplication results in overflow")
)

type Node interface {
	Frame() (canvas.Frame, error)
}

type Direction int

const (
	Row Direction = iota
	Column
)

type Div struct {
	Width, Height int
	Fg, Bg        color.Color
	Children      []Node
	Direction     Direction
}

func (d *Div) Frame() (canvas.Frame, error) {
	if d.Width <= 0 || d.Height <= 0 {
		return canvas.Frame{}, ErrInvalidDivDimensions
	}
	if d.Width > math.MaxInt/d.Height {
		return canvas.Frame{}, ErrDivAreaOverflow
	}
	cells := make([]canvas.Cell, d.Width*d.Height)
	for i := range cells {
		cell, err := canvas.NewCell(" ", d.Fg, d.Bg)
		if err != nil {
			return canvas.Frame{}, err
		}
		cells[i] = cell
	}
	switch d.Direction {
	case Row:
		err := d.renderRow(cells)
		if err != nil {
			return canvas.Frame{}, err
		}
	case Column:
		err := d.renderColumn(cells)
		if err != nil {
			return canvas.Frame{}, err
		}
	}
	return canvas.NewFrame(d.Width, d.Height, cells)
}

func (d *Div) renderRow(cells []canvas.Cell) error {
	offsetX := 0
	for _, child := range d.Children {
		frame, err := child.Frame()
		if err != nil {
			return err
		}
		for y := range frame.Height() {
			if y >= d.Height {
				break
			}
			row, err := frame.RowAt(y)
			if err != nil {
				return err
			}
			for x, cell := range row {
				if x+offsetX >= d.Width {
					break
				}
				if x+offsetX+cell.Width() > d.Width {
					break
				}
				idx := y*d.Width + (x + offsetX)
				if idx >= len(cells) {
					break
				}
				cells[idx] = cell
			}
		}
		offsetX += frame.Width()
	}
	return nil
}

func (d *Div) renderColumn(cells []canvas.Cell) error {
	offsetY := 0
	for _, child := range d.Children {
		frame, err := child.Frame()
		if err != nil {
			return err
		}
		for y := range frame.Height() {
			if y+offsetY >= d.Height {
				break
			}
			row, err := frame.RowAt(y)
			if err != nil {
				return err
			}
			for x, cell := range row {
				if x >= d.Width {
					break
				}
				if x+cell.Width() > d.Width {
					break
				}
				idx := (y+offsetY)*d.Width + x
				if idx >= len(cells) {
					break
				}
				cells[idx] = cell
			}
		}
		offsetY += frame.Height()
	}
	return nil
}

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
