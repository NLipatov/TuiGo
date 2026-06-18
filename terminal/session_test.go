package terminal

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/NLipatov/tuigo/color"
	"github.com/NLipatov/tuigo/core"
	"github.com/NLipatov/tuigo/internal/ansi"
	"github.com/NLipatov/tuigo/keyboard"
	"github.com/NLipatov/tuigo/mouse"
	"github.com/NLipatov/tuigo/terminal/internal/input"
	"github.com/NLipatov/tuigo/terminal/resize"
	"golang.org/x/term"
)

func TestNewSessionWiresRendererToSessionOutput(t *testing.T) {
	in, closeIn, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	defer func() {
		_ = in.Close()
	}()
	defer func() {
		_ = closeIn.Close()
	}()

	var out bytes.Buffer
	session, err := NewSession(context.Background(), in, &out)
	if err != nil {
		t.Fatalf("NewSession() error = %v", err)
	}

	cell, err := core.NewCell("x", color.FgRed, color.BgBlack)
	if err != nil {
		t.Fatalf("core.NewCell(%q) error = %v", "x", err)
	}
	frame, err := core.NewFrame(1, 1, []core.Cell{cell})
	if err != nil {
		t.Fatalf("core.NewFrame() error = %v", err)
	}

	if err := session.Render(frame); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !bytes.Contains(out.Bytes(), []byte("x")) {
		t.Fatalf("session output = %q, want rendered frame", out.String())
	}
}

func TestSessionCloseBeforeStartWritesNothing(t *testing.T) {
	var out bytes.Buffer
	session := Session{writer: &out}

	if err := session.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if got := out.String(); got != "" {
		t.Fatalf("session output = %q, want empty output", got)
	}
}

func TestSessionStartWritesNothingWhenRawModeSetupFails(t *testing.T) {
	var out bytes.Buffer
	session := Session{
		ctx:    context.Background(),
		writer: &out,
		device: Device{fd: -1},
	}

	_, err := session.Start()
	if err == nil {
		t.Fatal("Start() error = nil, want error")
	}
	if got := out.String(); got != "" {
		t.Fatalf("session output = %q, want empty output", got)
	}
}

func TestSessionStartRestoresTerminalOnSetupError(t *testing.T) {
	setupErr := errors.New("setup failed")
	writer := failOnceSessionWriter{err: setupErr}
	session := Session{
		ctx:    context.Background(),
		writer: &writer,
		device: Device{
			fd:           -1,
			initialState: &term.State{},
		},
	}

	_, err := session.Start()
	if !errors.Is(err, setupErr) {
		t.Fatalf("Start() error = %v, want wrapped %v", err, setupErr)
	}

	want := string(ansi.DISABLE_SGR_MOUSE) +
		string(ansi.DISABLE_MOUSE_DRAG) +
		string(ansi.DISABLE_MOUSE_REPORTING) +
		string(ansi.RESET) +
		string(ansi.SHOW_CURSOR) +
		string(ansi.EXIT_ALTERNATE_SCREEN)
	if got := writer.out.String(); got != want {
		t.Fatalf("session output = %q, want restore commands %q", got, want)
	}
}

func TestSessionCloseRestoresTerminalBestEffort(t *testing.T) {
	restoreErr := errors.New("restore failed")
	writer := failOnceSessionWriter{err: restoreErr}
	session := Session{
		writer: &writer,
		device: Device{
			fd:           -1,
			initialState: &term.State{},
		},
	}

	err := session.Close()
	if !errors.Is(err, restoreErr) {
		t.Fatalf("Close() error = %v, want wrapped %v", err, restoreErr)
	}

	want := string(ansi.DISABLE_MOUSE_DRAG) +
		string(ansi.DISABLE_MOUSE_REPORTING) +
		string(ansi.RESET) +
		string(ansi.SHOW_CURSOR) +
		string(ansi.EXIT_ALTERNATE_SCREEN)
	if got := writer.out.String(); got != want {
		t.Fatalf("session output = %q, want remaining restore commands %q", got, want)
	}
}

