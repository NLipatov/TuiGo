package input

import (
	"testing"

	"github.com/NLipatov/tuigo/keyboard"
	"github.com/NLipatov/tuigo/mouse"
)

func TestParserParsesRunes(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("aЖ"))
	want := []keyboard.KeyEvent{
		{Code: keyboard.KeyRune, Text: "a", Mod: keyboard.ModNone},
		{Code: keyboard.KeyRune, Text: "Ж", Mod: keyboard.ModNone},
	}
	assertEvents(t, got, want)
}

func TestParserWaitsForIncompleteUTF8Rune(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte{0xd0})
	assertEvents(t, got, nil)

	got = parser.Feed([]byte{0x96})
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyRune, Text: "Ж", Mod: keyboard.ModNone}})
}

func TestParserParsesControlBytes(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want keyboard.KeyEvent
	}{
		{name: "carriage return", in: []byte{'\r'}, want: keyboard.KeyEvent{Code: keyboard.KeyEnter}},
		{name: "line feed", in: []byte{'\n'}, want: keyboard.KeyEvent{Code: keyboard.KeyEnter}},
		{name: "tab", in: []byte{'\t'}, want: keyboard.KeyEvent{Code: keyboard.KeyTab}},
		{name: "delete backspace", in: []byte{0x7f}, want: keyboard.KeyEvent{Code: keyboard.KeyBackspace}},
		{name: "ctrl h backspace", in: []byte{0x08}, want: keyboard.KeyEvent{Code: keyboard.KeyBackspace}},
		{name: "ctrl c", in: []byte{0x03}, want: keyboard.KeyEvent{Code: keyboard.KeyRune, Text: "c", Mod: keyboard.ModCtrl}},
		{name: "ctrl d", in: []byte{0x04}, want: keyboard.KeyEvent{Code: keyboard.KeyRune, Text: "d", Mod: keyboard.ModCtrl}},
		{name: "ctrl z", in: []byte{0x1a}, want: keyboard.KeyEvent{Code: keyboard.KeyRune, Text: "z", Mod: keyboard.ModCtrl}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []keyboard.KeyEvent{tt.want})
		})
	}
}

func TestParserParsesEscapeSequences(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want keyboard.KeyEvent
	}{
		{name: "up", in: []byte("\x1b[A"), want: keyboard.KeyEvent{Code: keyboard.KeyUp}},
		{name: "down", in: []byte("\x1b[B"), want: keyboard.KeyEvent{Code: keyboard.KeyDown}},
		{name: "right", in: []byte("\x1b[C"), want: keyboard.KeyEvent{Code: keyboard.KeyRight}},
		{name: "left", in: []byte("\x1b[D"), want: keyboard.KeyEvent{Code: keyboard.KeyLeft}},
		{name: "home csi", in: []byte("\x1b[H"), want: keyboard.KeyEvent{Code: keyboard.KeyHome}},
		{name: "end csi", in: []byte("\x1b[F"), want: keyboard.KeyEvent{Code: keyboard.KeyEnd}},
		{name: "home ss3", in: []byte("\x1bOH"), want: keyboard.KeyEvent{Code: keyboard.KeyHome}},
		{name: "end ss3", in: []byte("\x1bOF"), want: keyboard.KeyEvent{Code: keyboard.KeyEnd}},
		{name: "delete", in: []byte("\x1b[3~"), want: keyboard.KeyEvent{Code: keyboard.KeyDelete}},
		{name: "page up", in: []byte("\x1b[5~"), want: keyboard.KeyEvent{Code: keyboard.KeyPageUp}},
		{name: "page down", in: []byte("\x1b[6~"), want: keyboard.KeyEvent{Code: keyboard.KeyPageDown}},
		{name: "shift tab", in: []byte("\x1b[Z"), want: keyboard.KeyEvent{Code: keyboard.KeyTab, Mod: keyboard.ModShift}},
		{name: "insert", in: []byte("\x1b[2~"), want: keyboard.KeyEvent{Code: keyboard.KeyInsert}},
		{name: "f1 ss3", in: []byte("\x1bOP"), want: keyboard.KeyEvent{Code: keyboard.KeyF1}},
		{name: "f4 ss3", in: []byte("\x1bOS"), want: keyboard.KeyEvent{Code: keyboard.KeyF4}},
		{name: "f5 csi", in: []byte("\x1b[15~"), want: keyboard.KeyEvent{Code: keyboard.KeyF5}},
		{name: "f12 csi", in: []byte("\x1b[24~"), want: keyboard.KeyEvent{Code: keyboard.KeyF12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []keyboard.KeyEvent{tt.want})
		})
	}
}

