package main

import (
	"context"
	"os"

	"github.com/NLipatov/tuigo/ansi"
	"github.com/NLipatov/tuigo/core"
	"github.com/NLipatov/tuigo/terminal"
	"github.com/NLipatov/tuigo/terminal/input"
)

type palette struct {
	blank  core.Cell
	logo   core.Cell
	shadow core.Cell
	title  ansi.Color
	hint   ansi.Color
	bg     ansi.Color
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session, err := terminal.NewSession(ctx, os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}

	events, err := session.Start()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	width, height, err := session.Size()
	if err != nil {
		panic(err)
	}
	if err := renderHello(&session, width, height); err != nil {
		panic(err)
	}

	for event := range events {
		switch event.Type {
		case terminal.EventKey:
			if quitRequested(event.Key) {
				cancel()
				return
			}
		case terminal.EventResize:
			width = event.Resize.Width
			height = event.Resize.Height
			if err := renderHello(&session, width, height); err != nil {
				panic(err)
			}
		case terminal.EventError:
			panic(event.Err)
		}
	}
}

func quitRequested(event input.Event) bool {
	if event.Code == input.KeyEsc {
		return true
	}
	if event.Code != input.KeyRune {
		return false
	}
	if event.Text == "q" {
		return true
	}
	return event.Text == "c" && event.Mod&input.ModCtrl != 0
}

func renderHello(session *terminal.Session, width, height int) error {
	frame, err := helloFrame(width, height)
	if err != nil {
		return err
	}
	return session.Render(frame)
}

func helloFrame(width, height int) (core.Frame, error) {
	colors := newPalette()
	cells := make([]core.Cell, width*height)
	for i := range cells {
		cells[i] = colors.blank
	}

	logoX := max(0, (width-logoWidth())/2)
	logoY := max(3, (height-logoHeight())/2-1)
	drawLogo(cells, width, height, logoX+1, logoY+1, colors.shadow)
	drawLogo(cells, width, height, logoX, logoY, colors.logo)
	drawText(cells, width, height, centeredX(width, helloTitle), logoY+logoHeight()+2, helloTitle, colors.title, colors.bg)
	drawText(cells, width, height, centeredX(width, "Press q, Esc, or Ctrl+C to quit"), height-2, "Press q, Esc, or Ctrl+C to quit", colors.hint, colors.bg)

	return core.NewFrame(width, height, cells)
}

func newPalette() palette {
	bg := mustColor(ansi.BG_BLACK)
	return palette{
		blank:  core.NewCell(' ', mustColor(ansi.FG_WHITE), bg),
		logo:   core.NewCell('=', mustColor(ansi.FG_BOLD_GREEN), bg),
		shadow: core.NewCell('.', mustColor(ansi.FG_HIGH_INTENSITY_BLACK), bg),
		title:  mustColor(ansi.FG_GREEN),
		hint:   mustColor(ansi.FG_HIGH_INTENSITY_BLACK),
		bg:     bg,
	}
}

func drawLogo(cells []core.Cell, width, height, left, top int, cell core.Cell) {
	for y, row := range helloLogo {
		for x, pixel := range row {
			if pixel != ' ' {
				putCell(cells, width, height, left+x, top+y, cell)
			}
		}
	}
}

func drawText(cells []core.Cell, width, height, left, y int, text string, fg, bg ansi.Color) {
	for x, char := range text {
		putCell(cells, width, height, left+x, y, core.NewCell(char, fg, bg))
	}
}

func putCell(cells []core.Cell, width, height, x, y int, cell core.Cell) {
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	cells[y*width+x] = cell
}

func centeredX(width int, text string) int {
	return max(0, (width-len(text))/2)
}

func logoWidth() int {
	width := 0
	for _, row := range helloLogo {
		width = max(width, len(row))
	}
	return width
}

func logoHeight() int {
	return len(helloLogo)
}

func mustColor(sequence ansi.ANSIEscapeSequence) ansi.Color {
	color, err := ansi.NewColor(sequence)
	if err != nil {
		panic(err)
	}
	return color
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var helloLogo = []string{
	"##  ## ###### ##     ##      ####",
	"##  ## ##     ##     ##     ##  ##",
	"##  ## ##     ##     ##     ##  ##",
	"###### #####  ##     ##     ##  ##",
	"##  ## ##     ##     ##     ##  ##",
	"##  ## ##     ##     ##     ##  ##",
	"##  ## ###### ###### ######  ####",
}

const helloTitle = "tuigo"
