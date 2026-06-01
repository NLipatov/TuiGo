package core

import "github.com/NLipatov/tuigo/ansi"

type Cell struct {
	fg, bg ansi.Color
	symbol rune
}

func NewCell(symbol rune, fg, bg ansi.Color) Cell {
	return Cell{
		symbol: symbol,
		fg:     fg,
		bg:     bg,
	}
}

func (c Cell) Foreground() ansi.Color {
	return c.fg
}

func (c Cell) Background() ansi.Color {
	return c.bg
}

func (c Cell) Symbol() rune {
	return c.symbol
}
