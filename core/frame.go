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
	if err := frame.validateCellLayout(); err != nil {
		return Frame{}, err
	}
	return frame, nil
}

func (frame Frame) validateCellLayout() error {
	for y := range frame.height {
		row, err := frame.RowAt(y)
		if err != nil {
			return err
		}
		for x := 0; x < frame.width; x++ {
			cell := row[x]
			switch cell.Width() {
			case 0:
				return FrameCellLayoutError{
					X:      x,
					Y:      y,
					Reason: "unexpected continuation block",
				}
			case 2:
				if x+1 == frame.width {
					return FrameCellLayoutError{
						X:      x,
						Y:      y,
						Reason: "missing continuation block",
					}
				}
				continuation := row[x+1]
				if continuation.Width() != 0 {
					return FrameCellLayoutError{
						X:      x,
						Y:      y,
						Reason: "missing continuation block",
					}
				}
				x += cell.Width() - 1
			}
		}
	}
	return nil
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
