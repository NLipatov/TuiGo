package benchmarks

import (
	"bytes"
	"io"
	"strings"
	"testing"

	vaxis "git.sr.ht/~rockorager/vaxis"
	"github.com/NLipatov/tuigo/ansi"
	"github.com/NLipatov/tuigo/core"
	"github.com/NLipatov/tuigo/terminal/render"
	"github.com/gdamore/tcell/v2"
)

const (
	benchWidth  = 120
	benchHeight = 40
)

type workload struct {
	name          string
	width         int
	height        int
	positions     []position
	baseSymbol    rune
	changedSymbol rune
	styleChange   bool
}

type position struct {
	x   int
	y   int
	idx int
}

func BenchmarkRenderer(b *testing.B) {
	for _, workload := range benchmarkWorkloads() {
		b.Run("tuigo/"+workload.name, func(b *testing.B) {
			benchmarkTuigoRenderer(b, workload)
		})
		b.Run("tcell/"+workload.name, func(b *testing.B) {
			benchmarkTcellSimulationScreen(b, workload)
		})
		b.Run("vaxis/"+workload.name, func(b *testing.B) {
			benchmarkVaxisRenderer(b, workload)
		})
	}
}

func benchmarkWorkloads() []workload {
	return []workload{
		{
			name:          "no-change/120x40",
			width:         benchWidth,
			height:        benchHeight,
			baseSymbol:    'x',
			changedSymbol: 'x',
		},
		{
			name:          "full-content-change/120x40",
			width:         benchWidth,
			height:        benchHeight,
			positions:     allPositions(benchWidth, benchHeight),
			baseSymbol:    'x',
			changedSymbol: 'y',
		},
		{
			name:          "one-cell-change/120x40",
			width:         benchWidth,
			height:        benchHeight,
			positions:     []position{{x: benchWidth - 1, y: benchHeight - 1, idx: benchWidth*benchHeight - 1}},
			baseSymbol:    'x',
			changedSymbol: 'y',
		},
		{
			name:          "run-40-change/120x40",
			width:         benchWidth,
			height:        benchHeight,
			positions:     runPositions(benchWidth, benchHeight, 40),
			baseSymbol:    'x',
			changedSymbol: 'y',
		},
		{
			name:          "random-5pct-change/120x40",
			width:         benchWidth,
			height:        benchHeight,
			positions:     randomPositions(benchWidth, benchHeight, benchWidth*benchHeight/20),
			baseSymbol:    'x',
			changedSymbol: 'y',
		},
		{
			name:          "full-style-change/120x40",
			width:         benchWidth,
			height:        benchHeight,
			positions:     allPositions(benchWidth, benchHeight),
			baseSymbol:    'x',
			changedSymbol: 'x',
			styleChange:   true,
		},
		{
			name:          "unicode-content-change/120x40",
			width:         benchWidth,
			height:        benchHeight,
			positions:     allPositions(benchWidth, benchHeight),
			baseSymbol:    'x',
			changedSymbol: 'Ж',
		},
	}
}

func benchmarkTuigoRenderer(b *testing.B, workload workload) {
	baseCell := tuigoCell(b, workload.baseSymbol, false)
	changedCell := tuigoCell(b, workload.changedSymbol, workload.styleChange)
	baseCells := make([]core.Cell, workload.width*workload.height)
	changedCells := make([]core.Cell, workload.width*workload.height)
	fillTuigoCells(baseCells, baseCell)
	fillTuigoCells(changedCells, baseCell)

	baseFrame := tuigoFrame(b, workload.width, workload.height, baseCells)
	changedFrame := tuigoFrame(b, workload.width, workload.height, changedCells)
	writer := countingWriter{}
	renderer := render.NewRenderer(&writer)
	if err := renderer.Render(baseFrame); err != nil {
		b.Fatalf("Render() error = %v", err)
	}

	writer.Reset()
	b.ReportAllocs()
	b.ResetTimer()

	useChanged := false
	for b.Loop() {
		useChanged = !useChanged
		if useChanged {
			applyTuigoCells(changedCells, workload.positions, changedCell)
			if err := renderer.Render(changedFrame); err != nil {
				b.Fatalf("Render() error = %v", err)
			}
			continue
		}
		applyTuigoCells(baseCells, workload.positions, baseCell)
		if err := renderer.Render(baseFrame); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
	if b.N > 0 {
		reportWorkload(b, workload)
		b.ReportMetric(float64(writer.Bytes())/float64(b.N), "bytes_out/op")
	}
}

func benchmarkTcellSimulationScreen(b *testing.B, workload workload) {
	baseStyle := tcellStyle(false)
	changedStyle := tcellStyle(workload.styleChange)
	baseText := string(workload.baseSymbol)
	changedText := string(workload.changedSymbol)
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		b.Fatalf("SimulationScreen.Init() error = %v", err)
	}
	defer screen.Fini()
	screen.SetSize(workload.width, workload.height)
	paintTcellScreen(screen, workload.width, workload.height, baseText, baseStyle)
	screen.Show()

	b.ReportAllocs()
	b.ResetTimer()

	useChanged := false
	for b.Loop() {
		useChanged = !useChanged
		text := baseText
		style := baseStyle
		if useChanged {
			text = changedText
			style = changedStyle
		}
		applyTcellCells(screen, workload, text, style)
		screen.Show()
	}
	reportWorkload(b, workload)
}

