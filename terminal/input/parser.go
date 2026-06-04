package input

import (
	"bytes"
	"strconv"
	"strings"
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

// FlushPendingEscape emits a pending ESC byte as KeyEsc after the caller's input timeout expires.
// Any buffered bytes after ESC are parsed as regular input.
func (i *Parser) FlushPendingEscape() ParseResult {
	if !i.hasPendingEscape() {
		return ParseResult{}
	}
	rest := append([]byte(nil), i.buf[1:]...)
	i.buf = i.buf[:0]
	result := i.Feed(rest)
	result.Events = append([]KeyEvent{{Code: KeyEsc}}, result.Events...)
	return result
}

func (i *Parser) Feed(buf []byte) ParseResult {
	i.buf = append(i.buf, buf...)
	events := make([]KeyEvent, 0)
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

func (i *Parser) hasPendingEscape() bool {
	return len(i.buf) > 0 && i.isEscapeByte(i.buf[0])
}

func (i *Parser) parseEvent(buf []byte) (KeyEvent, int, parseStatus) {
	if len(buf) == 0 {
		return KeyEvent{}, 0, parseNeedMore
	}
	if event, n, status := i.parseEscEvent(buf); status != parseNoMatch {
		return event, n, status
	}
	if event, n, status := i.parseControlEvent(buf); status != parseNoMatch {
		return event, n, status
	}
	return i.parseRuneEvent(buf)
}

func (i *Parser) parseEscEvent(buf []byte) (KeyEvent, int, parseStatus) {
	if len(buf) == 0 {
		return KeyEvent{}, 0, parseNeedMore
	}
	if !i.isEscapeByte(buf[0]) {
		return KeyEvent{}, 0, parseNoMatch
	}
	if len(buf) == 1 {
		return KeyEvent{}, 0, parseNeedMore
	}
	for _, candidate := range escapeSequences {
		if bytes.HasPrefix(buf, candidate.sequence) {
			return candidate.event, len(candidate.sequence), parseDone
		}
		if bytes.HasPrefix(candidate.sequence, buf) {
			return KeyEvent{}, 0, parseNeedMore
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
			return KeyEvent{}, 0, status
		}
		event.Mod |= ModAlt
		return event, n + 1, parseDone
	}
	event, n, status := i.parseRuneEvent(buf[1:])
	if status != parseDone {
		return KeyEvent{}, 0, status
	}
	event.Mod |= ModAlt
	return event, n + 1, parseDone
}

func (i *Parser) parseCSISequence(buf []byte) (KeyEvent, int, parseStatus) {
	finalIdx, ok := findCSIFinal(buf)
	if !ok {
		return KeyEvent{}, 0, parseNeedMore
	}
	params := string(buf[2:finalIdx])
	final := buf[finalIdx]
	event, ok := eventForCSIFinal(final, params)
	if !ok {
		return KeyEvent{Code: KeyUnknown}, finalIdx + 1, parseDone
	}
	return event, finalIdx + 1, parseDone
}

func findCSIFinal(buf []byte) (int, bool) {
	for i := 2; i < len(buf); i++ {
		if buf[i] >= 0x40 && buf[i] <= 0x7e {
			return i, true
		}
	}
	return 0, false
}

func (i *Parser) parseSS3Sequence(buf []byte) (KeyEvent, int, parseStatus) {
	for _, candidate := range ss3Sequences {
		if bytes.HasPrefix(buf, candidate.sequence) {
			return candidate.event, len(candidate.sequence), parseDone
		}
		if bytes.HasPrefix(candidate.sequence, buf) {
			return KeyEvent{}, 0, parseNeedMore
		}
	}
	return KeyEvent{Code: KeyUnknown}, len(ansi.SS3) + 1, parseDone
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

func eventForCSIFinal(final byte, params string) (KeyEvent, bool) {
	switch final {
	case 'A':
		return KeyEvent{Code: KeyUp, Mod: modFromCSIParams(params)}, true
	case 'B':
		return KeyEvent{Code: KeyDown, Mod: modFromCSIParams(params)}, true
	case 'C':
		return KeyEvent{Code: KeyRight, Mod: modFromCSIParams(params)}, true
	case 'D':
		return KeyEvent{Code: KeyLeft, Mod: modFromCSIParams(params)}, true
	case 'H':
		return KeyEvent{Code: KeyHome, Mod: modFromCSIParams(params)}, true
	case 'F':
		return KeyEvent{Code: KeyEnd, Mod: modFromCSIParams(params)}, true
	case 'Z':
		return KeyEvent{Code: KeyTab, Mod: ModShift}, true
	case '~':
		return eventForCSITilde(params)
	default:
		return KeyEvent{}, false
	}
}

func eventForCSITilde(params string) (KeyEvent, bool) {
	values := splitCSIParams(params)
	if len(values) == 0 {
		return KeyEvent{}, false
	}

	event := KeyEvent{Code: keyCodeForCSITilde(values[0])}
	if event.Code == KeyUnknown {
		return event, true
	}
	if len(values) > 1 {
		event.Mod = xtermModifier(values[1])
	}
	return event, true
}

//nolint:cyclop // This switch is a flat CSI numeric key mapping table.
func keyCodeForCSITilde(value int) KeyCode {
	switch value {
	case 2:
		return KeyInsert
	case 3:
		return KeyDelete
	case 5:
		return KeyPageUp
	case 6:
		return KeyPageDown
	case 11:
		return KeyF1
	case 12:
		return KeyF2
	case 13:
		return KeyF3
	case 14:
		return KeyF4
	case 15:
		return KeyF5
	case 17:
		return KeyF6
	case 18:
		return KeyF7
	case 19:
		return KeyF8
	case 20:
		return KeyF9
	case 21:
		return KeyF10
	case 23:
		return KeyF11
	case 24:
		return KeyF12
	default:
		return KeyUnknown
	}
}

func modFromCSIParams(params string) KeyMod {
	values := splitCSIParams(params)
	if len(values) < 2 {
		return ModNone
	}
	return xtermModifier(values[1])
}

func splitCSIParams(params string) []int {
	if params == "" {
		return nil
	}

	parts := strings.Split(params, ";")
	values := make([]int, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			values = append(values, 0)
			continue
		}
		value, err := strconv.Atoi(part)
		if err != nil {
			return nil
		}
		values = append(values, value)
	}
	return values
}

func xtermModifier(value int) KeyMod {
	bits := value - 1
	if bits <= 0 {
		return ModNone
	}

	var mod KeyMod
	if bits&1 != 0 {
		mod |= ModShift
	}
	if bits&2 != 0 {
		mod |= ModAlt
	}
	if bits&4 != 0 {
		mod |= ModCtrl
	}
	return mod
}
