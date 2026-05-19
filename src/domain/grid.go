package domain

import (
	"errors"
)

var (
	ErrOutOfGridBounds = errors.New("Out of grid bounds")
)

type Grid struct {
	height, width uint
	cells         []Cell
}

func NewGrid(width, height uint) Grid {
	return Grid{
		height: height,
		width:  width,
		cells:  make([]Cell, width*height),
	}
}

func (g *Grid) CellAt(x, y uint) (Cell, error) {
	idx, err := g.idx(x, y)
	if err != nil {
		return Cell{}, err
	}
	return g.cells[idx], nil
}

func (g *Grid) SetCellAt(x, y uint, cell Cell) error {
	idx, err := g.idx(x, y)
	if err != nil {
		return err
	}
	g.cells[idx] = cell
	return nil
}

func (g *Grid) idx(x, y uint) (uint, error) {
	if x >= g.width || y >= g.height {
		return 0, ErrOutOfGridBounds
	}
	idx := g.width*y + x
	if idx >= uint(len(g.cells)) {
		return 0, ErrOutOfGridBounds
	}
	return idx, nil
}
