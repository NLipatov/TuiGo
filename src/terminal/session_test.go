package terminal

import (
	"context"
	"errors"
	"testing"
	"time"
	"tuigo/terminal/input"
	"tuigo/terminal/resize"
)

func TestSessionEventLoopForwardsInputAndResizeEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := Session{ctx: ctx}
	resizeCh := make(chan resize.Event, 1)
	keyCh := make(chan input.Event, 1)
	listener := contextCanceledListener(ctx)
	events := session.runEventLoop(resizeCh, listener, keyCh, listener)

	resizeEvent := resize.Event{Width: 100, Height: 40}
	resizeCh <- resizeEvent
	got := receiveSessionEvent(t, events)
	if got.Type != EventResize {
		t.Fatalf("event type = %v, want %v", got.Type, EventResize)
	}
	if got.Resize != resizeEvent {
		t.Fatalf("resize event = %#v, want %#v", got.Resize, resizeEvent)
	}

	keyEvent := input.Event{Code: input.KeyCode(1), Text: "a", Mod: input.ModCtrl}
	keyCh <- keyEvent
	got = receiveSessionEvent(t, events)
	if got.Type != EventKey {
		t.Fatalf("event type = %v, want %v", got.Type, EventKey)
	}
	if got.Key != keyEvent {
		t.Fatalf("key event = %#v, want %#v", got.Key, keyEvent)
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
