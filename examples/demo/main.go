package main

import (
	"context"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/NLipatov/tuigo/ansi"
	"github.com/NLipatov/tuigo/core"
	"github.com/NLipatov/tuigo/terminal"
	"github.com/NLipatov/tuigo/terminal/input"
)

const (
	targetFPS           = 120
	targetFrameDuration = time.Second / targetFPS
	minDemoWidth        = 72
	minDemoHeight       = 24
	maxDemoWidth        = 120
	maxDemoHeight       = 24
	maxDrawSamples      = 72
)

var keyCodeLabels = [...]string{
	input.KeyUnknown:   "unknown",
	input.KeyRune:      "rune",
	input.KeyEnter:     "enter",
	input.KeyEsc:       "esc",
	input.KeyTab:       "tab",
	input.KeyBackspace: "backspace",
	input.KeyDelete:    "delete",
	input.KeyInsert:    "insert",
	input.KeyUp:        "up",
	input.KeyDown:      "down",
	input.KeyLeft:      "left",
	input.KeyRight:     "right",
	input.KeyHome:      "home",
	input.KeyEnd:       "end",
	input.KeyPageUp:    "page-up",
	input.KeyPageDown:  "page-down",
	input.KeyF1:        "f1",
	input.KeyF2:        "f2",
	input.KeyF3:        "f3",
	input.KeyF4:        "f4",
	input.KeyF5:        "f5",
	input.KeyF6:        "f6",
	input.KeyF7:        "f7",
	input.KeyF8:        "f8",
	input.KeyF9:        "f9",
	input.KeyF10:       "f10",
	input.KeyF11:       "f11",
	input.KeyF12:       "f12",
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
	defer func() {
		if err := session.Close(); err != nil {
			panic(err)
		}
	}()

	width, height, err := session.Size()
	if err != nil {
		panic(err)
	}

	if err := runDemo(ctx, cancel, &session, events, demoState{
		width: width, height: height, log: []string{"0000  ready"},
	}); err != nil {
		panic(err)
	}
}

func runDemo(
	ctx context.Context,
	cancel context.CancelFunc,
	session *terminal.Session,
	events <-chan terminal.Event,
	state demoState,
) error {
	var buffers frameBuffers
	if err := renderDemo(session, &state, &buffers); err != nil {
		return err
	}

	ticker := time.NewTicker(targetFrameDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			state.frame++
			if err := renderDemo(session, &state, &buffers); err != nil {
				return err
			}
		case event, ok := <-events:
			if !ok {
				return nil
			}
			if handleEvent(&state, cancel, event) {
				return nil
			}
			if err := renderDemo(session, &state, &buffers); err != nil {
				return err
			}
		}
	}
}

func handleEvent(state *demoState, cancel context.CancelFunc, event terminal.Event) bool {
	state.events++
	label := ""
	switch event.Type {
	case terminal.EventKey:
		key := keyLabel(event.Key)
		label = "key " + key
		if quitRequested(event.Key) {
			cancel()
			return true
		}
		if event.Key.Code == input.KeyRune && event.Key.Text == "r" {
			state.frame = 0
			state.events = 0
			state.log = []string{"0000  reset"}
			state.draw = 0
			state.draws = nil
			return false
		}
	case terminal.EventResize:
		state.width = event.Resize.Width
		state.height = event.Resize.Height
		label = "resize " + sizeLabel(state.width, state.height)
	case terminal.EventMouse:
		label = mouseLabel(event.Mouse)
	case terminal.EventError:
		panic(event.Err)
	}
	if label != "" {
		appendEvent(state, label)
	}
	return false
}

type demoState struct {
	width  int
	height int
	frame  int
	events int
	log    []string
	draw   time.Duration
	draws  []time.Duration
}

type frameBuffers struct {
	buffers [2]frameBuffer
	next    int
	width   int
	height  int
}

type frameBuffer struct {
	frame core.Frame
	cells []core.Cell
}

