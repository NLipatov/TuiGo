package input

import "github.com/NLipatov/tuigo/ansi"

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
