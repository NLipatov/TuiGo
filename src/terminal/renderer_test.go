package terminal

import (
	"bytes"
	"errors"
	"testing"
	"tuigo/ansi"
	"tuigo/domain"
)

func TestRendererRenderWritesCell(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	frame, err := domain.NewFrame(1, 1, []domain.Cell{domain.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(frame, &out)

	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "\x1b[1;1H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "x" + string(ansi.RESET)
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

	frame, err := domain.NewFrame(1, 1, []domain.Cell{domain.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(frame, &out)
	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.Render(); err != nil {
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

	firstFrame, err := domain.NewFrame(2, 1, []domain.Cell{
		domain.NewCell('x', fg, bg),
		domain.NewCell('y', fg, bg),
	})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}
	nextFrame, err := domain.NewFrame(2, 1, []domain.Cell{
		domain.NewCell('x', fg, bg),
		domain.NewCell('z', fg, bg),
	})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(firstFrame, &out)
	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.NextFrame(nextFrame); err != nil {
		t.Fatalf("NextFrame() error = %v", err)
	}
	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "\x1b[1;2H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "z" + string(ansi.RESET)
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

	firstFrame, err := domain.NewFrame(1, 1, []domain.Cell{domain.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}
	nextFrame, err := domain.NewFrame(2, 1, []domain.Cell{
		domain.NewCell('y', fg, bg),
		domain.NewCell('z', fg, bg),
	})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}

	var out bytes.Buffer
	renderer := NewRenderer(firstFrame, &out)
	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	out.Reset()
	if err := renderer.NextFrame(nextFrame); err != nil {
		t.Fatalf("NextFrame() error = %v", err)
	}
	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "\x1b[1;1H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "y" + string(ansi.RESET) +
		"\x1b[1;2H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "z" + string(ansi.RESET)
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

	frame, err := domain.NewFrame(1, 1, []domain.Cell{domain.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}

	writeErr := errors.New("write failed")
	writer := failOnceWriter{err: writeErr}
	renderer := NewRenderer(frame, &writer)

	if err := renderer.Render(); !errors.Is(err, writeErr) {
		t.Fatalf("Render() error = %v, want %v", err, writeErr)
	}

	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	want := "\x1b[1;1H" + string(ansi.FG_RED) + string(ansi.BG_BLACK) + "x" + string(ansi.RESET)
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

	frame, err := domain.NewFrame(1, 1, []domain.Cell{domain.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}

	renderer := NewRenderer(frame, discardWriter{})
	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = renderer.Render()
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

	frame, err := domain.NewFrame(1, 1, []domain.Cell{domain.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}

	renderer := NewRenderer(frame, discardWriter{})
	allocs := testing.AllocsPerRun(1000, func() {
		renderer.fullRepaint = true
		_ = renderer.Render()
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

	firstFrame, err := domain.NewFrame(1, 1, []domain.Cell{domain.NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}
	nextFrame, err := domain.NewFrame(1, 1, []domain.Cell{domain.NewCell('y', fg, bg)})
	if err != nil {
		t.Fatalf("domain.NewFrame() error = %v", err)
	}

	renderer := NewRenderer(firstFrame, discardWriter{})
	if err := renderer.Render(); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		_ = renderer.NextFrame(nextFrame)
		_ = renderer.Render()
		_ = renderer.NextFrame(firstFrame)
		_ = renderer.Render()
	})
	if allocs != 0 {
		t.Fatalf("allocations per changed-cell render = %.2f, want 0", allocs)
	}
}

func BenchmarkRendererRenderUnchangedFrame(b *testing.B) {
	frame := benchmarkFrame(b, 80, 24, 'x')
	renderer := NewRenderer(frame, discardWriter{})
	if err := renderer.Render(); err != nil {
		b.Fatalf("Render() error = %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := renderer.Render(); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
}

func BenchmarkRendererRenderFullFrame(b *testing.B) {
	frame := benchmarkFrame(b, 80, 24, 'x')
	renderer := NewRenderer(frame, discardWriter{})

	b.ReportAllocs()
	for b.Loop() {
		renderer.fullRepaint = true
		if err := renderer.Render(); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
}

func BenchmarkRendererRenderChangedCell(b *testing.B) {
	firstFrame := benchmarkFrame(b, 80, 24, 'x')
	nextFrame := benchmarkFrameWithLastCell(b, 80, 24, 'x', 'y')

	renderer := NewRenderer(firstFrame, discardWriter{})
	if err := renderer.Render(); err != nil {
		b.Fatalf("Render() error = %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := renderer.NextFrame(nextFrame); err != nil {
			b.Fatalf("NextFrame() error = %v", err)
		}
		if err := renderer.Render(); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
		if err := renderer.NextFrame(firstFrame); err != nil {
			b.Fatalf("NextFrame() error = %v", err)
		}
		if err := renderer.Render(); err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
}

func benchmarkFrame(b *testing.B, width, height int, symbol rune) domain.Frame {
	b.Helper()

	cells := make([]domain.Cell, width*height)
	cell := benchmarkCell(b, symbol)
	for i := range cells {
		cells[i] = cell
	}

	frame, err := domain.NewFrame(width, height, cells)
	if err != nil {
		b.Fatalf("domain.NewFrame() error = %v", err)
	}
	return frame
}

func benchmarkFrameWithLastCell(b *testing.B, width, height int, symbol, lastSymbol rune) domain.Frame {
	b.Helper()

	cells := make([]domain.Cell, width*height)
	cell := benchmarkCell(b, symbol)
	for i := range cells {
		cells[i] = cell
	}
	cells[len(cells)-1] = benchmarkCell(b, lastSymbol)

	frame, err := domain.NewFrame(width, height, cells)
	if err != nil {
		b.Fatalf("domain.NewFrame() error = %v", err)
	}
	return frame
}

func benchmarkCell(b *testing.B, symbol rune) domain.Cell {
	b.Helper()

	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		b.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		b.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}
	return domain.NewCell(symbol, fg, bg)
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
