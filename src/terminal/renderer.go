package terminal

import (
	"io"
	"strconv"
	"tuigo/ansi"
	"tuigo/domain"
	"unicode/utf8"
)

const estimatedCellBytes = 36

type renderStyle struct {
	fg, bg ansi.Color
	set    bool
}

type Renderer struct {
	frame, oldFrame domain.Frame
	fullRepaint     bool
	writer          io.Writer
	symbol          [utf8.UTFMax]byte
	out             []byte
	style           renderStyle
}

func NewRenderer(frame domain.Frame, writer io.Writer) Renderer {
	return Renderer{
		frame:       frame,
		fullRepaint: true,
		writer:      writer,
		out:         make([]byte, 0, frame.Height()*frame.Width()*estimatedCellBytes),
	}
}

func (r *Renderer) NextFrame(newFrame domain.Frame) error {
	if newFrame.Width() != r.frame.Width() ||
		newFrame.Height() != r.frame.Height() {
		r.frame = newFrame
		r.fullRepaint = true
		return nil
	}
	r.oldFrame = r.frame
	r.frame = newFrame
	return nil
}

func (r *Renderer) Render() error {
	r.style.set = false
	r.out = r.out[:0]
	if r.fullRepaint {
		if err := r.renderFullFrame(); err != nil {
			return err
		}
	} else {
		if err := r.renderDiffFrame(); err != nil {
			return err
		}
	}
	if len(r.out) > 0 {
		r.out = append(r.out, ansi.RESET...)
	}
	if err := r.flush(); err != nil {
		r.fullRepaint = true
		return err
	}
	r.fullRepaint = false
	r.oldFrame = r.frame
	return nil
}

func (r *Renderer) renderFullFrame() error {
	for y := range r.frame.Height() {
		for x := range r.frame.Width() {
			cell, err := r.frame.CellAt(x, y)
			if err != nil {
				return err
			}
			r.renderCell(x, y, cell)
		}
	}
	return nil
}

func (r *Renderer) renderDiffFrame() error {
	for y := range r.frame.Height() {
		for x := range r.frame.Width() {
			cell, err := r.frame.CellAt(x, y)
			if err != nil {
				return err
			}
			oldCell, err := r.oldFrame.CellAt(x, y)
			if err != nil {
				return err
			}
			if cell != oldCell {
				r.renderCell(x, y, cell)
			}
		}
	}
	return nil
}

func (r *Renderer) renderCell(x, y int, cell domain.Cell) {
	r.cursorMove(x, y)
	r.renderStyle(cell)
	n := utf8.EncodeRune(r.symbol[:], cell.Symbol())
	r.out = append(r.out, r.symbol[:n]...)
}

func (r *Renderer) renderStyle(cell domain.Cell) {
	if !r.style.set || r.style.fg != cell.Foreground() {
		r.out = append(r.out, cell.Foreground().String()...)
		r.style.fg = cell.Foreground()
	}
	if !r.style.set || r.style.bg != cell.Background() {
		r.out = append(r.out, cell.Background().String()...)
		r.style.bg = cell.Background()
	}
	r.style.set = true
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

func (r *Renderer) flush() error {
	if len(r.out) == 0 {
		return nil
	}
	_, err := r.writer.Write(r.out)
	r.out = r.out[:0]
	return err
}
