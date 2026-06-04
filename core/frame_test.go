package core

import (
	"errors"
	"testing"

	"github.com/NLipatov/tuigo/ansi"
)

func TestNewFrameExposesDimensionsAndCells(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "single cell",
			width:  1,
			height: 1,
		},
		{
			name:   "rectangular frame",
			width:  4,
			height: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cells := make([]Cell, tt.width*tt.height)
			for i := range cells {
				cells[i] = testCellWithSymbol(t, rune('a'+i))
			}
			frame, err := NewFrame(tt.width, tt.height, cells)
			if err != nil {
				t.Fatalf("NewFrame() error = %v", err)
			}

			if frame.Width() != tt.width {
				t.Fatalf("Width() = %d, want %d", frame.Width(), tt.width)
			}
			if frame.Height() != tt.height {
				t.Fatalf("Height() = %d, want %d", frame.Height(), tt.height)
			}
			got, err := frame.CellAt(tt.width-1, tt.height-1)
			if err != nil {
				t.Fatalf("CellAt() error = %v", err)
			}
			if want := cells[len(cells)-1]; got != want {
				t.Fatalf("CellAt() = %#v, want %#v", got, want)
			}
		})
	}
}

func TestNewFrameRejectsInvalidDimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		cells  []Cell
	}{
		{
			name:   "negative width",
			width:  -1,
			height: 1,
			cells:  nil,
		},
		{
			name:   "zero width",
			width:  0,
			height: 1,
			cells:  nil,
		},
		{
			name:   "negative height",
			width:  1,
			height: -1,
			cells:  nil,
		},
		{
			name:   "zero height",
			width:  1,
			height: 0,
			cells:  nil,
		},
		{
			name:   "cell count does not match dimensions",
			width:  2,
			height: 2,
			cells:  make([]Cell, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFrame(tt.width, tt.height, tt.cells)
			if !errors.Is(err, ErrInvalidFrameDimensions) {
				t.Fatalf("NewFrame() error = %v, want %v", err, ErrInvalidFrameDimensions)
			}
		})
	}
}

func TestFrameCellAtReadsCellsByRowMajorIndex(t *testing.T) {
	cells := []Cell{
		testCell(t, ansi.FG_BLACK),
		testCell(t, ansi.FG_RED),
		testCell(t, ansi.FG_GREEN),
		testCell(t, ansi.FG_YELLOW),
		testCell(t, ansi.FG_BLUE),
		testCell(t, ansi.FG_PURPLE),
		testCell(t, ansi.FG_CYAN),
		testCell(t, ansi.FG_WHITE),
		testCell(t, ansi.FG_BOLD_BLACK),
		testCell(t, ansi.FG_BOLD_RED),
		testCell(t, ansi.FG_BOLD_GREEN),
		testCell(t, ansi.FG_BOLD_YELLOW),
	}

	tests := []struct {
		name    string
		x       int
		y       int
		wantIdx int
	}{
		{
			name:    "top left",
			x:       0,
			y:       0,
			wantIdx: 0,
		},
		{
			name:    "same row",
			x:       2,
			y:       0,
			wantIdx: 2,
		},
		{
			name:    "next row",
			x:       0,
			y:       1,
			wantIdx: 4,
		},
		{
			name:    "bottom right",
			x:       3,
			y:       2,
			wantIdx: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := Frame{
				width:  4,
				height: 3,
				cells:  cells,
			}

			got, err := frame.CellAt(tt.x, tt.y)
			if err != nil {
				t.Fatalf("CellAt() error = %v", err)
			}
			if want := cells[tt.wantIdx]; got != want {
				t.Fatalf("CellAt() = %#v, want %#v", got, want)
			}
		})
	}
}

