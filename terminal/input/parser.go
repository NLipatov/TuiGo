package input

import (
	"bytes"
	"unicode/utf8"

	"github.com/NLipatov/tuigo/ansi"
)

type Parser struct {
	buf []byte
}

func NewParser() *Parser {
	return &Parser{
		buf: make([]byte, 0, utf8.UTFMax),
	}
}

func (i *Parser) Feed(buf []byte) ParseResult {
	i.buf = append(i.buf, buf...)
	events := make([]Event, 0)
	for len(i.buf) > 0 {
		event, n, status := i.parseEvent(i.buf)
		if status != parseDone {
			break
		}
		events = append(events, event)
		i.buf = i.buf[n:]
	}
	return ParseResult{
		Events:           events,
		HasPendingEscape: i.hasPendingEscape(),
	}
}

// FlushPendingEscape emits a pending ESC byte as KeyEsc after the caller's input timeout expires.
// Any buffered bytes after ESC are parsed as regular input.
func (i *Parser) FlushPendingEscape() ParseResult {
	if !i.hasPendingEscape() {
		return ParseResult{}
	}
	rest := append([]byte(nil), i.buf[1:]...)
	i.buf = i.buf[:0]
	result := i.Feed(rest)
	result.Events = append([]Event{{
		Type: EventTypeKey,
		Key:  KeyEvent{Code: KeyEsc},
	}}, result.Events...)
	return result
}

func (i *Parser) hasPendingEscape() bool {
	return len(i.buf) > 0 && i.isEscapeByte(i.buf[0])
}

func (i *Parser) parseEvent(buf []byte) (Event, int, parseStatus) {
	if len(buf) == 0 {
		return Event{}, 0, parseNeedMore
	}
	if event, n, status := i.parseEscEvent(buf); status != parseNoMatch {
		return event, n, status
	}
	if event, n, status := i.parseControlEvent(buf); status != parseNoMatch {
		return Event{Type: EventTypeKey, Key: event}, n, status
	}
	event, n, status := i.parseRuneEvent(buf)
	return Event{Type: EventTypeKey, Key: event}, n, status
}

func (i *Parser) parseEscEvent(buf []byte) (Event, int, parseStatus) {
	if len(buf) == 0 {
		return Event{}, 0, parseNeedMore
	}
	if !i.isEscapeByte(buf[0]) {
		return Event{}, 0, parseNoMatch
	}
	if len(buf) == 1 {
		return Event{}, 0, parseNeedMore
	}
	for _, candidate := range escapeSequences {
		if bytes.HasPrefix(buf, candidate.sequence) {
			return Event{Type: EventTypeKey, Key: candidate.event}, len(candidate.sequence), parseDone
		}
		if bytes.HasPrefix(candidate.sequence, buf) {
			return Event{}, 0, parseNeedMore
		}
	}
	if buf[1] == '[' {
		return i.parseCSISequence(buf)
	}
	if buf[1] == 'O' {
		return i.parseSS3Sequence(buf)
	}
	if event, n, status := i.parseControlEvent(buf[1:]); status != parseNoMatch {
		if status != parseDone {
			return Event{}, 0, status
		}
		event.Mod |= ModAlt
		return Event{Type: EventTypeKey, Key: event}, n + 1, parseDone
	}
	event, n, status := i.parseRuneEvent(buf[1:])
	if status != parseDone {
		return Event{}, 0, status
	}
	event.Mod |= ModAlt
	return Event{Type: EventTypeKey, Key: event}, n + 1, parseDone
}

func (i *Parser) parseCSISequence(buf []byte) (Event, int, parseStatus) {
	if event, n, status := legacyMouseEventFromCSI(buf); status != parseNoMatch {
		return event, n, status
	}
	finalIdx, ok := findCSIFinal(buf)
	if !ok {
		return Event{}, 0, parseNeedMore
	}
	params := string(buf[2:finalIdx])
	final := buf[finalIdx]
	if isSGRMouseSequence(params, final) {
		mouse, ok := mouseEventFromSGR(params, final)
		if !ok {
			return Event{Type: EventTypeKey, Key: KeyEvent{Code: KeyUnknown}}, finalIdx + 1, parseDone
		}
		return Event{Type: EventTypeMouse, Mouse: mouse}, finalIdx + 1, parseDone
	}
	event, ok := eventForCSIFinal(final, params)
	if !ok {
		return Event{Type: EventTypeKey, Key: KeyEvent{Code: KeyUnknown}}, finalIdx + 1, parseDone
	}
	return Event{Type: EventTypeKey, Key: event}, finalIdx + 1, parseDone
}

func findCSIFinal(buf []byte) (int, bool) {
	for i := 2; i < len(buf); i++ {
		if buf[i] >= 0x40 && buf[i] <= 0x7e {
			return i, true
		}
	}
	return 0, false
}

func (i *Parser) parseSS3Sequence(buf []byte) (Event, int, parseStatus) {
	for _, candidate := range ss3Sequences {
		if bytes.HasPrefix(buf, candidate.sequence) {
			return Event{Type: EventTypeKey, Key: candidate.event}, len(candidate.sequence), parseDone
		}
		if bytes.HasPrefix(candidate.sequence, buf) {
			return Event{}, 0, parseNeedMore
		}
	}
	return Event{Type: EventTypeKey, Key: KeyEvent{Code: KeyUnknown}}, len(ansi.SS3) + 1, parseDone
}

func (i *Parser) parseControlEvent(buf []byte) (KeyEvent, int, parseStatus) {
	if len(buf) == 0 {
		return KeyEvent{}, 0, parseNeedMore
	}

	b := buf[0]
	for _, candidate := range controlBytes {
		if b == candidate.b {
			return candidate.event, 1, parseDone
		}
	}
	if b >= 0x01 && b <= 0x1a {
		r := 'a' + rune(b) - 1
		return KeyEvent{Code: KeyRune, Text: string(r), Mod: ModCtrl}, 1, parseDone
	}
	return KeyEvent{}, 0, parseNoMatch
}

func (i *Parser) parseRuneEvent(buf []byte) (KeyEvent, int, parseStatus) {
	if !utf8.FullRune(buf) {
		return KeyEvent{}, 0, parseNeedMore
	}
	r, n := utf8.DecodeRune(buf)
	return KeyEvent{
		Code: KeyRune,
		Text: string(r),
		Mod:  ModNone,
	}, n, parseDone
}

func (i *Parser) isEscapeByte(b byte) bool {
	return byte(ansi.ESCAPEByte) == b
}
