package input

import (
	"bytes"
	"strconv"
	"strings"
	"tuigo/ansi"
	"unicode/utf8"
)

var escapeSequences = []struct {
	sequence []byte
	event    Event
}{
	{[]byte(ansi.CSI + "A"), Event{Code: KeyUp}},
	{[]byte(ansi.CSI + "B"), Event{Code: KeyDown}},
	{[]byte(ansi.CSI + "C"), Event{Code: KeyRight}},
	{[]byte(ansi.CSI + "D"), Event{Code: KeyLeft}},
	{[]byte(ansi.CSI + "H"), Event{Code: KeyHome}},
	{[]byte(ansi.CSI + "F"), Event{Code: KeyEnd}},
	{[]byte(ansi.CSI + "3~"), Event{Code: KeyDelete}},
	{[]byte(ansi.CSI + "5~"), Event{Code: KeyPageUp}},
	{[]byte(ansi.CSI + "6~"), Event{Code: KeyPageDown}},
	{[]byte(ansi.CSI + "Z"), Event{Code: KeyTab, Mod: ModShift}},
	{[]byte(ansi.CSI + "11~"), Event{Code: KeyF1}},
	{[]byte(ansi.CSI + "12~"), Event{Code: KeyF2}},
	{[]byte(ansi.CSI + "13~"), Event{Code: KeyF3}},
	{[]byte(ansi.CSI + "14~"), Event{Code: KeyF4}},
	{[]byte(ansi.CSI + "15~"), Event{Code: KeyF5}},
	{[]byte(ansi.CSI + "17~"), Event{Code: KeyF6}},
	{[]byte(ansi.CSI + "18~"), Event{Code: KeyF7}},
	{[]byte(ansi.CSI + "19~"), Event{Code: KeyF8}},
	{[]byte(ansi.CSI + "20~"), Event{Code: KeyF9}},
	{[]byte(ansi.CSI + "21~"), Event{Code: KeyF10}},
	{[]byte(ansi.CSI + "23~"), Event{Code: KeyF11}},
	{[]byte(ansi.CSI + "24~"), Event{Code: KeyF12}},
}

var ss3Sequences = []struct {
	sequence []byte
	event    Event
}{
	{[]byte(ansi.SS3 + "H"), Event{Code: KeyHome}},
	{[]byte(ansi.SS3 + "F"), Event{Code: KeyEnd}},
	{[]byte(ansi.SS3 + "P"), Event{Code: KeyF1}},
	{[]byte(ansi.SS3 + "Q"), Event{Code: KeyF2}},
	{[]byte(ansi.SS3 + "R"), Event{Code: KeyF3}},
	{[]byte(ansi.SS3 + "S"), Event{Code: KeyF4}},
}

var controlBytes = []struct {
	b     byte
	event Event
}{
	{'\r', Event{Code: KeyEnter}},
	{'\n', Event{Code: KeyEnter}},
	{'\t', Event{Code: KeyTab}},
	{0x7f, Event{Code: KeyBackspace}},
	{0x08, Event{Code: KeyBackspace}},
}

type parseStatus uint8

const (
	parseNoMatch parseStatus = iota
	parseNeedMore
	parseDone
)

type InputParser struct {
	buf []byte
}

func NewInputParser() *InputParser {
	return &InputParser{
		buf: make([]byte, 0, utf8.UTFMax),
	}
}

func (i *InputParser) Timeout() ParseResult {
	if i.needsTimeout() == false {
		return ParseResult{}
	}
	rest := append([]byte(nil), i.buf[1:]...)
	// drain i.buf
	i.buf = i.buf[:0]
	// create events from rest
	events := []Event{{Code: KeyEsc}}
	for _, res := range i.Feed(rest).Events {
		events = append(events, res)
	}
	return ParseResult{
		Events: events,
	}
}

func (i *InputParser) Feed(buf []byte) ParseResult {
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
		Events:       events,
		NeedsTimeout: i.needsTimeout(),
	}
}

func (i *InputParser) needsTimeout() bool {
	return len(i.buf) > 0 && i.isEscapeByte(i.buf[0])
}

func (i *InputParser) parseEvent(buf []byte) (Event, int, parseStatus) {
	if len(buf) == 0 {
		return Event{}, 0, parseNeedMore
	}
	if event, n, status := i.parseEscEvent(buf); status != parseNoMatch {
		return event, n, status
	}
	if event, n, status := i.parseControlEvent(buf); status != parseNoMatch {
		return event, n, status
	}
	return i.parseRuneEvent(buf)
}

