package domain

import (
	"errors"
	"tuigo/presentation/ansi"
)

var (
	ErrInvalidEscapeSequence = errors.New("Invalid Escape sequence")
)

type Color struct {
	escapeSequence ansi.ANSIEscapeSequence
}

func (c *Color) New(escapeSequence ansi.ANSIEscapeSequence) (Color, error) {
	if !escapeSequence.IsColor() {
		return Color{}, ErrInvalidEscapeSequence
	}
	return Color{
		escapeSequence: escapeSequence,
	}, nil
}