func TestFrameCellAtRejectsOutOfBoundsCoordinates(t *testing.T) {
	tests := []struct {
		name string
		x    int
		y    int
	}{
		{
			name: "x below zero",
			x:    -1,
			y:    0,
		},
		{
			name: "y below zero",
			x:    0,
			y:    -1,
		},
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
			frame, err := NewFrame(3, 2, make([]Cell, 6))
			if err != nil {
				t.Fatalf("NewFrame() error = %v", err)
			}

			_, err = frame.CellAt(tt.x, tt.y)
			if !errors.Is(err, ErrOutOfFrameBounds) {
				t.Fatalf("CellAt() error = %v, want %v", err, ErrOutOfFrameBounds)
			}
		})
	}
}

func TestFrameRowAtReadsCellsByRowMajorIndex(t *testing.T) {
	cells := []Cell{
		testCell(t, ansi.FG_BLACK),
		testCell(t, ansi.FG_RED),
		testCell(t, ansi.FG_GREEN),
		testCell(t, ansi.FG_YELLOW),
		testCell(t, ansi.FG_BLUE),
		testCell(t, ansi.FG_PURPLE),
	}

	frame := Frame{
		width:  3,
		height: 2,
		cells:  cells,
	}

	got, err := frame.RowAt(1)
	if err != nil {
		t.Fatalf("RowAt() error = %v", err)
	}

	want := cells[3:6]
	if len(got) != len(want) {
		t.Fatalf("len(RowAt()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("RowAt()[%d] = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestFrameRowAtLimitsCapacityToRow(t *testing.T) {
	frame, err := NewFrame(3, 2, make([]Cell, 6))
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
	}

	for y := range frame.Height() {
		row, err := frame.RowAt(y)
		if err != nil {
			t.Fatalf("RowAt(%d) error = %v", y, err)
		}
		if cap(row) != frame.Width() {
			t.Fatalf("cap(RowAt(%d)) = %d, want %d", y, cap(row), frame.Width())
		}
	}
}

func TestFrameRowAtRejectsOutOfBoundsRows(t *testing.T) {
	tests := []struct {
		name string
		y    int
	}{
		{
			name: "y below zero",
			y:    -1,
		},
		{
			name: "y equals height",
			y:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame, err := NewFrame(3, 2, make([]Cell, 6))
			if err != nil {
				t.Fatalf("NewFrame() error = %v", err)
			}

			_, err = frame.RowAt(tt.y)
			if !errors.Is(err, ErrOutOfFrameBounds) {
				t.Fatalf("RowAt() error = %v, want %v", err, ErrOutOfFrameBounds)
			}
		})
	}
}

func TestFrameCellAtDoesNotAllocate(t *testing.T) {
	frame := Frame{
		width:  2,
		height: 2,
		cells: []Cell{
			testCell(t, ansi.FG_BLACK),
			testCell(t, ansi.FG_RED),
			testCell(t, ansi.FG_GREEN),
			testCell(t, ansi.FG_BLUE),
		},
	}

	if _, err := frame.CellAt(1, 1); err != nil {
		t.Fatalf("CellAt() error = %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_, _ = frame.CellAt(1, 1)
	})
	if allocs != 0 {
		t.Fatalf("allocations per CellAt = %.2f, want 0", allocs)
	}
}

func TestFrameRowAtDoesNotAllocate(t *testing.T) {
	frame := Frame{
		width:  2,
		height: 2,
		cells: []Cell{
			testCell(t, ansi.FG_BLACK),
			testCell(t, ansi.FG_RED),
			testCell(t, ansi.FG_GREEN),
			testCell(t, ansi.FG_BLUE),
		},
	}

	if _, err := frame.RowAt(1); err != nil {
		t.Fatalf("RowAt() error = %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_, _ = frame.RowAt(1)
	})
	if allocs != 0 {
		t.Fatalf("allocations per RowAt = %.2f, want 0", allocs)
	}
}

func testCell(t *testing.T, sequence ansi.ANSIEscapeSequence) Cell {
	t.Helper()

	fg, err := ansi.NewColor(sequence)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", sequence, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	return NewCell('x', fg, bg)
}

func testCellWithSymbol(t *testing.T, symbol rune) Cell {
	t.Helper()

	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	return NewCell(symbol, fg, bg)
}
