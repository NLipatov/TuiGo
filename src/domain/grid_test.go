package domain

import (
	"errors"
	"testing"
	"tuigo/presentation/ansi"
)

func TestNewGrid(t *testing.T) {
	tests := []struct {
		name   string
		width  uint
		height uint
	}{
		{
			name:   "single cell",
			width:  1,
			height: 1,
		},
		{
			name:   "rectangular grid",
			width:  4,
			height: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGrid(tt.width, tt.height)

			if grid.width != tt.width {
				t.Fatalf("width = %d, want %d", grid.width, tt.width)
			}
			if grid.height != tt.height {
				t.Fatalf("height = %d, want %d", grid.height, tt.height)
			}
			if len(grid.cells) != int(tt.width*tt.height) {
				t.Fatalf("len(cells) = %d, want %d", len(grid.cells), tt.width*tt.height)
			}
		})
	}
}

func TestGridSetCellAtStoresCellsByRowMajorIndex(t *testing.T) {
	tests := []struct {
		name    string
		x       uint
		y       uint
		wantIdx uint
		cell    Cell
	}{
		{
			name:    "top left",
			x:       0,
			y:       0,
			wantIdx: 0,
			cell:    testCell(ansi.FG_RED),
		},
		{
			name:    "same row",
			x:       2,
			y:       0,
			wantIdx: 2,
			cell:    testCell(ansi.FG_GREEN),
		},
		{
			name:    "next row",
			x:       0,
			y:       1,
			wantIdx: 4,
			cell:    testCell(ansi.FG_BLUE),
		},
		{
			name:    "bottom right",
			x:       3,
			y:       2,
			wantIdx: 11,
			cell:    testCell(ansi.FG_WHITE),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGrid(4, 3)

			if err := grid.SetCellAt(tt.x, tt.y, tt.cell); err != nil {
				t.Fatalf("SetCellAt() error = %v", err)
			}

			if got := grid.cells[tt.wantIdx]; got != tt.cell {
				t.Fatalf("cells[%d] = %#v, want %#v", tt.wantIdx, got, tt.cell)
			}

			got, err := grid.CellAt(tt.x, tt.y)
			if err != nil {
				t.Fatalf("CellAt() error = %v", err)
			}
			if got != tt.cell {
				t.Fatalf("CellAt() = %#v, want %#v", got, tt.cell)
			}
		})
	}
}

func TestGridCellAtRejectsOutOfBoundsCoordinates(t *testing.T) {
	tests := []struct {
		name string
		x    uint
		y    uint
	}{
		{
			name: "x equals width",
			x:    3,
			y:    0,
		},
		{
			name: "y equals height",
			x:    0,
			y:    2,
		},
		{
			name: "x and y out of bounds",
			x:    3,
			y:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGrid(3, 2)

			_, err := grid.CellAt(tt.x, tt.y)
			if !errors.Is(err, ErrOutOfGridBounds) {
				t.Fatalf("CellAt() error = %v, want %v", err, ErrOutOfGridBounds)
			}
		})
	}
}

func TestGridSetCellAtRejectsOutOfBoundsCoordinates(t *testing.T) {
	tests := []struct {
		name string
		x    uint
		y    uint
	}{
		{
			name: "x equals width",
			x:    3,
			y:    0,
		},
		{
			name: "y equals height",
			x:    0,
			y:    2,
		},
		{
			name: "x and y out of bounds",
			x:    3,
			y:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGrid(3, 2)

			err := grid.SetCellAt(tt.x, tt.y, testCell(ansi.FG_RED))
			if !errors.Is(err, ErrOutOfGridBounds) {
				t.Fatalf("SetCellAt() error = %v, want %v", err, ErrOutOfGridBounds)
			}
		})
	}
}

func TestGridAccessDoesNotAllocate(t *testing.T) {
	grid := NewGrid(2, 2)
	cell := testCell(ansi.FG_RED)

	if err := grid.SetCellAt(1, 1, cell); err != nil {
		t.Fatalf("SetCellAt() error = %v", err)
	}
	if _, err := grid.CellAt(1, 1); err != nil {
		t.Fatalf("CellAt() error = %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = grid.SetCellAt(1, 1, cell)
		_, _ = grid.CellAt(1, 1)
	})
	if allocs != 0 {
		t.Fatalf("allocations per access = %.2f, want 0", allocs)
	}
}

func testCell(sequence ansi.ANSIEscapeSequence) Cell {
	return Cell{
		fg: Color{
			escapeSequence: sequence,
		},
	}
}