func TestParserFlushesPendingEscAfterTimeout(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b"))
	assertEvents(t, got, nil)
	if !got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = false, want true")
	}

	got = parser.FlushPendingEscape()
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyEsc}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesSplitEscapeSequence(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b"))
	assertEvents(t, got, nil)
	if !got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = false, want true")
	}

	got = parser.Feed([]byte("[A"))
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyUp}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesEscapeSequencesBeforeFollowingInput(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[Aa"))
	want := []keyboard.KeyEvent{
		{Code: keyboard.KeyUp},
		{Code: keyboard.KeyRune, Text: "a", Mod: keyboard.ModNone},
	}
	assertEvents(t, got, want)
}

func TestParserWaitsForIncompleteEscapeSequence(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b["))
	assertEvents(t, got, nil)
	if !got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = false, want true")
	}

	got = parser.Feed([]byte("A"))
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyUp}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesAltRune(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1bx"))
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyRune, Text: "x", Mod: keyboard.ModAlt}})
}

func TestParserParsesAltRuneSplitAcrossFeeds(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b"))
	assertEvents(t, got, nil)
	if !got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = false, want true")
	}

	got = parser.Feed([]byte("x"))
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyRune, Text: "x", Mod: keyboard.ModAlt}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesCSIModifiers(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want keyboard.KeyEvent
	}{
		{name: "ctrl up", in: []byte("\x1b[1;5A"), want: keyboard.KeyEvent{Code: keyboard.KeyUp, Mod: keyboard.ModCtrl}},
		{name: "shift left", in: []byte("\x1b[1;2D"), want: keyboard.KeyEvent{Code: keyboard.KeyLeft, Mod: keyboard.ModShift}},
		{name: "alt right", in: []byte("\x1b[1;3C"), want: keyboard.KeyEvent{Code: keyboard.KeyRight, Mod: keyboard.ModAlt}},
		{name: "shift ctrl page down", in: []byte("\x1b[6;6~"), want: keyboard.KeyEvent{Code: keyboard.KeyPageDown, Mod: keyboard.ModShift | keyboard.ModCtrl}},
		{name: "shift f5", in: []byte("\x1b[15;2~"), want: keyboard.KeyEvent{Code: keyboard.KeyF5, Mod: keyboard.ModShift}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []keyboard.KeyEvent{tt.want})
		})
	}
}

func TestParserParsesSGRMouseEvents(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want mouse.MouseEvent
	}{
		{
			name: "left press",
			in:   []byte("\x1b[<0;10;5M"),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonLeft, Action: mouse.MouseActionPress},
		},
		{
			name: "middle press",
			in:   []byte("\x1b[<1;10;5M"),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonMiddle, Action: mouse.MouseActionPress},
		},
		{
			name: "right press",
			in:   []byte("\x1b[<2;10;5M"),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonRight, Action: mouse.MouseActionPress},
		},
		{
			name: "release",
			in:   []byte("\x1b[<0;10;5m"),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonLeft, Action: mouse.MouseActionRelease},
		},
		{
			name: "drag",
			in:   []byte("\x1b[<32;10;5M"),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonLeft, Action: mouse.MouseActionDrag},
		},
		{
			name: "wheel up",
			in:   []byte("\x1b[<64;10;5M"),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonWheelUp, Action: mouse.MouseActionWheel},
		},
		{
			name: "wheel down",
			in:   []byte("\x1b[<65;10;5M"),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonWheelDown, Action: mouse.MouseActionWheel},
		},
		{
			name: "modifiers",
			in:   []byte("\x1b[<28;10;5M"),
			want: mouse.MouseEvent{
				X:      9,
				Y:      4,
				Button: mouse.MouseButtonLeft,
				Action: mouse.MouseActionPress,
				Mod:    keyboard.ModShift | keyboard.ModAlt | keyboard.ModCtrl,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			got := parser.Feed(tt.in)
			assertInputEvents(t, got, []Event{{Type: EventTypeMouse, Mouse: tt.want}})
		})
	}
}