func benchmarkVaxisRenderer(b *testing.B, workload workload) {
	console := newVaxisBenchmarkConsole(workload.width, workload.height)
	vx, err := vaxis.New(vaxis.Options{
		DisableKittyKeyboard: true,
		DisableMouse:         true,
		NoSignals:            true,
		WithConsole:          console,
	})
	if err != nil {
		b.Fatalf("vaxis.New() error = %v", err)
	}
	defer vx.Close()

	baseCell := vaxisCell(workload.baseSymbol, false)
	changedCell := vaxisCell(workload.changedSymbol, workload.styleChange)
	window := vx.Window()
	window.Fill(baseCell)
	vx.Render()

	console.ResetOutput()
	b.ReportAllocs()
	b.ResetTimer()

	useChanged := false
	for b.Loop() {
		useChanged = !useChanged
		cell := baseCell
		if useChanged {
			cell = changedCell
		}
		applyVaxisCells(window, workload, cell)
		vx.Render()
	}
	reportWorkload(b, workload)
	if b.N > 0 {
		b.ReportMetric(float64(console.Bytes())/float64(b.N), "bytes_out/op")
	}
}

func reportWorkload(b *testing.B, workload workload) {
	b.ReportMetric(float64(workload.width*workload.height), "cells/op")
	b.ReportMetric(float64(len(workload.positions)), "dirty_cells/op")
}

func allPositions(width, height int) []position {
	positions := make([]position, 0, width*height)
	for y := range height {
		for x := range width {
			positions = append(positions, position{
				x:   x,
				y:   y,
				idx: y*width + x,
			})
		}
	}
	return positions
}

func runPositions(width, height, count int) []position {
	y := height / 2
	start := (width - count) / 2
	positions := make([]position, 0, count)
	for x := start; x < start+count; x++ {
		positions = append(positions, position{
			x:   x,
			y:   y,
			idx: y*width + x,
		})
	}
	return positions
}

func randomPositions(width, height, count int) []position {
	positions := make([]position, 0, count)
	seen := make(map[int]struct{}, count)
	seed := uint32(1)
	for len(positions) < count {
		seed = seed*1664525 + 1013904223
		idx := int(seed % uint32(width*height))
		if _, ok := seen[idx]; ok {
			continue
		}
		seen[idx] = struct{}{}
		positions = append(positions, position{
			x:   idx % width,
			y:   idx / width,
			idx: idx,
		})
	}
	return positions
}

func fillTuigoCells(cells []core.Cell, cell core.Cell) {
	for i := range cells {
		cells[i] = cell
	}
}

func applyTuigoCells(cells []core.Cell, positions []position, cell core.Cell) {
	for _, pos := range positions {
		cells[pos.idx] = cell
	}
}

func paintTcellScreen(screen tcell.Screen, width, height int, text string, style tcell.Style) {
	line := strings.Repeat(text, width)
	for y := range height {
		screen.PutStrStyled(0, y, line, style)
	}
}

