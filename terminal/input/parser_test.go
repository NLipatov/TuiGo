package input

import (
	"testing"
)

func TestParserParsesRunes(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("aЖ"))
	want := []KeyEvent{
		{Code: KeyRune, Text: "a", Mod: ModNone},
		{Code: KeyRune, Text: "Ж", Mod: ModNone},
	}
	assertEvents(t, got, want)
}

func TestParserWaitsForIncompleteUTF8Rune(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte{0xd0})
	assertEvents(t, got, nil)

	got = parser.Feed([]byte{0x96})
	assertEvents(t, got, []KeyEvent{{Code: KeyRune, Text: "Ж", Mod: ModNone}})
}

func TestParserParsesControlBytes(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want KeyEvent
	}{
		{name: "carriage return", in: []byte{'\r'}, want: KeyEvent{Code: KeyEnter}},
		{name: "line feed", in: []byte{'\n'}, want: KeyEvent{Code: KeyEnter}},
		{name: "tab", in: []byte{'\t'}, want: KeyEvent{Code: KeyTab}},
		{name: "delete backspace", in: []byte{0x7f}, want: KeyEvent{Code: KeyBackspace}},
		{name: "ctrl h backspace", in: []byte{0x08}, want: KeyEvent{Code: KeyBackspace}},
		{name: "ctrl c", in: []byte{0x03}, want: KeyEvent{Code: KeyRune, Text: "c", Mod: ModCtrl}},
		{name: "ctrl d", in: []byte{0x04}, want: KeyEvent{Code: KeyRune, Text: "d", Mod: ModCtrl}},
		{name: "ctrl z", in: []byte{0x1a}, want: KeyEvent{Code: KeyRune, Text: "z", Mod: ModCtrl}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []KeyEvent{tt.want})
		})
	}
}