func (b *frameBuffers) ensure(width, height int) error {
	if b.width == width && b.height == height && b.buffers[0].cells != nil {
		return nil
	}
	if width <= 0 || height <= 0 {
		return core.ErrInvalidFrameDimensions
	}

	for idx := range b.buffers {
		cells := make([]core.Cell, width*height)
		frame, err := core.NewFrame(width, height, cells)
		if err != nil {
			return err
		}
		b.buffers[idx] = frameBuffer{
			frame: frame,
			cells: cells,
		}
	}
	b.next = 0
	b.width = width
	b.height = height
	return nil
}

func (b *frameBuffers) back() *frameBuffer {
	return &b.buffers[b.next]
}

func (b *frameBuffers) swap() {
	b.next = 1 - b.next
}

func renderDemo(session *terminal.Session, state *demoState, buffers *frameBuffers) error {
	if err := buffers.ensure(state.width, state.height); err != nil {
		return err
	}
	buffer := buffers.back()
	drawDemoFrame(buffer.cells, *state)

	start := time.Now()
	if err := session.Render(buffer.frame); err != nil {
		return err
	}
	buffers.swap()
	appendDrawSample(state, time.Since(start))
	return nil
}

func drawDemoFrame(cells []core.Cell, state demoState) {
	colors := newDemoPalette()
	for i := range cells {
		cells[i] = mustCell(" ", colors.fg, colors.bg)
	}

	if state.width < minDemoWidth || state.height < minDemoHeight {
		drawCompact(cells, state, colors)
		return
	}

	left, top, boxWidth, boxHeight := demoBounds(state.width, state.height)
	drawBox(cells, state.width, state.height, left, top, boxWidth, boxHeight, colors.fg, colors.bg)
	drawHeader(cells, state, left, top, boxWidth, colors)
	drawPanels(cells, state, left, top+5, colors)
	drawEventStream(cells, state, left+2, top+12, boxWidth-4, colors)
	drawSeparator(cells, state.width, state.height, left, top+boxHeight-3, boxWidth, colors.fg, colors.bg)
	drawTextClipped(cells, state.width, state.height, left+2, top+boxHeight-2, boxWidth-4, "q/esc quit   r reset", colors.fg, colors.bg)
}

func drawCompact(cells []core.Cell, state demoState, colors demoPalette) {
	drawText(cells, state.width, state.height, 1, 1, "tuigo", colors.accent, colors.bg)
	drawText(cells, state.width, state.height, 1, 3, "resize terminal to at least 72x24", colors.fg, colors.bg)
	drawText(cells, state.width, state.height, 1, 5, "q / esc / ctrl+c quit", colors.fg, colors.bg)
}

func drawHeader(cells []core.Cell, state demoState, left, top, width int, colors demoPalette) {
	title := "tuigo"
	mode := sizeLabel(state.width, state.height)
	drawText(cells, state.width, state.height, left+2, top+1, title, colors.accent, colors.bg)
	drawText(cells, state.width, state.height, left+width-len(mode)-2, top+1, mode, colors.fg, colors.bg)
	drawSeparator(cells, state.width, state.height, left, top+2, width, colors.fg, colors.bg)
}

func drawPanels(cells []core.Cell, state demoState, left, top int, colors demoPalette) {
	median := medianDrawDuration(state.draws)
	drawMetrics(cells, state.width, state.height, left+2, top, 20, "frame", []metric{
		{"size", sizeLabel(state.width, state.height)},
		{"cells", intLabel(state.width * state.height)},
		{"buffers", "2"},
	}, colors)
	drawMetrics(cells, state.width, state.height, left+26, top, 20, "render", []metric{
		{"frame", intLabel(state.frame)},
		{"last", formatDrawDuration(state.draw)},
		{"median", formatDrawDuration(median)},
	}, colors)
}

type metric struct {
	label string
	value string
}

func drawMetrics(cells []core.Cell, width, height, left, top, columnWidth int, title string, rows []metric, colors demoPalette) {
	drawTextClipped(cells, width, height, left, top, columnWidth, title, colors.accent, colors.bg)
	valueLeft := left + 10
	valueWidth := columnWidth - 10
	for idx, row := range rows {
		y := top + 2 + idx
		drawTextClipped(cells, width, height, left, y, 8, row.label, colors.fg, colors.bg)
		drawTextClipped(cells, width, height, valueLeft, y, valueWidth, row.value, colors.fg, colors.bg)
	}
}

