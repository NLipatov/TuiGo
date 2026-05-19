package domain

type Cell struct {
	fg, bg Color
	symbol rune
}

func NewCell(symbol rune, fg, bg Color) Cell {
	return Cell{
		symbol: symbol,
		fg:     fg,
		bg:     bg,
	}
}

func (c Cell) Foreground() Color {
	return c.fg
}

func (c Cell) Background() Color {
	return c.bg
}

func (c Cell) Symbol() rune {
	return c.symbol
}
