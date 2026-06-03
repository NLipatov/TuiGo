package render

import (
	"bytes"
	"errors"
	"testing"

	"github.com/NLipatov/tuigo/ansi"
	"github.com/NLipatov/tuigo/core"
)

func fullRepaintPrefix() string {
	return string(ansi.CLEAR_SCREEN) + string(ansi.CURSOR_HOME)
}

func TestRendererRenderWritesCell(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	frame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(&out)

	if err := renderer.Render(frame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := fullRepaintPrefix() + "\x1b[1;1H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "x"
	if got := out.String(); got != want {
		t.Fatalf("rendered output = %q, want %q", got, want)
	}
}

func TestRendererRenderWritesNothingWhenFrameIsUnchanged(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	frame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(&out)
	if err := renderer.Render(frame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.Render(frame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if got := out.String(); got != "" {
		t.Fatalf("rendered output = %q, want empty output", got)
	}
}

func TestRendererRenderWritesOnlyChangedCell(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	firstFrame, err := core.NewFrame(2, 1, []core.Cell{
		core.NewCell('x', fg, bg),
		core.NewCell('y', fg, bg),
	})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}
	nextFrame, err := core.NewFrame(2, 1, []core.Cell{
		core.NewCell('x', fg, bg),
		core.NewCell('z', fg, bg),
	})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(&out)
	if err := renderer.Render(firstFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.Render(nextFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "\x1b[1;2H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "z"
	if got := out.String(); got != want {
		t.Fatalf("rendered output = %q, want %q", got, want)
	}
}

func TestRendererRenderWritesCellWhenOnlyStyleChanges(t *testing.T) {
	red, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	green, err := ansi.NewColor(ansi.FG_GREEN)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_GREEN, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	firstFrame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', red, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}
	nextFrame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', green, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(&out)
	if err := renderer.Render(firstFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.Render(nextFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "\x1b[1;1H" + string(ansi.FG_GREEN) + string(ansi.BG_BLACK) + "x"
	if got := out.String(); got != want {
		t.Fatalf("rendered output = %q, want %q", got, want)
	}
}

func TestRendererRenderReappliesBackgroundAfterForegroundReset(t *testing.T) {
	title, err := ansi.NewColor(ansi.FG_BOLD_WHITE)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_BOLD_WHITE, err)
	}
	blank, err := ansi.NewColor(ansi.FG_WHITE)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_WHITE, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}
	frame, err := core.NewFrame(2, 1, []core.Cell{
		core.NewCell('t', title, bg),
		core.NewCell(' ', blank, bg),
	})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(&out)
	if err := renderer.Render(frame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := fullRepaintPrefix() + "\x1b[1;1H" + string(ansi.FG_BOLD_WHITE) + string(ansi.BG_BLACK) + "t" + string(ansi.FG_WHITE) + string(ansi.BG_BLACK) + " "
	if got := out.String(); got != want {
		t.Fatalf("rendered output = %q, want %q", got, want)
	}
}

func TestRendererRenderWritesAdjacentChangedCellsAsSingleRun(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}
	cell := func(symbol rune) core.Cell {
		return core.NewCell(symbol, fg, bg)
	}

	firstFrame, err := core.NewFrame(4, 1, []core.Cell{
		cell('a'),
		cell('b'),
		cell('c'),
		cell('d'),
	})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}
	nextFrame, err := core.NewFrame(4, 1, []core.Cell{
		cell('a'),
		cell('x'),
		cell('y'),
		cell('d'),
	})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(&out)
	if err := renderer.Render(firstFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.Render(nextFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "\x1b[1;2H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "xy"
	if got := out.String(); got != want {
		t.Fatalf("rendered output = %q, want %q", got, want)
	}
}

func TestRendererRenderWritesSeparatedChangedCellsAsSeparateRuns(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}
	cell := func(symbol rune) core.Cell {
		return core.NewCell(symbol, fg, bg)
	}

	firstFrame, err := core.NewFrame(4, 1, []core.Cell{
		cell('a'),
		cell('b'),
		cell('c'),
		cell('d'),
	})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}
	nextFrame, err := core.NewFrame(4, 1, []core.Cell{
		cell('a'),
		cell('x'),
		cell('c'),
		cell('y'),
	})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(&out)
	if err := renderer.Render(firstFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.Render(nextFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "\x1b[1;2H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "x" + "\x1b[1;4H" + "y"
	if got := out.String(); got != want {
		t.Fatalf("rendered output = %q, want %q", got, want)
	}
}

func TestRendererRenderWritesFullFrameAfterResize(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	firstFrame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}
	nextFrame, err := core.NewFrame(2, 1, []core.Cell{
		core.NewCell('y', fg, bg),
		core.NewCell('z', fg, bg),
	})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(&out)
	if err := renderer.Render(firstFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.Render(nextFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := fullRepaintPrefix() + "\x1b[1;1H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "yz"
	if got := out.String(); got != want {
		t.Fatalf("rendered output = %q, want %q", got, want)
	}
}

func TestRendererRenderRetriesFullFrameAfterWriteError(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	frame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	writeErr := errors.New("write failed")
	writer := failOnceWriter{err: writeErr}
	renderer := NewRenderer(&writer)

	if err := renderer.Render(frame); !errors.Is(err, writeErr) {
		t.Fatalf("Render() error = %v, want %v", err, writeErr)
	}

	if err := renderer.Render(frame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := fullRepaintPrefix() + "\x1b[1;1H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "x"
	if got := writer.out.String(); got != want {
		t.Fatalf("rendered output = %q, want %q", got, want)
	}
}

func TestRendererRenderDoesNotAllocateWhenFrameIsUnchanged(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	frame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	renderer := NewRenderer(discardWriter{})
	if err := renderer.Render(frame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = renderer.Render(frame)
	})
	if allocs != 0 {
		t.Fatalf("allocations per render = %.2f, want 0", allocs)
	}
}

func TestRendererRenderDoesNotAllocateWhenRenderingFullFrame(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	frame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	renderer := NewRenderer(discardWriter{})
	if err := renderer.Render(frame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	allocs := testing.AllocsPerRun(1000, func() {
		renderer.fullRepaint = true
		_ = renderer.Render(frame)
	})
	if allocs != 0 {
		t.Fatalf("allocations per full render = %.2f, want 0", allocs)
	}
}

func TestRendererRenderDoesNotAllocateWhenRenderingChangedCell(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	firstFrame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}
	nextFrame, err := core.NewFrame(1, 1, []core.Cell{core.NewCell('y', fg, bg)})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	renderer := NewRenderer(discardWriter{})
	if err := renderer.Render(firstFrame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = renderer.Render(nextFrame)
		_ = renderer.Render(firstFrame)
	})
	if allocs != 0 {
		t.Fatalf("allocations per changed-cell render = %.2f, want 0", allocs)
	}
}

func BenchmarkRendererRenderUnchangedFrame(b *testing.B) {
	frame := benchmarkFrame(b, 80, 24, 'x')
	renderer := NewRenderer(discardWriter{})
	if err := renderer.Render(frame); err != nil {
		b.Fatalf("Render() error = %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := renderer.Render(frame); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
}

func BenchmarkRendererRenderFullFrame(b *testing.B) {
	frame := benchmarkFrame(b, 80, 24, 'x')
	renderer := NewRenderer(discardWriter{})
	if err := renderer.Render(frame); err != nil {
		b.Fatalf("Render() error = %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		renderer.fullRepaint = true
		if err := renderer.Render(frame); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
}

func BenchmarkRendererRenderChangedCell(b *testing.B) {
	firstFrame := benchmarkFrame(b, 80, 24, 'x')
	nextFrame := benchmarkFrameWithLastCell(b, 80, 24, 'x', 'y')

	renderer := NewRenderer(discardWriter{})
	if err := renderer.Render(firstFrame); err != nil {
		b.Fatalf("Render() error = %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := renderer.Render(nextFrame); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
		if err := renderer.Render(firstFrame); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
}

func BenchmarkRendererRenderChangedRun(b *testing.B) {
	firstFrame := benchmarkFrame(b, 80, 24, 'x')
	nextFrame := benchmarkFrameWithRun(b, 80, 24, 'x', 'y', 1860, 40)

	renderer := NewRenderer(discardWriter{})
	if err := renderer.Render(firstFrame); err != nil {
		b.Fatalf("Render() error = %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := renderer.Render(nextFrame); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
		if err := renderer.Render(firstFrame); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
}

func benchmarkFrame(b *testing.B, width, height int, symbol rune) core.Frame {
	b.Helper()

	cells := make([]core.Cell, width*height)
	cell := benchmarkCell(b, symbol)
	for i := range cells {
		cells[i] = cell
	}

	frame, err := core.NewFrame(width, height, cells)
	if err != nil {
		b.Fatalf("core.NewFrame() error = %v", err)
	}
	return frame
}

func benchmarkFrameWithLastCell(b *testing.B, width, height int, symbol, lastSymbol rune) core.Frame {
	b.Helper()

	cells := make([]core.Cell, width*height)
	cell := benchmarkCell(b, symbol)
	for i := range cells {
		cells[i] = cell
	}
	cells[len(cells)-1] = benchmarkCell(b, lastSymbol)

	frame, err := core.NewFrame(width, height, cells)
	if err != nil {
		b.Fatalf("core.NewFrame() error = %v", err)
	}
	return frame
}

func benchmarkFrameWithRun(b *testing.B, width, height int, symbol, runSymbol rune, start, count int) core.Frame {
	b.Helper()

	cells := make([]core.Cell, width*height)
	if start < 0 || count < 0 || start+count > len(cells) {
		b.Fatalf("invalid run bounds: start=%d count=%d len=%d", start, count, len(cells))
	}
	cell := benchmarkCell(b, symbol)
	for i := range cells {
		cells[i] = cell
	}
	runCell := benchmarkCell(b, runSymbol)
	for i := range count {
		cells[start+i] = runCell
	}

	frame, err := core.NewFrame(width, height, cells)
	if err != nil {
		b.Fatalf("core.NewFrame() error = %v", err)
	}
	return frame
}

func benchmarkCell(b *testing.B, symbol rune) core.Cell {
	b.Helper()

	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		b.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		b.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}
	return core.NewCell(symbol, fg, bg)
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (discardWriter) WriteString(s string) (int, error) {
	return len(s), nil
}

type failOnceWriter struct {
	err    error
	failed bool
	out    bytes.Buffer
}

func (w *failOnceWriter) Write(p []byte) (int, error) {
	if !w.failed {
		w.failed = true
		return 0, w.err
	}
	return w.out.Write(p)
}
