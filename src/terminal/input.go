package terminal

import (
	"context"
	"errors"
	"io"
	"sync"
)

var (
	ErrNilContext   = errors.New("terminal: nil context")
	ErrNilReader    = errors.New("terminal: nil input reader")
	ErrNilParser    = errors.New("terminal: nil input parser")
	ErrNilEventChan = errors.New("terminal: nil input event channel")
)

type Parser interface {
	Feed(buf []byte) []Event
}

type Event struct {
	Code KeyCode
	Text string
	Mod  KeyMod
}

type Input struct {
	ctx          context.Context
	reader       io.ReadCloser
	buf          [64]byte
	parser       Parser
	out          chan<- Event
	once         sync.Once
	closeOnceErr error
}

func NewInput(ctx context.Context, reader io.ReadCloser, parser Parser, ch chan<- Event) (Input, error) {
	if ctx == nil {
		return Input{}, ErrNilContext
	}
	if reader == nil {
		return Input{}, ErrNilReader
	}
	if parser == nil {
		return Input{}, ErrNilParser
	}
	if ch == nil {
		return Input{}, ErrNilEventChan
	}
	return Input{
		ctx:    ctx,
		reader: reader,
		out:    ch,
		parser: parser,
	}, nil
}

func (i *Input) Listen() error {
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-i.ctx.Done():
			_ = i.Close()
		case <-done:
		}
	}()
	for {
		n, err := i.reader.Read(i.buf[:])
		if err != nil {
			if i.ctx.Err() != nil {
				return i.ctx.Err()
			}
			return err
		}
		events := i.parser.Feed(i.buf[:n])
		for _, event := range events {
			select {
			case <-i.ctx.Done():
				return i.ctx.Err()
			case i.out <- event:
			}
		}
	}
}

func (i *Input) Close() error {
	i.once.Do(func() {
		i.closeOnceErr = i.reader.Close()
	})
	return i.closeOnceErr
}
