package core

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidFrameDimensions = errors.New("invalid frame dimensions")
	ErrOutOfFrameBounds       = errors.New("out of frame bounds")
)

type FrameCellLayoutError struct {
	X, Y   int
	Reason string
}

func (f FrameCellLayoutError) Error() string {
	return fmt.Sprintf("invalid frame cell layout: %s. At x: %d, y: %d", f.Reason, f.X, f.Y)
}

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
	frame := Frame{
		height: height,
		width:  width,
		cells:  cells,
	}
	for y := range height {
		for x := 0; x < width; x++ {
			cell, err := frame.CellAt(x, y)
			if err != nil {
				return Frame{}, err
			}
			switch cell.Width() {
			case 0:
				return Frame{}, FrameCellLayoutError{
					X:      x,
					Y:      y,
					Reason: "unexpected continuation block",
				}
			case 2:
				if x+1 == width {
					return Frame{}, FrameCellLayoutError{
						X:      x,
						Y:      y,
						Reason: "missing continuation block",
					}
				}
				continuation, err := frame.CellAt(x+1, y)
				if err != nil {
					return Frame{}, err
				}
				if continuation.Width() != 0 {
					return Frame{}, FrameCellLayoutError{
						X:      x,
						Y:      y,
						Reason: "missing continuation block",
					}
				}
				x += cell.Width() - 1
			}
		}
	}
	return frame, nil
}

func (f Frame) Height() int {
	return f.height
}

func (f Frame) Width() int {
	return f.width
}

// RowAt returns a mutable, zero-copy view of row y backed by the frame's cell buffer.
func (f Frame) RowAt(y int) ([]Cell, error) {
	if y >= f.height || y < 0 {
		return nil, ErrOutOfFrameBounds
	}
	start := f.width * y
	end := start + f.width
	return f.cells[start:end:end], nil
}

// CellAt returns a copy of the cell at x,y.
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
