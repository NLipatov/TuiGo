package input

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewInputListenerRejectsNilDependencies(t *testing.T) {
	ctx := context.Background()
	reader := &scriptedReadCloser{}
	inputParser := NewParser()
	events := make(chan KeyEvent)

	tests := []struct {
		name   string
		ctx    context.Context
		reader io.ReadCloser
		parser EventParser
		events chan<- KeyEvent
		want   error
	}{
		{
			name:   "nil context",
			ctx:    nil,
			reader: reader,
			parser: inputParser,
			events: events,
			want:   ErrNilContext,
		},
		{
			name:   "nil reader",
			ctx:    ctx,
			reader: nil,
			parser: inputParser,
			events: events,
			want:   ErrNilReader,
		},
		{
			name:   "nil parser",
			ctx:    ctx,
			reader: reader,
			parser: nil,
			events: events,
			want:   ErrNilParser,
		},
		{
			name:   "nil event channel",
			ctx:    ctx,
			reader: reader,
			parser: inputParser,
			events: nil,
			want:   ErrNilEventChan,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewListener(tt.ctx, tt.reader, tt.parser, tt.events)
			if !errors.Is(err, tt.want) {
				t.Fatalf("NewInputListener() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestInputListenerListenFeedsParserAndEmitsEvents(t *testing.T) {
	readErr := errors.New("read failed")
	reader := &scriptedReadCloser{
		reads: [][]byte{
			[]byte("ab"),
			[]byte("c"),
		},
		err: readErr,
	}
	out := make(chan KeyEvent, 3)
	input, err := NewListener(context.Background(), reader, NewParser(), out)
	if err != nil {
		t.Fatalf("NewInputListener() error = %v", err)
	}

	err = input.Listen()
	if !errors.Is(err, readErr) {
		t.Fatalf("Listen() error = %v, want %v", err, readErr)
	}

	wantEvents := []KeyEvent{
		{Code: KeyRune, Text: "a", Mod: ModNone},
		{Code: KeyRune, Text: "b", Mod: ModNone},
		{Code: KeyRune, Text: "c", Mod: ModNone},
	}
	for idx, want := range wantEvents {
		got := <-out
		if got != want {
			t.Fatalf("event %d = %#v, want %#v", idx, got, want)
		}
	}
	select {
	case got := <-out:
		t.Fatalf("unexpected extra event: %#v", got)
	default:
	}
}

func TestInputListenerListenFlushesPendingEscapeAfterTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := newBlockingAfterReadsReadCloser([][]byte{[]byte("\x1b")})
	out := make(chan KeyEvent, 1)
	input, err := NewListener(ctx, reader, NewParser(), out)
	if err != nil {
		t.Fatalf("NewInputListener() error = %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- input.Listen()
	}()

	got := receiveInputEvent(t, out)
	if want := (KeyEvent{Code: KeyEsc}); got != want {
		t.Fatalf("event = %#v, want %#v", got, want)
	}

	cancel()
	select {
	case err := <-errCh:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Listen() error = %v, want %v", err, context.Canceled)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Listen to return")
	}
}

func TestInputListenerListenReturnsReadErrorWithoutCallingParser(t *testing.T) {
	readErr := errors.New("read failed")
	reader := &scriptedReadCloser{err: readErr}
	out := make(chan KeyEvent, 1)
	input, err := NewListener(context.Background(), reader, NewParser(), out)
	if err != nil {
		t.Fatalf("NewInputListener() error = %v", err)
	}

	err = input.Listen()
	if !errors.Is(err, readErr) {
		t.Fatalf("Listen() error = %v, want %v", err, readErr)
	}
	select {
	case got := <-out:
		t.Fatalf("unexpected event: %#v", got)
	default:
	}
}

func TestInputListenerListenReturnsContextErrorAndClosesReader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := newBlockingReadCloser()
	input, err := NewListener(ctx, reader, NewParser(), make(chan KeyEvent))
	if err != nil {
		t.Fatalf("NewInputListener() error = %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- input.Listen()
	}()

	select {
	case <-reader.readStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for reader to block")
	}

	cancel()

	select {
	case err := <-errCh:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Listen() error = %v, want %v", err, context.Canceled)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Listen to return")
	}

	if got := reader.closeCount.Load(); got != 1 {
		t.Fatalf("reader Close calls = %d, want 1", got)
	}
}

func TestInputListenerCloseIsIdempotent(t *testing.T) {
	closeErr := errors.New("close failed")
	reader := &scriptedReadCloser{closeErr: closeErr}
	input, err := NewListener(
		context.Background(),
		reader,
		NewParser(),
		make(chan KeyEvent),
	)
	if err != nil {
		t.Fatalf("NewInputListener() error = %v", err)
	}

	if err := input.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("first Close() error = %v, want %v", err, closeErr)
	}
	if err := input.Close(); !errors.Is(err, closeErr) {
		t.Fatalf("second Close() error = %v, want %v", err, closeErr)
	}
	if reader.closeCount != 1 {
		t.Fatalf("reader Close calls = %d, want 1", reader.closeCount)
	}
}

type scriptedReadCloser struct {
	reads      [][]byte
	err        error
	closeErr   error
	closeCount int
}

func (r *scriptedReadCloser) Read(p []byte) (int, error) {
	if len(r.reads) == 0 {
		return 0, r.err
	}
	n := copy(p, r.reads[0])
	r.reads = r.reads[1:]
	return n, nil
}

func (r *scriptedReadCloser) Close() error {
	r.closeCount++
	return r.closeErr
}

type blockingReadCloser struct {
	readStarted chan struct{}
	closed      chan struct{}
	closeOnce   sync.Once
	startOnce   sync.Once
	closeCount  atomic.Int32
}

func newBlockingReadCloser() *blockingReadCloser {
	return &blockingReadCloser{
		readStarted: make(chan struct{}),
		closed:      make(chan struct{}),
	}
}

func (r *blockingReadCloser) Read([]byte) (int, error) {
	r.startOnce.Do(func() {
		close(r.readStarted)
	})
	<-r.closed
	return 0, errors.New("reader closed")
}

func (r *blockingReadCloser) Close() error {
	r.closeOnce.Do(func() {
		r.closeCount.Add(1)
		close(r.closed)
	})
	return nil
}

type blockingAfterReadsReadCloser struct {
	reads     [][]byte
	closed    chan struct{}
	closeOnce sync.Once
}

func newBlockingAfterReadsReadCloser(reads [][]byte) *blockingAfterReadsReadCloser {
	return &blockingAfterReadsReadCloser{
		reads:  reads,
		closed: make(chan struct{}),
	}
}

func (r *blockingAfterReadsReadCloser) Read(p []byte) (int, error) {
	if len(r.reads) > 0 {
		n := copy(p, r.reads[0])
		r.reads = r.reads[1:]
		return n, nil
	}
	<-r.closed
	return 0, errors.New("reader closed")
}

func (r *blockingAfterReadsReadCloser) Close() error {
	r.closeOnce.Do(func() {
		close(r.closed)
	})
	return nil
}

func receiveInputEvent(t *testing.T, events <-chan KeyEvent) KeyEvent {
	t.Helper()
	select {
	case event := <-events:
		return event
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for input event")
		return KeyEvent{}
	}
}
