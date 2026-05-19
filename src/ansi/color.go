package ansi

import "errors"

var (
	ErrInvalidEscapeSequence = errors.New("Invalid Escape sequence")
)

type Color struct {
	escapeSequence ANSIEscapeSequence
}

func NewColor(escapeSequence ANSIEscapeSequence) (Color, error) {
	if !escapeSequence.IsColor() {
		return Color{}, ErrInvalidEscapeSequence
	}
	return Color{
		escapeSequence: escapeSequence,
	}, nil
}

func (c Color) String() string {
	return string(c.escapeSequence)
}
