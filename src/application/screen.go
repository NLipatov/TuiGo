package application

import (
	"io"
	"tuigo/ansi"
)

type Screen struct {
	writer io.StringWriter
}

func NewScreen(writer io.StringWriter) Screen {
	return Screen{
		writer: writer,
	}
}

func (s *Screen) Start() error {
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

func (s *Screen) Close() error {
	if err := s.ansiCommand(ansi.SHOW_CURSOR); err != nil {
		return err
	}
	if err := s.ansiCommand(ansi.EXIT_ALTERNATE_SCREEN); err != nil {
		return err
	}
	return nil
}

func (s *Screen) ansiCommand(command ansi.ANSIEscapeSequence) error {
	if _, err := s.writer.WriteString(string(command)); err != nil {
		return err
	}
	return nil
}
