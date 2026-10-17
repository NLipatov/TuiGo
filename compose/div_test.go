package compose

import (
	"errors"
	"math"
	"slices"
	"testing"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
)

func TestDivLaysOutChildren(t *testing.T) {
	parentBlank := divTestCell(t, " ", color.FgWhite, color.BgBlue)
	childBlank := divTestCell(t, " ", color.FgGreen, color.BgPurple)
	tests := []struct {
		name string
		div  *Div
		want [][]canvas.Cell
	}{
		{
			name: "row places a taller child beside text",
			div: &Div{
				Width: 4, Height: 2, Fg: color.FgWhite, Bg: color.BgBlue,
				Direction: Row,
				Children: []Node{
					&Text{Content: "A", Fg: color.FgRed, Bg: color.BgYellow},
					&Div{Width: 2, Height: 2, Fg: color.FgGreen, Bg: color.BgPurple,
						Children: []Node{&Text{Content: "B", Fg: color.FgGreen, Bg: color.BgPurple}}},
				},
			},
			want: [][]canvas.Cell{
				{divTestCell(t, "A", color.FgRed, color.BgYellow), divTestCell(t, "B", color.FgGreen, color.BgPurple), childBlank, parentBlank},
				{parentBlank, childBlank, childBlank, parentBlank},
			},
		},
		{
			name: "column places children on successive rows",
			div: &Div{
				Width: 3, Height: 3, Fg: color.FgWhite, Bg: color.BgBlue,
				Direction: Column,
				Children: []Node{
					&Text{Content: "AB", Fg: color.FgRed, Bg: color.BgYellow},
					&Text{Content: "C", Fg: color.FgGreen, Bg: color.BgPurple},
				},
			},
			want: [][]canvas.Cell{
				{divTestCell(t, "A", color.FgRed, color.BgYellow), divTestCell(t, "B", color.FgRed, color.BgYellow), parentBlank},
				{divTestCell(t, "C", color.FgGreen, color.BgPurple), parentBlank, parentBlank},
				{parentBlank, parentBlank, parentBlank},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertDivFrame(t, tt.div, tt.want)
		})
	}
}

func TestDivClipsWideChildAtRightEdge(t *testing.T) {
	blank := divTestCell(t, " ", color.FgWhite, color.BgBlue)
	letter := divTestCell(t, "A", color.FgRed, color.BgYellow)
	wide := divTestCell(t, "🙂", color.FgRed, color.BgYellow)
	for _, tt := range []struct {
		name  string
		width int
		want  []canvas.Cell
	}{
		{name: "wide glyph fits", width: 3, want: []canvas.Cell{letter, wide, {}}},
		{name: "wide glyph is clipped", width: 2, want: []canvas.Cell{letter, blank}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			div := &Div{
				Width: tt.width, Height: 1, Fg: color.FgWhite, Bg: color.BgBlue,
				Direction: Row,
				Children: []Node{
					&Text{Content: "A", Fg: color.FgRed, Bg: color.BgYellow},
					&Text{Content: "🙂", Fg: color.FgRed, Bg: color.BgYellow},
				},
			}
			assertDivFrame(t, div, [][]canvas.Cell{tt.want})
		})
	}
}

func TestDivRejectsInvalidDimensions(t *testing.T) {
	for _, tt := range []struct {
		name          string
		width, height int
		want          error
	}{
		{name: "zero width", width: 0, height: 1, want: ErrInvalidDimensions},
		{name: "area overflow", width: math.MaxInt, height: 2, want: ErrAreaOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := (&Div{Width: tt.width, Height: tt.height}).Frame()
			if !errors.Is(err, tt.want) {
				t.Fatalf("Frame() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func assertDivFrame(t *testing.T, div *Div, want [][]canvas.Cell) {
	t.Helper()
	frame, err := div.Frame()
	if err != nil {
		t.Fatalf("Frame() error = %v", err)
	}
	if frame.Width() != len(want[0]) || frame.Height() != len(want) {
		t.Fatalf("Frame() size = %dx%d, want %dx%d", frame.Width(), frame.Height(), len(want[0]), len(want))
	}
	for y, expected := range want {
		row, err := frame.RowAt(y)
		if err != nil {
			t.Fatalf("RowAt(%d) error = %v", y, err)
		}
		if !slices.Equal(row, expected) {
			t.Errorf("row %d = %#v, want %#v", y, row, expected)
		}
	}
}

func divTestCell(t *testing.T, glyph string, fg, bg color.Color) canvas.Cell {
	t.Helper()
	cell, err := canvas.NewCell(glyph, fg, bg)
	if err != nil {
		t.Fatalf("NewCell(%q) error = %v", glyph, err)
	}
	return cell
}
