package tuigo

import (
	"io"
	"strconv"
	"tuigo/ansi"
	"unicode/utf8"
)

type RenderWriter interface {
	io.Writer
	io.StringWriter
}

type Renderer struct {
	frame, oldFrame Frame
	firstRender     bool
	writer          RenderWriter
	cursor          [32]byte
	symbol          [utf8.UTFMax]byte
}

func NewRenderer(frame Frame, writer RenderWriter) Renderer {
	return Renderer{
		frame:       frame,
		firstRender: true,
		writer:      writer,
	}
}

func (r *Renderer) NextFrame(newFrame Frame) error {
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
				if err := r.renderCell(x, y, cell); err != nil {
					return err
				}
			} else {
				oldCell, err := r.oldFrame.CellAt(x, y)
				if err != nil {
					return err
				}
				if cell != oldCell {
					if err := r.renderCell(x, y, cell); err != nil {
						return err
					}
				}
			}
		}
	}
	if r.firstRender {
		r.firstRender = false
	}
	r.oldFrame = r.frame
	return nil
}

func (r *Renderer) renderCell(x, y int, cell Cell) error {
	if err := r.moveCursor(x, y); err != nil {
		return err
	}
	if _, err := r.writer.WriteString(cell.Foreground().String()); err != nil {
		return err
	}
	if _, err := r.writer.WriteString(cell.Background().String()); err != nil {
		return err
	}

	n := utf8.EncodeRune(r.symbol[:], cell.Symbol())
	if _, err := r.writer.Write(r.symbol[:n]); err != nil {
		return err
	}

	_, err := r.writer.WriteString(string(ansi.RESET))
	return err
}

func (r *Renderer) moveCursor(x, y int) error {
	buf := r.cursor[:0]
	buf = append(buf, "\x1b["...)
	buf = strconv.AppendInt(buf, int64(y+1), 10)
	buf = append(buf, ';')
	buf = strconv.AppendInt(buf, int64(x+1), 10)
	buf = append(buf, 'H')

	_, err := r.writer.Write(buf)
	return err
}