func TestSessionEventLoopForwardsInputAndResizeEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := Session{ctx: ctx}
	resizeCh := make(chan resize.Event, 1)
	inputCh := make(chan input.Event, 1)
	listener := contextCanceledListener(ctx)
	events := session.runEventLoop(resizeCh, listener, inputCh, listener)

	resizeEvent := resize.Event{Width: 100, Height: 40}
	resizeCh <- resizeEvent
	got := receiveSessionEvent(t, events)
	if got.Type != EventResize {
		t.Fatalf("event type = %v, want %v", got.Type, EventResize)
	}
	if got.Resize != resizeEvent {
		t.Fatalf("resize event = %#v, want %#v", got.Resize, resizeEvent)
	}

	keyEvent := keyboard.KeyEvent{Code: keyboard.KeyCode(1), Text: "a", Mod: keyboard.ModCtrl}
	inputCh <- input.Event{Type: input.EventTypeKey, Key: keyEvent}
	got = receiveSessionEvent(t, events)
	if got.Type != EventKey {
		t.Fatalf("event type = %v, want %v", got.Type, EventKey)
	}
	if got.Key != keyEvent {
		t.Fatalf("key event = %#v, want %#v", got.Key, keyEvent)
	}

	mouseEvent := mouse.MouseEvent{
		X:      10,
		Y:      20,
		Button: mouse.MouseButtonLeft,
		Action: mouse.MouseActionPress,
		Mod:    keyboard.ModShift,
	}
	inputCh <- input.Event{Type: input.EventTypeMouse, Mouse: mouseEvent}
	got = receiveSessionEvent(t, events)
	if got.Type != EventMouse {
		t.Fatalf("event type = %v, want %v", got.Type, EventMouse)
	}
	if got.Mouse != mouseEvent {
		t.Fatalf("mouse event = %#v, want %#v", got.Mouse, mouseEvent)
	}
}

func TestSessionEventLoopEmitsListenerErrors(t *testing.T) {
	resizeErr := errors.New("resize failed")
	keyErr := errors.New("key failed")
	tests := []struct {
		name      string
		resizeErr error
		keyErr    error
		want      error
	}{
		{
			name:      "resize listener",
			resizeErr: resizeErr,
			want:      resizeErr,
		},
		{
			name:   "key listener",
			keyErr: keyErr,
			want:   keyErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			session := Session{ctx: ctx}
			resizeListener := contextCanceledListener(ctx)
			if tt.resizeErr != nil {
				resizeListener = listenerFunc(func() error { return tt.resizeErr })
			}
			keyListener := contextCanceledListener(ctx)
			if tt.keyErr != nil {
				keyListener = listenerFunc(func() error { return tt.keyErr })
			}
			events := session.runEventLoop(
				make(chan resize.Event),
				resizeListener,
				make(chan input.Event),
				keyListener,
			)

			got := receiveSessionEvent(t, events)
			if got.Type != EventError {
				t.Fatalf("event type = %v, want %v", got.Type, EventError)
			}
			if !errors.Is(got.Err, tt.want) {
				t.Fatalf("event error = %v, want %v", got.Err, tt.want)
			}
		})
	}
}

func TestSessionEventLoopDoesNotEmitErrorOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	session := Session{ctx: ctx}
	listener := contextCanceledListener(ctx)
	events := session.runEventLoop(
		make(chan resize.Event),
		listener,
		make(chan input.Event),
		listener,
	)

	cancel()

	select {
	case event, ok := <-events:
		if ok {
			t.Fatalf("unexpected event after context cancel: %#v", event)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event channel to close")
	}
}

func receiveSessionEvent(t *testing.T, events <-chan Event) Event {
	t.Helper()
	select {
	case event, ok := <-events:
		if !ok {
			t.Fatal("event channel closed")
		}
		return event
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
		return Event{}
	}
}

type listenerFunc func() error

func (f listenerFunc) Listen() error {
	return f()
}

func contextCanceledListener(ctx context.Context) eventListener {
	return listenerFunc(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
}

type failOnceSessionWriter struct {
	err    error
	failed bool
	out    bytes.Buffer
}

func (w *failOnceSessionWriter) Write(p []byte) (int, error) {
	if !w.failed {
		w.failed = true
		return 0, w.err
	}
	return w.out.Write(p)
}

func (w *failOnceSessionWriter) WriteString(s string) (int, error) {
	if !w.failed {
		w.failed = true
		return 0, w.err
	}
	return w.out.WriteString(s)
}
