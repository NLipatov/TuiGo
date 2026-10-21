package compose

import (
	"errors"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
)

var (
	ErrNegativePadding      = errors.New("negative padding")
	ErrUnsupportedDirection = errors.New("unsupported direction")
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
	if d.Padding.Up < 0 || d.Padding.Down < 0 ||
		d.Padding.Left < 0 || d.Padding.Right < 0 {
		return canvas.Frame{}, ErrNegativePadding
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
	if err := d.renderChildren(cells); err != nil {
		return canvas.Frame{}, err
	}
	return canvas.NewFrame(d.Width, d.Height, cells)
}

func (d *Div) renderChildren(cells []canvas.Cell) error {
	offsetY, offsetX := 0+d.Padding.Up, 0+d.Padding.Left
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
			for x := 0; x < len(row); x += row[x].Width() {
				cell := row[x]
				if cell.Foreground() == color.FgTransparent ||
					cell.Background() == color.BgTransparent {
					continue
				}
				if x+offsetX+cell.Width() > d.Width-d.Padding.Right {
					break
				}
				idx := (y+offsetY)*d.Width + x + d.Padding.Left
				if idx >= len(cells) {
					break
				}
				cells[idx] = cell
			}
		}
		switch d.Direction {
		case Row:
			offsetX += frame.Width()
		case Column:
			offsetY += frame.Height()
		default:
			return ErrUnsupportedDirection
		}
	}
	return nil
}
