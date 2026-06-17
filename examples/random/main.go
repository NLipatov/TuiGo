package main

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NLipatov/tuigo/ansi"
	"github.com/NLipatov/tuigo/core"
	"github.com/NLipatov/tuigo/terminal"
	"github.com/NLipatov/tuigo/terminal/input"
)

const (
	asciiMin = 40
	asciiMax = 125
)

var fgEscapeSequences = []ansi.ANSIEscapeSequence{
	ansi.FG_RED,
	ansi.FG_GREEN,
	ansi.FG_YELLOW,
	ansi.FG_BLUE,
	ansi.FG_PURPLE,
	ansi.FG_CYAN,
	ansi.FG_WHITE,
	ansi.FG_BOLD_BLACK,
	ansi.FG_BOLD_RED,
	ansi.FG_BOLD_GREEN,
	ansi.FG_BOLD_YELLOW,
	ansi.FG_BOLD_BLUE,
	ansi.FG_BOLD_PURPLE,
	ansi.FG_BOLD_CYAN,
	ansi.FG_BOLD_WHITE,
}

type demo struct {
	idx           int
	width         int
	height        int
	frameCount    int
	fps           int
	lastFPSUpdate time.Time
	frames        [2]core.Frame
	cells         [2][]core.Cell
	cellVariants  []core.Cell
	headerGlyphs  [128]core.Cell
	headerCells   []core.Cell
	rng           *rand.Rand
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session, err := terminal.NewSession(ctx, os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()

	events, err := session.Start()
	if err != nil {
		panic(err)
	}

	demo, err := newDemo(session)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-events:
			if !ok {
				if ctx.Err() != nil {
					return
				}
				panic("events closed")
			}
			switch event.Type {
			case terminal.EventKey:
				if quitRequested(event.Key) {
					return
				}
			case terminal.EventResize:
				demo, err = newDemo(session)
				if err != nil {
					panic(err)
				}
			case terminal.EventError:
				panic(event.Err)
			}
		default:
			if err := demo.tick(&session); err != nil {
				panic(err)
			}
		}
	}
}

