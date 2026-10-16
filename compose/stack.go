package compose

import (
	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
)

type Stack struct {
	Width, Height int
	Children      []Node
}

func (s *Stack) Frame() (canvas.Frame, error) {
	if err := validateDimensions(s.Width, s.Height); err != nil {
		return canvas.Frame{}, err
	}
	cells := make([]canvas.Cell, s.Width*s.Height)
	for _, child := range s.Children {
		frame, err := child.Frame()
		if err != nil {
			return canvas.Frame{}, err
		}
		for y := range frame.Height() {
			if y >= s.Height {
				break
			}
			frameRow, err := frame.RowAt(y)
			if err != nil {
				return canvas.Frame{}, err
			}
			cellsRow := cells[y*s.Width : (y+1)*s.Width]
			if err := overlayRow(cellsRow, frameRow); err != nil {
				return canvas.Frame{}, err
			}
		}
	}
	return canvas.NewFrame(s.Width, s.Height, cells)
}

func overlayRow(dst, src []canvas.Cell) error {
	for start := 0; start < len(src); start += src[start].Width() {
		if start >= len(dst) {
			break
		}
		cell := src[start]
		if cell.Foreground() == color.FgTransparent ||
			cell.Background() == color.BgTransparent {
			continue
		}
		end := start + cell.Width()
		if end > len(dst) {
			break
		}
		if start > 0 && dst[start].Width() == 0 {
			if err := clearGlyphCrossingBoundary(dst, start); err != nil {
				return err
			}
		}
		if end < len(dst) && dst[end].Width() == 0 {
			if err := clearGlyphCrossingBoundary(dst, end); err != nil {
				return err
			}
		}
		copy(dst[start:end], src[start:end])
	}
	return nil
}

func clearGlyphCrossingBoundary(dst []canvas.Cell, boundary int) error {
	head := boundary
	for head > 0 && dst[head].Width() == 0 {
		head--
	}
	tail := head + dst[head].Width()
	if tail > boundary {
		blank, err := canvas.NewCell(" ", dst[head].Foreground(), dst[head].Background())
		if err != nil {
			return err
		}
		for ; head < tail; head++ {
			dst[head] = blank
		}
	}
	return nil
}