func (i *InputParser) parseEscEvent(buf []byte) (Event, int, parseStatus) {
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
		// full escape sequence prefix
		if bytes.HasPrefix(buf, candidate.sequence) {
			return candidate.event, len(candidate.sequence), parseDone
		}
		// partial escape sequence prefix
		if bytes.HasPrefix(candidate.sequence, buf) {
			return Event{}, 0, parseNeedMore
		}
	}
	// if CSI
	if buf[1] == '[' {
		return i.parseCSISequence(buf)
	}
	// If SS3
	if buf[1] == 'O' {
		return i.parseSS3Sequence(buf)
	}
	// Is control byte
	if event, n, status := i.parseControlEvent(buf[1:]); status != parseNoMatch {
		if status != parseDone {
			return Event{}, 0, status
		}
		event.Mod |= ModAlt
		return event, n + 1, parseDone
	}
	// Is rune
	event, n, status := i.parseRuneEvent(buf[1:])
	if status != parseDone {
		return Event{}, 0, status
	}
	event.Mod |= ModAlt
	return event, n + 1, parseDone
}

func (i *InputParser) parseCSISequence(buf []byte) (Event, int, parseStatus) {
	finalIdx, ok := findCSIFinal(buf)
	if !ok {
		return Event{}, 0, parseNeedMore
	}
	params := string(buf[2:finalIdx])
	final := buf[finalIdx]
	event, ok := eventForCSIFinal(final, params)
	if !ok {
		return Event{Code: KeyUnknown}, finalIdx + 1, parseDone
	}
	return event, finalIdx + 1, parseDone
}

func (i *InputParser) parseSS3Sequence(buf []byte) (Event, int, parseStatus) {
	for _, candidate := range ss3Sequences {
		if bytes.HasPrefix(buf, candidate.sequence) {
			return candidate.event, len(candidate.sequence), parseDone
		}
		if bytes.HasPrefix(candidate.sequence, buf) {
			return Event{}, 0, parseNeedMore
		}
	}
	return Event{Code: KeyUnknown}, len(ansi.SS3) + 1, parseDone
}

func (i *InputParser) parseControlEvent(buf []byte) (Event, int, parseStatus) {
	if len(buf) == 0 {
		return Event{}, 0, parseNeedMore
	}

	b := buf[0]
	for _, candidate := range controlBytes {
		if b == candidate.b {
			return candidate.event, 1, parseDone
		}
	}
	if b >= 0x01 && b <= 0x1a {
		r := 'a' + rune(b) - 1
		return Event{Code: KeyRune, Text: string(r), Mod: ModCtrl}, 1, parseDone
	}
	return Event{}, 0, parseNoMatch
}

func (i *InputParser) parseRuneEvent(buf []byte) (Event, int, parseStatus) {
	if !utf8.FullRune(buf) {
		return Event{}, 0, parseNeedMore
	}
	r, n := utf8.DecodeRune(buf)
	return Event{
		Code: KeyRune,
		Text: string(r),
		Mod:  ModNone,
	}, n, parseDone
}

func (i *InputParser) isEscapeByte(b byte) bool {
	return byte(ansi.ESCAPEByte) == b
}

func findCSIFinal(buf []byte) (int, bool) {
	for i := 2; i < len(buf); i++ {
		if buf[i] >= 0x40 && buf[i] <= 0x7e {
			return i, true
		}
	}
	return 0, false
}

func eventForCSIFinal(final byte, params string) (Event, bool) {
	switch final {
	case 'A':
		return Event{Code: KeyUp, Mod: modFromCSIParams(params)}, true
	case 'B':
		return Event{Code: KeyDown, Mod: modFromCSIParams(params)}, true
	case 'C':
		return Event{Code: KeyRight, Mod: modFromCSIParams(params)}, true
	case 'D':
		return Event{Code: KeyLeft, Mod: modFromCSIParams(params)}, true
	case 'H':
		return Event{Code: KeyHome, Mod: modFromCSIParams(params)}, true
	case 'F':
		return Event{Code: KeyEnd, Mod: modFromCSIParams(params)}, true
	case 'Z':
		return Event{Code: KeyTab, Mod: ModShift}, true
	case '~':
		return eventForCSITilde(params)
	default:
		return Event{}, false
	}
}

func eventForCSITilde(params string) (Event, bool) {
	values := splitCSIParams(params)
	if len(values) == 0 {
		return Event{}, false
	}

	event := Event{Code: keyCodeForCSITilde(values[0])}
	if event.Code == KeyUnknown {
		return event, true
	}
	if len(values) > 1 {
		event.Mod = xtermModifier(values[1])
	}
	return event, true
}

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
