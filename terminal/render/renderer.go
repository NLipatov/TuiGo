package render

import (
	"io"
	"strconv"

	"github.com/NLipatov/tuigo/canvas"
	"github.com/NLipatov/tuigo/color"
	"github.com/NLipatov/tuigo/internal/ansi"
)

type renderStyle struct {
	fg, bg color.Color
	set    bool
}

type Renderer struct {
	frame       canvas.Frame
	fullRepaint bool
	writer      io.Writer
	out         []byte
	style       renderStyle
}

func NewRenderer(writer io.Writer) *Renderer {
	return &Renderer{
		fullRepaint: true,
		writer:      writer,
	}
}

// Render writes the changes in frame. The frame's backing cell buffer is not
// retained and may be mutated after Render returns.
func (r *Renderer) Render(frame canvas.Frame) error {
	if err := r.ensurePreviousFrame(frame); err != nil {
		return err
	}
	r.style.set = false
	r.out = r.out[:0]
	if r.fullRepaint {
		if err := r.renderFullFrame(frame); err != nil {
			return err
		}
	} else {
		if err := r.renderDiffFrame(frame); err != nil {
			return err
		}
	}
	if err := r.flush(); err != nil {
		r.fullRepaint = true
		return err
	}
	r.fullRepaint = false
	return nil
}

func (r *Renderer) ensurePreviousFrame(frame canvas.Frame) error {
	if frame.Width() == r.frame.Width() && frame.Height() == r.frame.Height() {
		return nil
	}

	cells := make([]canvas.Cell, frame.Width()*frame.Height())
	for y := range frame.Height() {
		row, err := frame.RowAt(y)
		if err != nil {
			return err
		}
		start := y * frame.Width()
		copy(cells[start:start+frame.Width()], row)
	}
	previous, err := canvas.NewFrame(frame.Width(), frame.Height(), cells)
	if err != nil {
		return err
	}
	r.frame = previous
	r.fullRepaint = true
	return nil
}

func (r *Renderer) renderFullFrame(frame canvas.Frame) error {
	r.out = append(r.out, ansi.CLEAR_SCREEN...)
	r.out = append(r.out, ansi.CURSOR_HOME...)

	for y := range frame.Height() {
		row, err := frame.RowAt(y)
		if err != nil {
			return err
		}
		previousRow, err := r.frame.RowAt(y)
		if err != nil {
			return err
		}
		r.renderRow(0, y, row)
		copy(previousRow, row)
	}
	return nil
}

func (r *Renderer) renderDiffFrame(frame canvas.Frame) error {
	for y := range frame.Height() {
		row, err := frame.RowAt(y)
		if err != nil {
			return err
		}
		previousRow, err := r.frame.RowAt(y)
		if err != nil {
			return err
		}
		for x := 0; x < len(row); {
			if row[x] == previousRow[x] {
				x++
				continue
			}
			start := x
			for x < len(row) && row[x] != previousRow[x] {
				x++
			}
			r.renderRow(start, y, row[start:x])
			copy(previousRow[start:x], row[start:x])
		}
	}
	return nil
}

func (r *Renderer) renderRow(x, y int, cells []canvas.Cell) {
	r.cursorMove(x, y)
	for _, cell := range cells {
		if cell.Width() == 0 {
			continue
		}
		r.renderStyle(cell)
		r.out = append(r.out, cell.Glyph()...)
	}
}

// cursorMove appends a CSI cursor-position command.
// Terminal coordinates are 1-based and ordered as row;column, so frame x,y
// becomes y+1;x+1. For example, x=9 y=4 appends "\x1b[5;10H".
func (r *Renderer) cursorMove(x, y int) {
	r.out = append(r.out, ansi.CSI...)
	r.out = strconv.AppendInt(r.out, int64(y+1), 10)
	r.out = append(r.out, ';')
	r.out = strconv.AppendInt(r.out, int64(x+1), 10)
	r.out = append(r.out, 'H')
}

func (r *Renderer) renderStyle(cell canvas.Cell) {
	fgChanged := !r.style.set || r.style.fg != cell.Foreground()
	if fgChanged {
		r.out = append(r.out, colorEscape(cell.Foreground())...)
		r.style.fg = cell.Foreground()
	}
	if !r.style.set || fgChanged || r.style.bg != cell.Background() {
		r.out = append(r.out, colorEscape(cell.Background())...)
		r.style.bg = cell.Background()
	}
	r.style.set = true
}

func (r *Renderer) flush() error {
	if len(r.out) == 0 {
		return nil
	}
	_, err := r.writer.Write(r.out)
	r.out = r.out[:0]
	return err
}
