package terminal

import (
	"context"
	"errors"
	"io"
	"os"
	"tuigo/ansi"
	"tuigo/terminal/resize"
)

var (
	ErrNilContext = errors.New("nil context")
	ErrNilInput   = errors.New("nil input")
	ErrNilOutput  = errors.New("nil output")
)

type Session struct {
	ctx     context.Context
	in      *os.File
	out     io.Writer
	device  Device
	started bool
	events  chan Event
}

func NewSession(ctx context.Context, in *os.File, out io.Writer) (Session, error) {
	if ctx == nil {
		return Session{}, ErrNilContext
	}
	if in == nil {
		return Session{}, ErrNilInput
	}
	if out == nil {
		return Session{}, ErrNilOutput
	}
	return Session{
		ctx:    ctx,
		in:     in,
		out:    out,
		device: NewDevice(int(in.Fd())),
	}, nil
}

func (s *Session) Start() (<-chan Event, error) {
	if s.started {
		return s.events, nil
	}
	if err := s.startTerminal(); err != nil {
		return nil, err
	}
	s.events = s.startEventLoop()
	s.started = true
	return s.events, nil
}

func (s *Session) startTerminal() error {
	if err := s.ansiCommand(ansi.ENTER_ALTERNATE_SCREEN); err != nil {
		return err
	}
	if err := s.ansiCommand(ansi.HIDE_CURSOR); err != nil {
		return err
	}
	if err := s.ansiCommand(ansi.CLEAR_SCREEN); err != nil {
		return err
	}
	if err := s.ansiCommand(ansi.CURSOR_HOME); err != nil {
		return err
	}
	return nil
}

func (s *Session) startEventLoop() chan Event {
	resizeIn := make(chan resize.Event)
	out := make(chan Event)
	go func() {
		defer close(out)
		for {
			select {
			case <-s.ctx.Done():
				return
			case event, ok := <-resizeIn:
				if !ok {
					return
				}
				select {
				case <-s.ctx.Done():
					return
				case out <- Event{
					Type:   EventResize,
					Resize: event,
				}:
				}
			}
		}
	}()
	resizeListener := resize.NewListener(s.ctx, resizeIn, &s.device)
	go func() {
		// ToDo: handle error
		_ = resizeListener.Listen()
	}()
	return out
}

func (s *Session) Close() error {
	if err := s.ansiCommand(ansi.SHOW_CURSOR); err != nil {
		return err
	}
	if err := s.ansiCommand(ansi.EXIT_ALTERNATE_SCREEN); err != nil {
		return err
	}
	return nil
}

func (s *Session) ansiCommand(command ansi.ANSIEscapeSequence) error {
	if writer, ok := s.out.(io.StringWriter); ok {
		_, err := writer.WriteString(string(command))
		return err
	}
	_, err := s.out.Write([]byte(command))
	return err
}