func drawEventStream(cells []core.Cell, state demoState, left, top, streamWidth int, colors demoPalette) {
	drawText(cells, state.width, state.height, left, top, "event stream", colors.accent, colors.bg)
	drawTextClipped(cells, state.width, state.height, left, top+2, streamWidth, "events "+intLabel(state.events), colors.fg, colors.bg)
	start := max(0, len(state.log)-4)
	for idx, line := range state.log[start:] {
		drawTextClipped(cells, state.width, state.height, left, top+3+idx, streamWidth, line, colors.fg, colors.bg)
	}
}

func appendEvent(state *demoState, label string) {
	state.log = append(state.log, zeroPadInt(state.events, 4)+"  "+label)
	if len(state.log) > 5 {
		state.log = state.log[len(state.log)-5:]
	}
}

func appendDrawSample(state *demoState, draw time.Duration) {
	state.draw = draw
	state.draws = append(state.draws, draw)
	if len(state.draws) > maxDrawSamples {
		state.draws = state.draws[len(state.draws)-maxDrawSamples:]
	}
}

func formatDrawDuration(draw time.Duration) string {
	if draw <= 0 {
		return "--"
	}
	if draw < time.Millisecond {
		return int64Label(draw.Microseconds()) + "µs"
	}
	if draw >= time.Second {
		return fixedDuration(draw, time.Second) + "s"
	}
	return fixedDuration(draw, time.Millisecond) + "ms"
}

func fixedDuration(draw, unit time.Duration) string {
	hundredths := int64((draw + unit/200) / (unit / 100))
	whole := hundredths / 100
	fraction := hundredths % 100
	return int64Label(whole) + "." + zeroPadInt64(fraction, 2)
}

func sizeLabel(width, height int) string {
	return intLabel(width) + "x" + intLabel(height)
}

func intLabel(value int) string {
	return strconv.Itoa(value)
}

func int64Label(value int64) string {
	return strconv.FormatInt(value, 10)
}

func zeroPadInt(value, width int) string {
	return zeroPadInt64(int64(value), width)
}

func zeroPadInt64(value int64, width int) string {
	text := int64Label(value)
	if len(text) >= width {
		return text
	}
	return strings.Repeat("0", width-len(text)) + text
}