func TestParserParsesEscapeSequences(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want KeyEvent
	}{
		{name: "up", in: []byte("\x1b[A"), want: KeyEvent{Code: KeyUp}},
		{name: "down", in: []byte("\x1b[B"), want: KeyEvent{Code: KeyDown}},
		{name: "right", in: []byte("\x1b[C"), want: KeyEvent{Code: KeyRight}},
		{name: "left", in: []byte("\x1b[D"), want: KeyEvent{Code: KeyLeft}},
		{name: "home csi", in: []byte("\x1b[H"), want: KeyEvent{Code: KeyHome}},
		{name: "end csi", in: []byte("\x1b[F"), want: KeyEvent{Code: KeyEnd}},
		{name: "home ss3", in: []byte("\x1bOH"), want: KeyEvent{Code: KeyHome}},
		{name: "end ss3", in: []byte("\x1bOF"), want: KeyEvent{Code: KeyEnd}},
		{name: "delete", in: []byte("\x1b[3~"), want: KeyEvent{Code: KeyDelete}},
		{name: "page up", in: []byte("\x1b[5~"), want: KeyEvent{Code: KeyPageUp}},
		{name: "page down", in: []byte("\x1b[6~"), want: KeyEvent{Code: KeyPageDown}},
		{name: "shift tab", in: []byte("\x1b[Z"), want: KeyEvent{Code: KeyTab, Mod: ModShift}},
		{name: "insert", in: []byte("\x1b[2~"), want: KeyEvent{Code: KeyInsert}},
		{name: "f1 ss3", in: []byte("\x1bOP"), want: KeyEvent{Code: KeyF1}},
		{name: "f4 ss3", in: []byte("\x1bOS"), want: KeyEvent{Code: KeyF4}},
		{name: "f5 csi", in: []byte("\x1b[15~"), want: KeyEvent{Code: KeyF5}},
		{name: "f12 csi", in: []byte("\x1b[24~"), want: KeyEvent{Code: KeyF12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []KeyEvent{tt.want})
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
	assertEvents(t, got, []KeyEvent{{Code: KeyEsc}})
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
	assertEvents(t, got, []KeyEvent{{Code: KeyUp}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesEscapeSequencesBeforeFollowingInput(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[Aa"))
	want := []KeyEvent{
		{Code: KeyUp},
		{Code: KeyRune, Text: "a", Mod: ModNone},
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
	assertEvents(t, got, []KeyEvent{{Code: KeyUp}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesAltRune(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1bx"))
	assertEvents(t, got, []KeyEvent{{Code: KeyRune, Text: "x", Mod: ModAlt}})
}

func TestParserParsesAltRuneSplitAcrossFeeds(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b"))
	assertEvents(t, got, nil)
	if !got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = false, want true")
	}

	got = parser.Feed([]byte("x"))
	assertEvents(t, got, []KeyEvent{{Code: KeyRune, Text: "x", Mod: ModAlt}})
	if got.HasPendingEscape {
		t.Fatalf("HasPendingEscape = true, want false")
	}
}

func TestParserParsesCSIModifiers(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want KeyEvent
	}{
		{name: "ctrl up", in: []byte("\x1b[1;5A"), want: KeyEvent{Code: KeyUp, Mod: ModCtrl}},
		{name: "shift left", in: []byte("\x1b[1;2D"), want: KeyEvent{Code: KeyLeft, Mod: ModShift}},
		{name: "alt right", in: []byte("\x1b[1;3C"), want: KeyEvent{Code: KeyRight, Mod: ModAlt}},
		{name: "shift ctrl page down", in: []byte("\x1b[6;6~"), want: KeyEvent{Code: KeyPageDown, Mod: ModShift | ModCtrl}},
		{name: "shift f5", in: []byte("\x1b[15;2~"), want: KeyEvent{Code: KeyF5, Mod: ModShift}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []KeyEvent{tt.want})
		})
	}
}

func TestParserParsesSGRMouseEvents(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want MouseEvent
	}{
		{
			name: "left press",
			in:   []byte("\x1b[<0;10;5M"),
			want: MouseEvent{X: 9, Y: 4, Button: MouseButtonLeft, Action: MouseActionPress},
		},
		{
			name: "middle press",
			in:   []byte("\x1b[<1;10;5M"),
			want: MouseEvent{X: 9, Y: 4, Button: MouseButtonMiddle, Action: MouseActionPress},
		},
		{
			name: "right press",
			in:   []byte("\x1b[<2;10;5M"),
			want: MouseEvent{X: 9, Y: 4, Button: MouseButtonRight, Action: MouseActionPress},
		},
		{
			name: "release",
			in:   []byte("\x1b[<0;10;5m"),
			want: MouseEvent{X: 9, Y: 4, Button: MouseButtonLeft, Action: MouseActionRelease},
		},
		{
			name: "drag",
			in:   []byte("\x1b[<32;10;5M"),
			want: MouseEvent{X: 9, Y: 4, Button: MouseButtonLeft, Action: MouseActionDrag},
		},
		{
			name: "wheel up",
			in:   []byte("\x1b[<64;10;5M"),
			want: MouseEvent{X: 9, Y: 4, Button: MouseButtonWheelUp, Action: MouseActionWheel},
		},
		{
			name: "wheel down",
			in:   []byte("\x1b[<65;10;5M"),
			want: MouseEvent{X: 9, Y: 4, Button: MouseButtonWheelDown, Action: MouseActionWheel},
		},
		{
			name: "modifiers",
			in:   []byte("\x1b[<28;10;5M"),
			want: MouseEvent{
				X:      9,
				Y:      4,
				Button: MouseButtonLeft,
				Action: MouseActionPress,
				Mod:    ModShift | ModAlt | ModCtrl,
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
		Mouse: MouseEvent{
			X:      9,
			Y:      4,
			Button: MouseButtonLeft,
			Action: MouseActionPress,
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
			Mouse: MouseEvent{
				X:      9,
				Y:      4,
				Button: MouseButtonWheelUp,
				Action: MouseActionWheel,
			},
		},
		{
			Type: EventTypeKey,
			Key:  KeyEvent{Code: KeyRune, Text: "a", Mod: ModNone},
		},
	})
}

func TestParserParsesMalformedSGRMouseEventAsUnknownKey(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[<0;10M"))
	assertEvents(t, got, []KeyEvent{{Code: KeyUnknown}})
}

func TestParserParsesUnknownCSISequence(t *testing.T) {
	parser := NewParser()

	got := parser.Feed([]byte("\x1b[999z"))
	assertEvents(t, got, []KeyEvent{{Code: KeyUnknown}})
}

func assertEvents(t *testing.T, got ParseResult, want []KeyEvent) {
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
