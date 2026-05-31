package input

import (
	"testing"
)

func TestInputParserParsesRunes(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte("aЖ"))
	want := []Event{
		{Code: KeyRune, Text: "a", Mod: ModNone},
		{Code: KeyRune, Text: "Ж", Mod: ModNone},
	}
	assertEvents(t, got, want)
}

func TestInputParserWaitsForIncompleteUTF8Rune(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte{0xd0})
	assertEvents(t, got, nil)

	got = parser.Feed([]byte{0x96})
	assertEvents(t, got, []Event{{Code: KeyRune, Text: "Ж", Mod: ModNone}})
}

func TestInputParserParsesControlBytes(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want Event
	}{
		{name: "carriage return", in: []byte{'\r'}, want: Event{Code: KeyEnter}},
		{name: "line feed", in: []byte{'\n'}, want: Event{Code: KeyEnter}},
		{name: "tab", in: []byte{'\t'}, want: Event{Code: KeyTab}},
		{name: "delete backspace", in: []byte{0x7f}, want: Event{Code: KeyBackspace}},
		{name: "ctrl h backspace", in: []byte{0x08}, want: Event{Code: KeyBackspace}},
		{name: "ctrl c", in: []byte{0x03}, want: Event{Code: KeyRune, Text: "c", Mod: ModCtrl}},
		{name: "ctrl d", in: []byte{0x04}, want: Event{Code: KeyRune, Text: "d", Mod: ModCtrl}},
		{name: "ctrl z", in: []byte{0x1a}, want: Event{Code: KeyRune, Text: "z", Mod: ModCtrl}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewInputParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []Event{tt.want})
		})
	}
}

func TestInputParserParsesEscapeSequences(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want Event
	}{
		{name: "up", in: []byte("\x1b[A"), want: Event{Code: KeyUp}},
		{name: "down", in: []byte("\x1b[B"), want: Event{Code: KeyDown}},
		{name: "right", in: []byte("\x1b[C"), want: Event{Code: KeyRight}},
		{name: "left", in: []byte("\x1b[D"), want: Event{Code: KeyLeft}},
		{name: "home csi", in: []byte("\x1b[H"), want: Event{Code: KeyHome}},
		{name: "end csi", in: []byte("\x1b[F"), want: Event{Code: KeyEnd}},
		{name: "home ss3", in: []byte("\x1bOH"), want: Event{Code: KeyHome}},
		{name: "end ss3", in: []byte("\x1bOF"), want: Event{Code: KeyEnd}},
		{name: "delete", in: []byte("\x1b[3~"), want: Event{Code: KeyDelete}},
		{name: "page up", in: []byte("\x1b[5~"), want: Event{Code: KeyPageUp}},
		{name: "page down", in: []byte("\x1b[6~"), want: Event{Code: KeyPageDown}},
		{name: "shift tab", in: []byte("\x1b[Z"), want: Event{Code: KeyTab, Mod: ModShift}},
		{name: "insert", in: []byte("\x1b[2~"), want: Event{Code: KeyInsert}},
		{name: "f1 ss3", in: []byte("\x1bOP"), want: Event{Code: KeyF1}},
		{name: "f4 ss3", in: []byte("\x1bOS"), want: Event{Code: KeyF4}},
		{name: "f5 csi", in: []byte("\x1b[15~"), want: Event{Code: KeyF5}},
		{name: "f12 csi", in: []byte("\x1b[24~"), want: Event{Code: KeyF12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewInputParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []Event{tt.want})
		})
	}
}

func TestInputParserTimesOutPendingEsc(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte("\x1b"))
	assertEvents(t, got, nil)
	if !got.NeedsTimeout {
		t.Fatalf("NeedsTimeout = false, want true")
	}

	got = parser.Timeout()
	assertEvents(t, got, []Event{{Code: KeyEsc}})
	if got.NeedsTimeout {
		t.Fatalf("NeedsTimeout = true, want false")
	}
}

func TestInputParserParsesSplitEscapeSequence(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte("\x1b"))
	assertEvents(t, got, nil)
	if !got.NeedsTimeout {
		t.Fatalf("NeedsTimeout = false, want true")
	}

	got = parser.Feed([]byte("[A"))
	assertEvents(t, got, []Event{{Code: KeyUp}})
	if got.NeedsTimeout {
		t.Fatalf("NeedsTimeout = true, want false")
	}
}

func TestInputParserParsesEscapeSequencesBeforeFollowingInput(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte("\x1b[Aa"))
	want := []Event{
		{Code: KeyUp},
		{Code: KeyRune, Text: "a", Mod: ModNone},
	}
	assertEvents(t, got, want)
}

func TestInputParserWaitsForIncompleteEscapeSequence(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte("\x1b["))
	assertEvents(t, got, nil)
	if !got.NeedsTimeout {
		t.Fatalf("NeedsTimeout = false, want true")
	}

	got = parser.Feed([]byte("A"))
	assertEvents(t, got, []Event{{Code: KeyUp}})
	if got.NeedsTimeout {
		t.Fatalf("NeedsTimeout = true, want false")
	}
}

func TestInputParserParsesAltRune(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte("\x1bx"))
	assertEvents(t, got, []Event{{Code: KeyRune, Text: "x", Mod: ModAlt}})
}

func TestInputParserParsesAltRuneSplitAcrossFeeds(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte("\x1b"))
	assertEvents(t, got, nil)
	if !got.NeedsTimeout {
		t.Fatalf("NeedsTimeout = false, want true")
	}

	got = parser.Feed([]byte("x"))
	assertEvents(t, got, []Event{{Code: KeyRune, Text: "x", Mod: ModAlt}})
	if got.NeedsTimeout {
		t.Fatalf("NeedsTimeout = true, want false")
	}
}

func TestInputParserParsesCSIModifiers(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want Event
	}{
		{name: "ctrl up", in: []byte("\x1b[1;5A"), want: Event{Code: KeyUp, Mod: ModCtrl}},
		{name: "shift left", in: []byte("\x1b[1;2D"), want: Event{Code: KeyLeft, Mod: ModShift}},
		{name: "alt right", in: []byte("\x1b[1;3C"), want: Event{Code: KeyRight, Mod: ModAlt}},
		{name: "shift ctrl page down", in: []byte("\x1b[6;6~"), want: Event{Code: KeyPageDown, Mod: ModShift | ModCtrl}},
		{name: "shift f5", in: []byte("\x1b[15;2~"), want: Event{Code: KeyF5, Mod: ModShift}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewInputParser()
			got := parser.Feed(tt.in)
			assertEvents(t, got, []Event{tt.want})
		})
	}
}

func TestInputParserParsesUnknownCSISequence(t *testing.T) {
	parser := NewInputParser()

	got := parser.Feed([]byte("\x1b[999z"))
	assertEvents(t, got, []Event{{Code: KeyUnknown}})
}

func assertEvents(t *testing.T, got ParseResult, want []Event) {
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
