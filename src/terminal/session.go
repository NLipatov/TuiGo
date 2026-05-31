package terminal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"tuigo/ansi"
	"tuigo/terminal/input"
	"tuigo/terminal/resize"
)

var (
	ErrNilContext = errors.New("terminal: nil context")
	ErrNilInput   = errors.New("terminal: nil input")
	ErrNilOutput  = errors.New("terminal: nil output")
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
	if err := s.setupTerminal(); err != nil {
		return nil, err
	}
	events, err := s.startEventLoop()
	if err != nil {
		if unsetTerminalErr := s.restoreTerminal(); unsetTerminalErr != nil {
			return nil, fmt.Errorf("failed to start event loop: %w; terminal was not restored: %v", err, unsetTerminalErr)
		}
		return nil, err
	}
	s.events = events
	s.started = true
	return s.events, nil
}

func (s *Session) setupTerminal() error {
	if err := s.device.EnableRawMode(); err != nil {
		return err
	}
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

func (s *Session) restoreTerminal() error {
	if err := s.ansiCommand(ansi.SHOW_CURSOR); err != nil {
		return err
	}
	if err := s.ansiCommand(ansi.EXIT_ALTERNATE_SCREEN); err != nil {
		return err
	}
	if err := s.device.RestoreInitialMode(); err != nil {
		return err
	}
	return nil
}

func (s *Session) startEventLoop() (chan Event, error) {
	resizeCh := make(chan resize.Event)
	resizeListener := resize.NewListener(s.ctx, resizeCh, &s.device)
	keyCh := make(chan input.Event)
	keyListener, err := input.NewListener(s.ctx, s.in, input.NewInputParser(), keyCh)
	if err != nil {
		return nil, err
	}
	outCh := make(chan Event)
	errCh := make(chan error, 2)
	go func() {
		defer close(outCh)
		for {
			select {
			case <-s.ctx.Done():
				return
			case err, ok := <-errCh:
				if !ok {
					return
				}
				select {
				case <-s.ctx.Done():
					return
				case outCh <- Event{
					Type: EventError,
					Err:  err,
				}:
				}
			case event, ok := <-resizeCh:
				if !ok {
					return
				}
				select {
				case <-s.ctx.Done():
					return
				case outCh <- Event{
					Type:   EventResize,
					Resize: event,
				}:
				}
			case event, ok := <-keyCh:
				if !ok {
					return
				}
				select {
				case <-s.ctx.Done():
					return
				case outCh <- Event{
					Type: EventKey,
					Key:  event,
				}:
				}
			}
		}
	}()
	go func() {
		if err := resizeListener.Listen(); err != nil {
			if s.ctx.Err() != nil {
				return
			}
			select {
			case <-s.ctx.Done():
				return
			case errCh <- err:
			}
		}
	}()
	go func() {
		if err := keyListener.Listen(); err != nil {
			if s.ctx.Err() != nil {
				return
			}
			select {
			case <-s.ctx.Done():
				return
			case errCh <- err:
			}
		}
	}()
	return outCh, nil
}

func (s *Session) Close() error {
	return s.restoreTerminal()
}

func (s *Session) ansiCommand(command ansi.ANSIEscapeSequence) error {
	if writer, ok := s.out.(io.StringWriter); ok {
		_, err := writer.WriteString(string(command))
		return err
	}
	_, err := s.out.Write([]byte(command))
	return err
}
