package terminal

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewInputRejectsNilDependencies(t *testing.T) {
	ctx := context.Background()
	reader := &scriptedReadCloser{}
	parser := &recordingParser{}
	events := make(chan Event)

	tests := []struct {
		name   string
		ctx    context.Context
		reader io.ReadCloser
		parser Parser
		events chan<- Event
		want   error
	}{
		{
			name:   "nil context",
			ctx:    nil,
			reader: reader,
			parser: parser,
			events: events,
			want:   ErrNilContext,
		},
		{
			name:   "nil reader",
			ctx:    ctx,
			reader: nil,
			parser: parser,
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
			parser: parser,
			events: nil,
			want:   ErrNilEventChan,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewInput(tt.ctx, tt.reader, tt.parser, tt.events)
			if !errors.Is(err, tt.want) {
				t.Fatalf("NewInput() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestNewInputAcceptsValidDependencies(t *testing.T) {
	_, err := NewInput(
		context.Background(),
		&scriptedReadCloser{},
		&recordingParser{},
		make(chan Event),
	)
	if err != nil {
		t.Fatalf("NewInput() error = %v", err)
	}
}

func TestInputListenFeedsParserAndEmitsEvents(t *testing.T) {
	readErr := errors.New("read failed")
	reader := &scriptedReadCloser{
		reads: [][]byte{
			[]byte("ab"),
			[]byte("c"),
		},
		err: readErr,
	}
	parser := &recordingParser{
		events: func(buf []byte) []Event {
			events := make([]Event, 0, len(buf))
			for _, b := range buf {
				events = append(events, Event{Code: KeyCode(b)})
			}
			return events
		},
	}
	out := make(chan Event, 3)
	input, err := NewInput(context.Background(), reader, parser, out)
	if err != nil {
		t.Fatalf("NewInput() error = %v", err)
	}

	err = input.Listen()
	if !errors.Is(err, readErr) {
		t.Fatalf("Listen() error = %v, want %v", err, readErr)
	}

	wantFeeds := [][]byte{
		[]byte("ab"),
		[]byte("c"),
	}
	if len(parser.feeds) != len(wantFeeds) {
		t.Fatalf("parser feeds = %q, want %q", parser.feeds, wantFeeds)
	}
	for idx, want := range wantFeeds {
		if got := string(parser.feeds[idx]); got != string(want) {
			t.Fatalf("parser feed %d = %q, want %q", idx, got, want)
		}
	}

	wantEvents := []Event{
		{Code: KeyCode('a')},
		{Code: KeyCode('b')},
		{Code: KeyCode('c')},
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

func TestInputListenReturnsReadErrorWithoutCallingParser(t *testing.T) {
	readErr := errors.New("read failed")
	reader := &scriptedReadCloser{err: readErr}
	parser := &recordingParser{}
	out := make(chan Event, 1)
	input, err := NewInput(context.Background(), reader, parser, out)
	if err != nil {
		t.Fatalf("NewInput() error = %v", err)
	}

	err = input.Listen()
	if !errors.Is(err, readErr) {
		t.Fatalf("Listen() error = %v, want %v", err, readErr)
	}
	if len(parser.feeds) != 0 {
		t.Fatalf("parser feeds = %q, want none", parser.feeds)
	}
	select {
	case got := <-out:
		t.Fatalf("unexpected event: %#v", got)
	default:
	}
}

func TestInputListenReturnsContextErrorAndClosesReader(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := newBlockingReadCloser()
	input, err := NewInput(ctx, reader, &recordingParser{}, make(chan Event))
	if err != nil {
		t.Fatalf("NewInput() error = %v", err)
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

func TestInputCloseIsIdempotent(t *testing.T) {
	closeErr := errors.New("close failed")
	reader := &scriptedReadCloser{closeErr: closeErr}
	input, err := NewInput(
		context.Background(),
		reader,
		&recordingParser{},
		make(chan Event),
	)
	if err != nil {
		t.Fatalf("NewInput() error = %v", err)
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

type recordingParser struct {
	feeds  [][]byte
	events func([]byte) []Event
}

func (p *recordingParser) Feed(buf []byte) []Event {
	p.feeds = append(p.feeds, append([]byte(nil), buf...))
	if p.events == nil {
		return nil
	}
	return p.events(buf)
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
