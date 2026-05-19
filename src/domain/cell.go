package domain

type Cell struct {
	fg, bg Color
	symbol rune
}

func NewCell(fg, bg Color, symbol rune) Cell {
	return Cell{
		fg:     fg,
		bg:     bg,
		symbol: symbol,
	}
}

func (c *Cell) Foreground() Color {
	return c.fg
}

func (c *Cell) Background() Color {
	return c.bg
}

func (c *Cell) Symbol() rune {
	return c.symbol
}