func newDemo(session terminal.Session) (*demo, error) {
	width, height, err := session.Size()
	if err != nil {
		return nil, err
	}

	fg, err := ansi.NewColor(ansi.FG_WHITE)
	if err != nil {
		return nil, err
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		return nil, err
	}
	blank, err := core.NewCellWithWidth(" ", 1, fg, bg)
	if err != nil {
		return nil, err
	}
	frames, cells, err := newFrames(width, height, blank)
	if err != nil {
		return nil, err
	}
	variants, err := newCellVariants(bg)
	if err != nil {
		return nil, err
	}
	headerFG, err := ansi.NewColor(ansi.FG_GREEN)
	if err != nil {
		return nil, err
	}
	headerBG, err := ansi.NewColor(ansi.BG_HIGH_INTENSITY_YELLOW)
	if err != nil {
		return nil, err
	}
	headerGlyphs, err := newHeaderGlyphs(headerFG, headerBG)
	if err != nil {
		return nil, err
	}

	demo := &demo{
		width:         width,
		height:        height,
		frames:        frames,
		cells:         cells,
		cellVariants:  variants,
		headerGlyphs:  headerGlyphs,
		lastFPSUpdate: time.Now(),
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	demo.updateHeader()
	demo.drawHeader(0)
	demo.drawHeader(1)
	return demo, nil
}

func newFrames(width, height int, blank core.Cell) ([2]core.Frame, [2][]core.Cell, error) {
	var frames [2]core.Frame
	var cells [2][]core.Cell

	for idx := range frames {
		frame, frameCells, err := newFrame(width, height, blank)
		if err != nil {
			return [2]core.Frame{}, [2][]core.Cell{}, err
		}
		frames[idx] = frame
		cells[idx] = frameCells
	}
	return frames, cells, nil
}

func newFrame(width, height int, blank core.Cell) (core.Frame, []core.Cell, error) {
	cells := newCells(width*height, blank)
	frame, err := core.NewFrame(width, height, cells)
	if err != nil {
		return core.Frame{}, nil, err
	}
	return frame, cells, nil
}

func newCells(size int, blank core.Cell) []core.Cell {
	cells := make([]core.Cell, size)
	for idx := range cells {
		cells[idx] = blank
	}
	return cells
}

func newCellVariants(bg ansi.Color) ([]core.Cell, error) {
	fgPalette, err := newFGPalette()
	if err != nil {
		return nil, err
	}

	variants := make([]core.Cell, 0, len(fgPalette)*(asciiMax-asciiMin+1))
	for _, fg := range fgPalette {
		for r := asciiMin; r <= asciiMax; r++ {
			cell, err := core.NewCellWithWidth(string(rune(r)), 1, fg, bg)
			if err != nil {
				return nil, err
			}
			variants = append(variants, cell)
		}
	}
	return variants, nil
}

func newFGPalette() ([]ansi.Color, error) {
	palette := make([]ansi.Color, 0, len(fgEscapeSequences))
	for _, sequence := range fgEscapeSequences {
		fg, err := ansi.NewColor(sequence)
		if err != nil {
			return nil, err
		}
		palette = append(palette, fg)
	}
	return palette, nil
}

func newHeaderGlyphs(fg, bg ansi.Color) ([128]core.Cell, error) {
	var glyphs [128]core.Cell
	for ch := byte(32); ch <= 126; ch++ {
		cell, err := core.NewCellWithWidth(string(rune(ch)), 1, fg, bg)
		if err != nil {
			return [128]core.Cell{}, err
		}
		glyphs[ch] = cell
	}
	return glyphs, nil
}

func (d *demo) tick(session *terminal.Session) error {
	cellIdx, hasCell := d.randomCellIndex()
	var cell core.Cell
	if hasCell {
		cell = d.cellVariants[d.rng.Intn(len(d.cellVariants))]
		d.cells[d.idx][cellIdx] = cell
	}

	headerUpdated := d.updateFPS()
	if headerUpdated {
		d.drawHeader(d.idx)
	}
	if err := session.Render(d.frames[d.idx]); err != nil {
		return err
	}

	d.idx ^= 1
	if hasCell {
		d.cells[d.idx][cellIdx] = cell
	}
	if headerUpdated {
		d.drawHeader(d.idx)
	}
	return nil
}

func (d *demo) randomCellIndex() (int, bool) {
	total := d.width * d.height
	headerWidth := min(len(d.headerCells), d.width)
	if total <= headerWidth {
		return 0, false
	}
	return headerWidth + d.rng.Intn(total-headerWidth), true
}

func (d *demo) updateFPS() bool {
	d.frameCount++
	now := time.Now()
	if now.Sub(d.lastFPSUpdate) < time.Second {
		return false
	}

	d.fps, d.frameCount = d.frameCount, 0
	d.lastFPSUpdate = now
	d.updateHeader()
	return true
}

func (d *demo) updateHeader() {
	headerText := "q/esc to quit | FPS: " + fpsLabel(d.fps)
	if cap(d.headerCells) < len(headerText) {
		d.headerCells = make([]core.Cell, 0, len(headerText))
	}
	d.headerCells = d.headerCells[:0]
	for idx := range len(headerText) {
		d.headerCells = append(d.headerCells, d.headerGlyphs[headerText[idx]])
	}
}

func (d *demo) drawHeader(frameIdx int) {
	header := d.headerCells
	if len(header) > d.width {
		header = header[:d.width]
	}
	copy(d.cells[frameIdx], header)
}

func fpsLabel(fps int) string {
	const width = 9
	const cap = 99999999
	if fps <= 0 {
		return padLeft("--", width)
	}
	if fps > cap {
		return ">" + strings.Repeat("9", width-1)
	}
	return padLeft(strconv.Itoa(fps), width)
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func quitRequested(event input.KeyEvent) bool {
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
