package tuigo

import (
	"errors"
)

var (
	ErrInvalidFrameDimensions = errors.New("Invalid frame dimensions")
	ErrOutOfFrameBounds       = errors.New("Out of frame bounds")
)

type Frame struct {
	height, width int
	cells         []Cell
}

func NewFrame(width, height int, cells []Cell) (Frame, error) {
	if width <= 0 || height <= 0 {
		return Frame{}, ErrInvalidFrameDimensions
	}
	if len(cells) != width*height {
		return Frame{}, ErrInvalidFrameDimensions
	}
	return Frame{
		height: height,
		width:  width,
		cells:  cells,
	}, nil
}

func (f Frame) Height() int {
	return f.height
}

func (f Frame) Width() int {
	return f.width
}

func (f Frame) CellAt(x, y int) (Cell, error) {
	idx, err := f.idx(x, y)
	if err != nil {
		return Cell{}, err
	}
	return f.cells[idx], nil
}

func (f Frame) idx(x, y int) (int, error) {
	if x < 0 || y < 0 {
		return 0, ErrOutOfFrameBounds
	}
	if x >= f.width || y >= f.height {
		return 0, ErrOutOfFrameBounds
	}
	idx := f.width*y + x
	if idx >= len(f.cells) {
		return 0, ErrOutOfFrameBounds
	}
	return idx, nil
}