func medianDrawDuration(samples []time.Duration) time.Duration {
	if len(samples) == 0 {
		return 0
	}
	sorted := slices.Clone(samples)
	slices.Sort(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}

type demoPalette struct {
	bg     ansi.Color
	fg     ansi.Color
	accent ansi.Color
}

func newDemoPalette() demoPalette {
	return demoPalette{
		bg:     mustColor(ansi.BG_BLACK),
		fg:     mustColor(ansi.FG_GREEN),
		accent: mustColor(ansi.FG_BOLD_GREEN),
	}
}

func demoBounds(width, height int) (int, int, int, int) {
	boxWidth := min(width, maxDemoWidth)
	boxHeight := min(height, maxDemoHeight)
	left := max(0, (width-boxWidth)/2)
	top := max(0, (height-boxHeight)/2)
	return left, top, boxWidth, boxHeight
}

func drawBox(cells []core.Cell, width, height, left, top, boxWidth, boxHeight int, fg, bg ansi.Color) {
	if boxWidth < 2 || boxHeight < 2 {
		return
	}
	right := left + boxWidth - 1
	bottom := top + boxHeight - 1
	putCell(cells, width, height, left, top, mustCell("┌", fg, bg))
	putCell(cells, width, height, right, top, mustCell("┐", fg, bg))
	putCell(cells, width, height, left, bottom, mustCell("└", fg, bg))
	putCell(cells, width, height, right, bottom, mustCell("┘", fg, bg))
	for x := left + 1; x < right; x++ {
		putCell(cells, width, height, x, top, mustCell("─", fg, bg))
		putCell(cells, width, height, x, bottom, mustCell("─", fg, bg))
	}
	for y := top + 1; y < bottom; y++ {
		putCell(cells, width, height, left, y, mustCell("│", fg, bg))
		putCell(cells, width, height, right, y, mustCell("│", fg, bg))
	}
}

func drawSeparator(cells []core.Cell, width, height, left, y, lineWidth int, fg, bg ansi.Color) {
	putCell(cells, width, height, left, y, mustCell("├", fg, bg))
	putCell(cells, width, height, left+lineWidth-1, y, mustCell("┤", fg, bg))
	for x := 1; x < lineWidth-1; x++ {
		putCell(cells, width, height, left+x, y, mustCell("─", fg, bg))
	}
}

func drawText(cells []core.Cell, width, height, left, y int, text string, fg, bg ansi.Color) {
	for x, char := range []rune(text) {
		putCell(cells, width, height, left+x, y, mustCell(string(char), fg, bg))
	}
}

func drawTextClipped(cells []core.Cell, width, height, left, y, maxWidth int, text string, fg, bg ansi.Color) {
	drawText(cells, width, height, left, y, trimLabel(text, maxWidth), fg, bg)
}

func putCell(cells []core.Cell, width, height, x, y int, cell core.Cell) {
	if x < 0 || y < 0 || x >= width || y >= height {
		return
	}
	cells[y*width+x] = cell
}

func keyLabel(event input.KeyEvent) string {
	var text string
	if event.Code == input.KeyRune {
		text = keyTextLabel(event.Text)
	} else {
		text = keyCodeLabel(event.Code)
	}
	if event.Mod == input.ModNone {
		return text
	}
	return modLabel(event.Mod) + "+" + text
}

func keyCodeLabel(code input.KeyCode) string {
	idx := int(code)
	if idx >= 0 && idx < len(keyCodeLabels) && keyCodeLabels[idx] != "" {
		return keyCodeLabels[idx]
	}
	return "key-" + intLabel(idx)
}

func keyTextLabel(text string) string {
	if text == " " {
		return "space"
	}
	return text
}

func mouseLabel(event input.MouseEvent) string {
	label := "mouse " + mouseButtonLabel(event.Button) + " " +
		mouseActionLabel(event.Action) + " " +
		intLabel(event.X) + "," + intLabel(event.Y)
	if event.Mod == input.ModNone {
		return label
	}
	return modLabel(event.Mod) + " " + label
}

func mouseButtonLabel(button input.MouseButton) string {
	switch button {
	case input.MouseButtonLeft:
		return "left"
	case input.MouseButtonMiddle:
		return "middle"
	case input.MouseButtonRight:
		return "right"
	case input.MouseButtonWheelUp:
		return "wheel-up"
	case input.MouseButtonWheelDown:
		return "wheel-down"
	default:
		return "unknown"
	}
}

func mouseActionLabel(action input.MouseAction) string {
	switch action {
	case input.MouseActionPress:
		return "press"
	case input.MouseActionRelease:
		return "release"
	case input.MouseActionDrag:
		return "drag"
	case input.MouseActionWheel:
		return "wheel"
	default:
		return "unknown"
	}
}

func modLabel(mod input.KeyMod) string {
	parts := make([]string, 0, 3)
	if mod&input.ModCtrl != 0 {
		parts = append(parts, "ctrl")
	}
	if mod&input.ModAlt != 0 {
		parts = append(parts, "alt")
	}
	if mod&input.ModShift != 0 {
		parts = append(parts, "shift")
	}
	return strings.Join(parts, "+")
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

func trimLabel(text string, limit int) string {
	if limit <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	if limit == 1 {
		return "…"
	}
	return string(runes[:limit-1]) + "…"
}

func mustColor(sequence ansi.ANSIEscapeSequence) ansi.Color {
	color, err := ansi.NewColor(sequence)
	if err != nil {
		panic(err)
	}
	return color
}

func mustCell(text string, fg, bg ansi.Color) core.Cell {
	cell, err := core.NewCell(text, fg, bg)
	if err != nil {
		panic(err)
	}
	return cell
}
