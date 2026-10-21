package compose

import (
	"errors"
	"math"
	"slices"
	"testing"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
)

func TestStackOverlaysSmallerText(t *testing.T) {
	base := &Div{Width: 4, Height: 2, Fg: color.FgWhite, Bg: color.BgBlue}
	text := &Text{Content: "Hi", Fg: color.FgRed, Bg: color.BgYellow}
	stack := &Stack{Width: 4, Height: 2, Children: []Node{base, text}}
	blank := stackTestCell(t, " ", color.FgWhite, color.BgBlue)

	assertStackFrame(t, stack, [][]canvas.Cell{
		{stackTestCell(t, "H", color.FgRed, color.BgYellow), stackTestCell(t, "i", color.FgRed, color.BgYellow), blank, blank},
		{blank, blank, blank, blank},
	})
}

func TestStackOverlaysNestedTransparentStack(t *testing.T) {
	layer := &Div{
		Width: 2, Height: 1, Fg: color.FgTransparent, Bg: color.BgTransparent,
		Children: []Node{&Text{Content: "X", Fg: color.FgRed, Bg: color.BgYellow}},
	}
	inner := &Stack{Width: 2, Height: 1, Children: []Node{layer}}
	outer := &Stack{
		Width: 2, Height: 1,
		Children: []Node{
			&Div{Width: 2, Height: 1, Fg: color.FgWhite, Bg: color.BgBlue},
			inner,
		},
	}

	assertStackFrame(t, outer, [][]canvas.Cell{{
		stackTestCell(t, "X", color.FgRed, color.BgYellow),
		stackTestCell(t, " ", color.FgWhite, color.BgBlue),
	}})
}

func TestStackClipsOversizedChild(t *testing.T) {
	base := &Div{Width: 2, Height: 1, Fg: color.FgWhite, Bg: color.BgBlue}
	child := &Div{
		Width: 3, Height: 2, Fg: color.FgRed, Bg: color.BgYellow,
		Children: []Node{&Text{Content: "ABC", Fg: color.FgRed, Bg: color.BgYellow}},
	}
	stack := &Stack{Width: 2, Height: 1, Children: []Node{base, child}}

	assertStackFrame(t, stack, [][]canvas.Cell{{
		stackTestCell(t, "A", color.FgRed, color.BgYellow),
		stackTestCell(t, "B", color.FgRed, color.BgYellow),
	}})
}

func TestStackClipsWideGraphemeAtRightEdge(t *testing.T) {
	base := &Div{Width: 2, Height: 1, Fg: color.FgWhite, Bg: color.BgBlue}
	text := &Text{Content: "A🙂", Fg: color.FgRed, Bg: color.BgYellow}
	stack := &Stack{Width: 2, Height: 1, Children: []Node{base, text}}

	assertStackFrame(t, stack, [][]canvas.Cell{{
		stackTestCell(t, "A", color.FgRed, color.BgYellow),
		stackTestCell(t, " ", color.FgWhite, color.BgBlue),
	}})
}

func TestStackClearsPartiallyCoveredWideGraphemes(t *testing.T) {
	transparent := func(width int, children ...Node) *Div {
		return &Div{
			Width: width, Height: 1, Fg: color.FgTransparent, Bg: color.BgTransparent,
			Children: children,
		}
	}
	letter := &Text{Content: "A", Fg: color.FgRed, Bg: color.BgYellow}
	wide := &Text{Content: "🔥", Fg: color.FgRed, Bg: color.BgYellow}
	lowerBlank := stackTestCell(t, " ", color.FgWhite, color.BgBlue)
	upperLetter := stackTestCell(t, "A", color.FgRed, color.BgYellow)

	tests := []struct {
		name    string
		base    string
		overlay Node
		want    []canvas.Cell
	}{
		{
			name:    "replace wide grapheme head",
			base:    "🙂",
			overlay: transparent(2, letter),
			want:    []canvas.Cell{upperLetter, lowerBlank},
		},
		{
			name:    "replace wide grapheme continuation",
			base:    "🙂",
			overlay: transparent(2, transparent(1), letter),
			want:    []canvas.Cell{lowerBlank, upperLetter},
		},
		{
			name:    "overlap two wide graphemes",
			base:    "🙂🙂",
			overlay: transparent(4, transparent(1), wide),
			want: []canvas.Cell{
				lowerBlank, stackTestCell(t, "🔥", color.FgRed, color.BgYellow), {}, lowerBlank,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := &Text{Content: tt.base, Fg: color.FgWhite, Bg: color.BgBlue}
			stack := &Stack{Width: len(tt.want), Height: 1, Children: []Node{base, tt.overlay}}
			assertStackFrame(t, stack, [][]canvas.Cell{tt.want})
		})
	}
}

