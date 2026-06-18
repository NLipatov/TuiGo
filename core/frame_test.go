package core

import (
	"errors"
	"testing"

	"github.com/NLipatov/tuigo/color"
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
				cells[i] = testCellWithGlyph(t, rune('a'+i))
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

func TestNewFrameAcceptsWideCellWithContinuation(t *testing.T) {
	cells := []Cell{
		testCellWithGlyph(t, 'A'),
		testCellWithGlyph(t, '🙂'),
		{},
		testCellWithGlyph(t, 'B'),
	}

	frame, err := NewFrame(4, 1, cells)
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
	}

	got, err := frame.CellAt(2, 0)
	if err != nil {
		t.Fatalf("CellAt() error = %v", err)
	}
	if got != (Cell{}) {
		t.Fatalf("CellAt() = %#v, want continuation block", got)
	}
}

func TestNewFrameRejectsInvalidCellLayout(t *testing.T) {
	blank := testCellWithGlyph(t, ' ')
	wide := testCellWithGlyph(t, '🙂')

	tests := []struct {
		name       string
		width      int
		cells      []Cell
		wantX      int
		wantReason string
	}{
		{
			name:       "unexpected continuation at row start",
			width:      3,
			cells:      []Cell{{}, blank, blank},
			wantX:      0,
			wantReason: "unexpected continuation block",
		},
		{
			name:       "unexpected continuation after normal cell",
			width:      3,
			cells:      []Cell{blank, {}, blank},
			wantX:      1,
			wantReason: "unexpected continuation block",
		},
		{
			name:       "wide cell at row end",
			width:      3,
			cells:      []Cell{blank, blank, wide},
			wantX:      2,
			wantReason: "missing continuation block",
		},
		{
			name:       "wide cell followed by normal cell",
			width:      3,
			cells:      []Cell{blank, wide, blank},
			wantX:      1,
			wantReason: "missing continuation block",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewFrame(tt.width, 1, tt.cells)
			var layoutErr FrameCellLayoutError
			if !errors.As(err, &layoutErr) {
				t.Fatalf("NewFrame() error = %v, want FrameCellLayoutError", err)
			}
			if layoutErr.X != tt.wantX || layoutErr.Y != 0 || layoutErr.Reason != tt.wantReason {
				t.Fatalf(
					"FrameCellLayoutError = %#v, want x=%d y=0 reason=%q",
					layoutErr,
					tt.wantX,
					tt.wantReason,
				)
			}
		})
	}
}

func TestFrameCellAtReadsCellsByRowMajorIndex(t *testing.T) {
	cells := []Cell{
		testCell(t, color.FgBlack),
		testCell(t, color.FgRed),
		testCell(t, color.FgGreen),
		testCell(t, color.FgYellow),
		testCell(t, color.FgBlue),
		testCell(t, color.FgPurple),
		testCell(t, color.FgCyan),
		testCell(t, color.FgWhite),
		testCell(t, color.FgBoldBlack),
		testCell(t, color.FgBoldRed),
		testCell(t, color.FgBoldGreen),
		testCell(t, color.FgBoldYellow),
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
			frame, err := newBlankFrame(t, 3, 2)
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
		testCell(t, color.FgBlack),
		testCell(t, color.FgRed),
		testCell(t, color.FgGreen),
		testCell(t, color.FgYellow),
		testCell(t, color.FgBlue),
		testCell(t, color.FgPurple),
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
	frame, err := newBlankFrame(t, 3, 2)
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
			frame, err := newBlankFrame(t, 3, 2)
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
			testCell(t, color.FgBlack),
			testCell(t, color.FgRed),
			testCell(t, color.FgGreen),
			testCell(t, color.FgBlue),
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
			testCell(t, color.FgBlack),
			testCell(t, color.FgRed),
			testCell(t, color.FgGreen),
			testCell(t, color.FgBlue),
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

func newBlankFrame(t *testing.T, width, height int) (Frame, error) {
	t.Helper()

	blank := testCellWithGlyph(t, ' ')
	cells := make([]Cell, width*height)
	for i := range cells {
		cells[i] = blank
	}
	return NewFrame(width, height, cells)
}

func testCell(t *testing.T, fg color.Color) Cell {
	t.Helper()

	cell, err := NewCell("x", fg, color.BgBlack)
	if err != nil {
		t.Fatalf("NewCell(%q) error = %v", "x", err)
	}
	return cell
}

func testCellWithGlyph(t *testing.T, symbol rune) Cell {
	t.Helper()

	cell, err := NewCell(string(symbol), color.FgRed, color.BgBlack)
	if err != nil {
		t.Fatalf("NewCell(%q) error = %v", string(symbol), err)
	}
	return cell
}
