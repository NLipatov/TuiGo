package application

import (
	"io"
	"tuigo/ansi"
)

type Session struct {
	writer io.StringWriter
}

func NewSession(writer io.StringWriter) Session {
	return Session{
		writer: writer,
	}
}

func (s *Session) Start() error {
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
	if _, err := s.writer.WriteString(string(command)); err != nil {
		return err
	}
	return nil
}