func TestParserParsesSplitSGRMouseEvent(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[<0;10"))
	assertInputEvents(t, got, nil)
	if !got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = false, want true")
	}

	got = parser.Feed([]byte(";5M"))
	assertInputEvents(t, got, []Event{{
		Type: EventTypeMouse,
		Mouse: mouse.MouseEvent{
			X:      9,
			Y:      4,
			Button: mouse.MouseButtonLeft,
			Action: mouse.MouseActionPress,
		},
	}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesSGRMouseEventBeforeFollowingInput(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[<64;10;5Ma"))
	assertInputEvents(t, got, []Event{
		{
			Type: EventTypeMouse,
			Mouse: mouse.MouseEvent{
				X:      9,
				Y:      4,
				Button: mouse.MouseButtonWheelUp,
				Action: mouse.MouseActionWheel,
			},
		},
		{
			Type: EventTypeKey,
			Key:  keyboard.KeyEvent{Code: keyboard.KeyRune, Text: "a", Mod: keyboard.ModNone},
		},
	})
}

func TestParserParsesLegacyMouseEvents(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want mouse.MouseEvent
	}{
		{
			name: "left press",
			in:   legacyMouseInput(0, 10, 5),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonLeft, Action: mouse.MouseActionPress},
		},
		{
			name: "middle press",
			in:   legacyMouseInput(1, 10, 5),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonMiddle, Action: mouse.MouseActionPress},
		},
		{
			name: "right press",
			in:   legacyMouseInput(2, 10, 5),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonRight, Action: mouse.MouseActionPress},
		},
		{
			name: "release",
			in:   legacyMouseInput(3, 10, 5),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonUnknown, Action: mouse.MouseActionRelease},
		},
		{
			name: "drag",
			in:   legacyMouseInput(32, 10, 5),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonLeft, Action: mouse.MouseActionDrag},
		},
		{
			name: "wheel up",
			in:   legacyMouseInput(64, 10, 5),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonWheelUp, Action: mouse.MouseActionWheel},
		},
		{
			name: "wheel down",
			in:   legacyMouseInput(65, 10, 5),
			want: mouse.MouseEvent{X: 9, Y: 4, Button: mouse.MouseButtonWheelDown, Action: mouse.MouseActionWheel},
		},
		{
			name: "modifiers",
			in:   legacyMouseInput(28, 10, 5),
			want: mouse.MouseEvent{
				X:      9,
				Y:      4,
				Button: mouse.MouseButtonLeft,
				Action: mouse.MouseActionPress,
				Mod:    keyboard.ModShift | keyboard.ModAlt | keyboard.ModCtrl,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			got := parser.Feed(tt.in)
			assertInputEvents(t, got, []Event{{Type: EventTypeMouse, Mouse: tt.want}})
		})
	}
}

func TestParserParsesSplitLegacyMouseEvent(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[M"))
	assertInputEvents(t, got, nil)
	if !got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = false, want true")
	}

	got = parser.Feed(legacyMouseInputTail(0, 10, 5))
	assertInputEvents(t, got, []Event{{
		Type: EventTypeMouse,
		Mouse: mouse.MouseEvent{
			X:      9,
			Y:      4,
			Button: mouse.MouseButtonLeft,
			Action: mouse.MouseActionPress,
		},
	}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesLegacyMouseEventBeforeFollowingInput(t *testing.T) {
	parser := NewParser()

	got := parser.Feed(append(legacyMouseInput(64, 10, 5), 'a'))
	assertInputEvents(t, got, []Event{
		{
			Type: EventTypeMouse,
			Mouse: mouse.MouseEvent{
				X:      9,
				Y:      4,
				Button: mouse.MouseButtonWheelUp,
				Action: mouse.MouseActionWheel,
			},
		},
		{
			Type: EventTypeKey,
			Key:  keyboard.KeyEvent{Code: keyboard.KeyRune, Text: "a", Mod: keyboard.ModNone},
		},
	})
}

func TestParserParsesMalformedLegacyMouseEventAsUnknownKey(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte{'\x1b', '[', 'M', 32, 32, 32})
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyUnknown}})
}

func TestParserParsesMalformedSGRMouseEventAsUnknownKey(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[<0;10M"))
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyUnknown}})
}

func TestParserParsesUnknownCSISequence(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[999z"))
	assertEvents(t, got, []keyboard.KeyEvent{{Code: keyboard.KeyUnknown}})
}

func assertEvents(t *testing.T, got ParseResult, want []keyboard.KeyEvent) {
	t.Helper()

	events := make([]Event, 0, len(want))
	for _, event := range want {
		events = append(events, Event{Type: EventTypeKey, Key: event})
	}
	assertInputEvents(t, got, events)
}

func assertInputEvents(t *testing.T, got ParseResult, want []Event) {
	t.Helper()

	if len(got.Events) != len(want) {
		t.Fatalf("events = %#v, want %#v", got, want)
	}
	for i := range want {
		if got.Events[i] != want[i] {
			t.Fatalf("event %d = %#v, want %#v", i, got.Events[i], want[i])
		}
	}
}

func legacyMouseInput(code, x, y int) []byte {
	return append([]byte{'\x1b', '[', 'M'}, legacyMouseInputTail(code, x, y)...)
}

func legacyMouseInputTail(code, x, y int) []byte {
	return []byte{byte(code + 32), byte(x + 32), byte(y + 32)}
}
