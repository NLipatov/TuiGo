package terminal

import (
	"io"
	"strconv"
	"tuigo/ansi"
	"tuigo/domain"
	"unicode/utf8"
)

type Renderer struct {
	frame, oldFrame domain.Frame
	firstRender     bool
	writer          io.Writer
	symbol          [utf8.UTFMax]byte
	out             []byte
}

func NewRenderer(frame domain.Frame, writer io.Writer) Renderer {
	return Renderer{
		frame:       frame,
		firstRender: true,
		writer:      writer,
	}
}

func (r *Renderer) NextFrame(newFrame domain.Frame) error {
	if newFrame.Width() != r.frame.Width() ||
		newFrame.Height() != r.frame.Height() {
		r.frame = newFrame
		r.firstRender = true
		return nil
	}
	r.oldFrame = r.frame
	r.frame = newFrame
	return nil
}

func (r *Renderer) Render() error {
	for y := range r.frame.Height() {
		for x := range r.frame.Width() {
			cell, err := r.frame.CellAt(x, y)
			if err != nil {
				return err
			}
			if r.firstRender {
				r.renderCell(x, y, cell)
			} else {
				oldCell, err := r.oldFrame.CellAt(x, y)
				if err != nil {
					return err
				}
				if cell != oldCell {
					r.renderCell(x, y, cell)
				}
			}
		}
	}
	if r.firstRender {
		r.firstRender = false
	}
	r.oldFrame = r.frame
	if len(r.out) > 0 {
		if _, err := r.writer.Write(r.out); err != nil {
			return err
		}
		r.out = r.out[:0]
	}
	return nil
}

func (r *Renderer) renderCell(x, y int, cell domain.Cell) {
	r.moveCursor(x, y)
	r.out = append(r.out, cell.Foreground().String()...)
	r.out = append(r.out, cell.Background().String()...)
	n := utf8.EncodeRune(r.symbol[:], cell.Symbol())
	r.out = append(r.out, r.symbol[:n]...)
	r.out = append(r.out, ansi.RESET...)
}

// moveCursor writes a CSI cursor-position command.
// Terminal coordinates are 1-based and ordered as row;column, so frame x,y
// becomes y+1;x+1. For example, x=9 y=4 writes "\x1b[5;10H".
func (r *Renderer) moveCursor(x, y int) {
	r.out = append(r.out, ansi.CSI...)
	r.out = strconv.AppendInt(r.out, int64(y+1), 10)
	r.out = append(r.out, ';')
	r.out = strconv.AppendInt(r.out, int64(x+1), 10)
	r.out = append(r.out, 'H')
}