func TestStackClearsCrossingWidthThreeGlyphs(t *testing.T) {
	// Synthetic widths exercise frame layout independently of terminal width calculation.
	lowerHead, err := canvas.NewCellWithWidth("L", 3, color.FgWhite, color.BgBlue)
	if err != nil {
		t.Fatal(err)
	}
	upperHead, err := canvas.NewCellWithWidth("U", 3, color.FgRed, color.BgYellow)
	if err != nil {
		t.Fatal(err)
	}
	transparent := stackTestCell(t, " ", color.FgTransparent, color.BgTransparent)
	lower, err := canvas.NewFrame(6, 1, []canvas.Cell{lowerHead, {}, {}, lowerHead, {}, {}})
	if err != nil {
		t.Fatal(err)
	}
	upper, err := canvas.NewFrame(6, 1, []canvas.Cell{transparent, upperHead, {}, {}, transparent, transparent})
	if err != nil {
		t.Fatal(err)
	}

	stack := &Stack{Width: 6, Height: 1, Children: []Node{stackFrameNode{lower}, stackFrameNode{upper}}}
	lowerBlank := stackTestCell(t, " ", color.FgWhite, color.BgBlue)
	assertStackFrame(t, stack, [][]canvas.Cell{{lowerBlank, upperHead, {}, {}, lowerBlank, lowerBlank}})
}

func TestStackSkipsCellWithEitherTransparentColor(t *testing.T) {
	tests := []struct {
		name    string
		base    string
		overlay *Text
		want    []canvas.Cell
	}{
		{
			name:    "transparent foreground",
			base:    "B",
			overlay: &Text{Content: "A", Fg: color.FgTransparent, Bg: color.BgYellow},
			want:    []canvas.Cell{stackTestCell(t, "B", color.FgWhite, color.BgBlue)},
		},
		{
			name:    "transparent background",
			base:    "B",
			overlay: &Text{Content: "A", Fg: color.FgRed, Bg: color.BgTransparent},
			want:    []canvas.Cell{stackTestCell(t, "B", color.FgWhite, color.BgBlue)},
		},
		{
			name:    "transparent wide grapheme",
			base:    "BC",
			overlay: &Text{Content: "🙂", Fg: color.FgTransparent, Bg: color.BgYellow},
			want: []canvas.Cell{
				stackTestCell(t, "B", color.FgWhite, color.BgBlue),
				stackTestCell(t, "C", color.FgWhite, color.BgBlue),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := &Text{Content: tt.base, Fg: color.FgWhite, Bg: color.BgBlue}
			stack := &Stack{Width: len(tt.want), Height: 1, Children: []Node{base, tt.overlay}}
			assertStackFrame(t, stack, [][]canvas.Cell{tt.want})
		})
	}
}

func TestStackResetsUncoveredCellsBetweenFrames(t *testing.T) {
	text := &Text{Content: "🙂", Fg: color.FgWhite, Bg: color.BgBlue}
	stack := &Stack{Width: 2, Height: 1, Children: []Node{text}}
	assertStackFrame(t, stack, [][]canvas.Cell{{stackTestCell(t, "🙂", color.FgWhite, color.BgBlue), {}}})

	text.Fg = color.FgTransparent
	blank := stackTestCell(t, " ", color.FgTransparent, color.BgTransparent)
	assertStackFrame(t, stack, [][]canvas.Cell{{blank, blank}})
}

func TestStackRejectsInvalidDimensions(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
		want          error
	}{
		{name: "zero width", width: 0, height: 1, want: ErrInvalidDimensions},
		{name: "negative height", width: 1, height: -1, want: ErrInvalidDimensions},
		{name: "area overflow", width: math.MaxInt, height: 2, want: ErrAreaOverflow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := (&Stack{Width: tt.width, Height: tt.height}).Frame()
			if !errors.Is(err, tt.want) {
				t.Fatalf("Frame() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func assertStackFrame(t *testing.T, stack *Stack, want [][]canvas.Cell) {
	t.Helper()
	frame, err := stack.Frame()
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

func stackTestCell(t *testing.T, glyph string, fg, bg color.Color) canvas.Cell {
	t.Helper()
	cell, err := canvas.NewCell(glyph, fg, bg)
	if err != nil {
		t.Fatalf("NewCell(%q) error = %v", glyph, err)
	}
	return cell
}

type stackFrameNode struct {
	frame canvas.Frame
}

func (n stackFrameNode) Frame() (canvas.Frame, error) {
	return n.frame, nil
}