func applyTcellCells(screen tcell.Screen, workload workload, text string, style tcell.Style) {
	if len(workload.positions) == 0 {
		return
	}
	if len(workload.positions) == workload.width*workload.height {
		paintTcellScreen(screen, workload.width, workload.height, text, style)
		return
	}
	if start, y, ok := contiguousRowRun(workload.positions); ok {
		screen.PutStrStyled(start, y, strings.Repeat(text, len(workload.positions)), style)
		return
	}
	for _, pos := range workload.positions {
		screen.Put(pos.x, pos.y, text, style)
	}
}

func applyVaxisCells(window vaxis.Window, workload workload, cell vaxis.Cell) {
	if len(workload.positions) == 0 {
		return
	}
	if len(workload.positions) == workload.width*workload.height {
		window.Fill(cell)
		return
	}
	for _, pos := range workload.positions {
		window.SetCell(pos.x, pos.y, cell)
	}
}

func contiguousRowRun(positions []position) (int, int, bool) {
	if len(positions) == 0 {
		return 0, 0, false
	}
	start := positions[0].x
	y := positions[0].y
	for i, pos := range positions {
		if pos.y != y || pos.x != start+i {
			return 0, 0, false
		}
	}
	return start, y, true
}

func tuigoFrame(b *testing.B, width, height int, cells []core.Cell) core.Frame {
	b.Helper()

	frame, err := core.NewFrame(width, height, cells)
	if err != nil {
		b.Fatalf("core.NewFrame() error = %v", err)
	}
	return frame
}

func tuigoCell(b *testing.B, symbol rune, changedStyle bool) core.Cell {
	b.Helper()

	fgSequence := ansi.FG_RED
	if changedStyle {
		fgSequence = ansi.FG_GREEN
	}
	fg, err := ansi.NewColor(fgSequence)
	if err != nil {
		b.Fatalf("ansi.NewColor(%q) error = %v", fgSequence, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		b.Fatalf("ansi.NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}
	return core.NewCell(symbol, fg, bg)
}

func tcellStyle(changed bool) tcell.Style {
	fg := tcell.ColorRed
	if changed {
		fg = tcell.ColorGreen
	}
	return tcell.StyleDefault.Foreground(fg).Background(tcell.ColorBlack)
}

func vaxisCell(symbol rune, changedStyle bool) vaxis.Cell {
	fg := vaxis.IndexColor(1)
	if changedStyle {
		fg = vaxis.IndexColor(2)
	}
	return vaxis.Cell{
		Character: vaxis.Character{
			Grapheme: string(symbol),
			Width:    1,
		},
		Style: vaxis.Style{
			Foreground: fg,
			Background: vaxis.IndexColor(0),
		},
	}
}

type countingWriter struct {
	bytes int64
}

func (w *countingWriter) Write(p []byte) (int, error) {
	w.bytes += int64(len(p))
	return len(p), nil
}

func (w *countingWriter) Bytes() int64 {
	return w.bytes
}

func (w *countingWriter) Reset() {
	w.bytes = 0
}

type vaxisBenchmarkConsole struct {
	input  *bytes.Reader
	output countingWriter
	width  int
	height int
}

func newVaxisBenchmarkConsole(width, height int) *vaxisBenchmarkConsole {
	return &vaxisBenchmarkConsole{
		input:  bytes.NewReader([]byte("\x1b[?62;4;22c")),
		width:  width,
		height: height,
	}
}

func (c *vaxisBenchmarkConsole) Read(p []byte) (int, error) {
	if c.input.Len() > 0 {
		return c.input.Read(p)
	}
	return 0, io.EOF
}

func (c *vaxisBenchmarkConsole) Write(p []byte) (int, error) {
	return c.output.Write(p)
}

func (c *vaxisBenchmarkConsole) Fd() uintptr {
	return ^uintptr(0)
}

func (c *vaxisBenchmarkConsole) SetRaw() error {
	return nil
}

func (c *vaxisBenchmarkConsole) Reset() error {
	return nil
}

func (c *vaxisBenchmarkConsole) ResetOutput() {
	c.output.Reset()
}

func (c *vaxisBenchmarkConsole) Size() (int, int, int, int, error) {
	return c.width, c.height, 0, 0, nil
}

func (c *vaxisBenchmarkConsole) Close() error {
	return nil
}

func (c *vaxisBenchmarkConsole) Bytes() int64 {
	return c.output.Bytes()
}
