package compose

import (
	"errors"
	"math"
)

var (
	ErrInvalidDimensions = errors.New("invalid dimensions")
	ErrAreaOverflow      = errors.New("given width and height multiplication results in overflow")
)

func validateDimensions(w, h int) error {
	if w <= 0 || h <= 0 {
		return ErrInvalidDimensions
	}
	if w > math.MaxInt/h {
		return ErrAreaOverflow
	}
	return nil
}
