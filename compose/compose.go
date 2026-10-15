package compose

import (
	"errors"
	"math"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
	"github.com/rivo/uniseg"
)

var (
	ErrInvalidDimensions = errors.New("invalid dimensions")
	ErrEmptyText         = errors.New("text element should have non-empty content")
	ErrAreaOverflow      = errors.New("given width and height multiplication results in overflow")
)

type Node interface {
	Frame() (canvas.Frame, error)
}

type Direction int

const (
	Row Direction = iota
	Column
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

type Div struct {
	Width, Height int
	Fg, Bg        color.Color
	Children      []Node
	Direction     Direction
}

func (d *Div) Frame() (canvas.Frame, error) {
	if err := validateDimensions(d.Width, d.Height); err != nil {
		return canvas.Frame{}, err
	}
	cells := make([]canvas.Cell, d.Width*d.Height)
	for i := range cells {
		cell, err := canvas.NewCell(" ", d.Fg, d.Bg)
		if err != nil {
			return canvas.Frame{}, err
		}
		cells[i] = cell
	}
	switch d.Direction {
	case Row:
		err := d.renderRow(cells)
		if err != nil {
			return canvas.Frame{}, err
		}
	case Column:
		err := d.renderColumn(cells)
		if err != nil {
			return canvas.Frame{}, err
		}
	}
	return canvas.NewFrame(d.Width, d.Height, cells)
}

func (d *Div) renderRow(cells []canvas.Cell) error {
	offsetX := 0
	for _, child := range d.Children {
		frame, err := child.Frame()
		if err != nil {
			return err
		}
		for y := range frame.Height() {
			if y >= d.Height {
				break
			}
			row, err := frame.RowAt(y)
			if err != nil {
				return err
			}
			for x, cell := range row {
				if x+offsetX >= d.Width {
					break
				}
				if x+offsetX+cell.Width() > d.Width {
					break
				}
				idx := y*d.Width + (x + offsetX)
				if idx >= len(cells) {
					break
				}
				cells[idx] = cell
			}
		}
		offsetX += frame.Width()
	}
	return nil
}

func (d *Div) renderColumn(cells []canvas.Cell) error {
	offsetY := 0
	for _, child := range d.Children {
		frame, err := child.Frame()
		if err != nil {
			return err
		}
		for y := range frame.Height() {
			if y+offsetY >= d.Height {
				break
			}
			row, err := frame.RowAt(y)
			if err != nil {
				return err
			}
			for x, cell := range row {
				if x >= d.Width {
					break
				}
				if x+cell.Width() > d.Width {
					break
				}
				idx := (y+offsetY)*d.Width + x
				if idx >= len(cells) {
					break
				}
				cells[idx] = cell
			}
		}
		offsetY += frame.Height()
	}
	return nil
}

type Text struct {
	Content string
	Fg, Bg  color.Color
}

func (t *Text) Frame() (canvas.Frame, error) {
	if len(t.Content) == 0 {
		return canvas.Frame{}, ErrEmptyText
	}
	gs := uniseg.NewGraphemes(t.Content)
	cells := make([]canvas.Cell, 0, len(t.Content))
	for gs.Next() {
		cell, err := canvas.NewCellWithWidth(gs.Str(), gs.Width(), t.Fg, t.Bg)
		if err != nil {
			return canvas.Frame{}, err
		}
		cells = append(cells, cell)
		for range gs.Width() - 1 {
			cells = append(cells, canvas.Cell{})
		}
	}
	return canvas.NewFrame(len(cells), 1, cells)
}

func validateDimensions(w, h int) error {
	if w <= 0 || h <= 0 {
		return ErrInvalidDimensions
	}
	if w > math.MaxInt/h {
		return ErrAreaOverflow
	}
	return nil
}
