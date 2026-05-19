package tuigo

import (
	"bytes"
	"testing"
	"tuigo/ansi"
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

	frame, err := NewFrame(1, 1, []Cell{NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
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

	frame, err := NewFrame(1, 1, []Cell{NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
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

	firstFrame, err := NewFrame(2, 1, []Cell{
		NewCell('x', fg, bg),
		NewCell('y', fg, bg),
	})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
	}
	nextFrame, err := NewFrame(2, 1, []Cell{
		NewCell('x', fg, bg),
		NewCell('z', fg, bg),
	})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
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

	firstFrame, err := NewFrame(1, 1, []Cell{NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
	}
	nextFrame, err := NewFrame(2, 1, []Cell{
		NewCell('y', fg, bg),
		NewCell('z', fg, bg),
	})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
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

func TestRendererRenderDoesNotAllocateWhenFrameIsUnchanged(t *testing.T) {
	fg, err := ansi.NewColor(ansi.FG_RED)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.FG_RED, err)
	}
	bg, err := ansi.NewColor(ansi.BG_BLACK)
	if err != nil {
		t.Fatalf("NewColor(%q) error = %v", ansi.BG_BLACK, err)
	}

	frame, err := NewFrame(1, 1, []Cell{NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
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

	frame, err := NewFrame(1, 1, []Cell{NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
	}

	renderer := NewRenderer(frame, discardWriter{})
	allocs := testing.AllocsPerRun(1000, func() {
		renderer.firstRender = true
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

	firstFrame, err := NewFrame(1, 1, []Cell{NewCell('x', fg, bg)})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
	}
	nextFrame, err := NewFrame(1, 1, []Cell{NewCell('y', fg, bg)})
	if err != nil {
		t.Fatalf("NewFrame() error = %v", err)
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

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (discardWriter) WriteString(s string) (int, error) {
	return len(s), nil
}
