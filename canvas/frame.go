package canvas

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
		for x := 0; x < frame.width; x += row[x].Width() {
			cell := row[x]
			if cell.Width() == 0 {
				return FrameCellLayoutError{
					X:      x,
					Y:      y,
					Reason: "unexpected continuation block",
				}
			}
			head, end := x+1, x+cell.Width()
			if end > frame.width {
				return FrameCellLayoutError{
					X:      x,
					Y:      y,
					Reason: "could not fit grapheme on frame: too narrow frame",
				}
			}
			for head < end {
				if row[head].Width() != 0 {
					return FrameCellLayoutError{
						X:      x,
						Y:      y,
						Reason: "missing continuation block",
					}
				}
				head++
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
