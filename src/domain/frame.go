package domain

import (
	"errors"
)

var (
	ErrOutOfFrameBounds = errors.New("Out of frame bounds")
)

type Frame struct {
	height, width uint
	cells         []Cell
}

func NewFrame(width, height uint) Frame {
	return Frame{
		height: height,
		width:  width,
		cells:  make([]Cell, width*height),
	}
}

func (f *Frame) Height() uint {
	return f.height
}

func (f *Frame) Width() uint {
	return f.width
}

func (f *Frame) CellAt(x, y uint) (Cell, error) {
	idx, err := f.idx(x, y)
	if err != nil {
		return Cell{}, err
	}
	return f.cells[idx], nil
}

func (f *Frame) idx(x, y uint) (uint, error) {
	if x >= f.width || y >= f.height {
		return 0, ErrOutOfFrameBounds
	}
	idx := f.width*y + x
	if idx >= uint(len(f.cells)) {
		return 0, ErrOutOfFrameBounds
	}
	return idx, nil
}
