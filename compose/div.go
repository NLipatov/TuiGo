package compose

import (
	"errors"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
)

var (
	ErrNegativePadding = errors.New("negative padding")
)

type Direction int

const (
	Row Direction = iota
	Column
)

type Padding struct {
	Up, Down, Left, Right int
}

type Div struct {
	Width, Height int
	Padding       Padding
	Fg, Bg        color.Color
	Children      []Node
	Direction     Direction
	cells         []canvas.Cell
}

func (d *Div) Frame() (canvas.Frame, error) {
	if err := validateDimensions(d.Width, d.Height); err != nil {
		return canvas.Frame{}, err
	}
	if err := d.validatePadding(); err != nil {
		return canvas.Frame{}, err
	}
	if len(d.cells) != d.Width*d.Height {
		d.cells = make([]canvas.Cell, d.Width*d.Height)
	}
	cells := d.cells
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

func (d *Div) validatePadding() error {
	if d.Padding.Up < 0 || d.Padding.Down < 0 ||
		d.Padding.Left < 0 || d.Padding.Right < 0 {
		return ErrNegativePadding
	}
	return nil
}

func (d *Div) renderRow(cells []canvas.Cell) error {
	offsetX := 0 + d.Padding.Left
	for _, child := range d.Children {
		frame, err := child.Frame()
		if err != nil {
			return err
		}
		for y := range frame.Height() {
			if y+d.Padding.Up >= d.Height-d.Padding.Down {
				break
			}
			row, err := frame.RowAt(y)
			if err != nil {
				return err
			}
			for x, cell := range row {
				if x+offsetX+cell.Width() > d.Width-d.Padding.Right {
					break
				}
				idx := (y+d.Padding.Up)*d.Width + (x + offsetX)
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
	offsetY := 0 + d.Padding.Up
	for _, child := range d.Children {
		frame, err := child.Frame()
		if err != nil {
			return err
		}
		for y := range frame.Height() {
			if y+offsetY >= d.Height-d.Padding.Down {
				break
			}
			row, err := frame.RowAt(y)
			if err != nil {
				return err
			}
			for x, cell := range row {
				if x+d.Padding.Left+cell.Width() > d.Width-d.Padding.Right {
					break
				}
				idx := (y+offsetY)*d.Width + x + d.Padding.Left
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
