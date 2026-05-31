package input

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

const inputTimeout = time.Millisecond * 50

var (
	ErrNilContext     = errors.New("terminal: nil context")
	ErrNilReader      = errors.New("terminal: nil input reader")
	ErrNilParser      = errors.New("terminal: nil input parser")
	ErrNilEventChan   = errors.New("terminal: nil input event channel")
	ErrReadChanClosed = errors.New("terminal: input read channel closed")
)

type ParseResult struct {
	Events       []Event
	NeedsTimeout bool
}

type Parser interface {
	Feed([]byte) ParseResult
	Timeout() ParseResult
}

type Listener struct {
	ctx             context.Context
	reader          io.ReadCloser
	parser          Parser
	out             chan<- Event
	closeReaderOnce sync.Once
	closeReaderErr  error
}

func NewListener(ctx context.Context, reader io.ReadCloser, parser Parser, ch chan<- Event) (Listener, error) {
	if ctx == nil {
		return Listener{}, ErrNilContext
	}
	if reader == nil {
		return Listener{}, ErrNilReader
	}
	if parser == nil {
		return Listener{}, ErrNilParser
	}
	if ch == nil {
		return Listener{}, ErrNilEventChan
	}
	return Listener{
		ctx:    ctx,
		reader: reader,
		out:    ch,
		parser: parser,
	}, nil
}

func (i *Listener) Listen() error {
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-i.ctx.Done():
			_ = i.Close()
		case <-done:
		}
	}()
	type readResult struct {
		N    int
		Data [64]byte
		Err  error
	}
	readCh := make(chan readResult)
	go func() {
		defer close(readCh)
		for {
			var result readResult
			result.N, result.Err = i.reader.Read(result.Data[:])
			select {
			case <-i.ctx.Done():
				return
			case readCh <- result:
			}
			if result.Err != nil {
				return
			}
		}
	}()
	timeout := time.NewTimer(inputTimeout)
	i.stopTimer(timeout)
	for {
		select {
		case read, ok := <-readCh:
			if !ok {
				if i.ctx.Err() != nil {
					return i.ctx.Err()
				}
				return ErrReadChanClosed
			}
			if read.Err != nil {
				if i.ctx.Err() != nil {
					return i.ctx.Err()
				}
				return read.Err
			}
			results := i.parser.Feed(read.Data[:read.N])
			if err := i.handleParserResult(results, timeout); err != nil {
				return err
			}
		case <-timeout.C:
			results := i.parser.Timeout()
			if err := i.handleParserResult(results, timeout); err != nil {
				return err
			}
		}
	}
}

func (i *Listener) handleParserResult(result ParseResult, timeout *time.Timer) error {
	for _, event := range result.Events {
		select {
		case <-i.ctx.Done():
			return i.ctx.Err()
		case i.out <- event:
		}
	}
	if result.NeedsTimeout {
		i.stopTimer(timeout)
		_ = timeout.Reset(inputTimeout)
		return nil
	}
	i.stopTimer(timeout)
	return nil
}

func (i *Listener) stopTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
}

func (i *Listener) Close() error {
	i.closeReaderOnce.Do(func() {
		// TODO: avoid closing process stdin when stopping the input listener.
		i.closeReaderErr = i.reader.Close()
	})
	return i.closeReaderErr
}
