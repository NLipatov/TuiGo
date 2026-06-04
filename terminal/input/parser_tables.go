package input

import "github.com/NLipatov/tuigo/ansi"

var escapeSequences = []struct {
	sequence []byte
	event    KeyEvent
}{
	{[]byte(ansi.CSI + "A"), KeyEvent{Code: KeyUp}},
	{[]byte(ansi.CSI + "B"), KeyEvent{Code: KeyDown}},
	{[]byte(ansi.CSI + "C"), KeyEvent{Code: KeyRight}},
	{[]byte(ansi.CSI + "D"), KeyEvent{Code: KeyLeft}},
	{[]byte(ansi.CSI + "H"), KeyEvent{Code: KeyHome}},
	{[]byte(ansi.CSI + "F"), KeyEvent{Code: KeyEnd}},
	{[]byte(ansi.CSI + "3~"), KeyEvent{Code: KeyDelete}},
	{[]byte(ansi.CSI + "5~"), KeyEvent{Code: KeyPageUp}},
	{[]byte(ansi.CSI + "6~"), KeyEvent{Code: KeyPageDown}},
	{[]byte(ansi.CSI + "Z"), KeyEvent{Code: KeyTab, Mod: ModShift}},
	{[]byte(ansi.CSI + "11~"), KeyEvent{Code: KeyF1}},
	{[]byte(ansi.CSI + "12~"), KeyEvent{Code: KeyF2}},
	{[]byte(ansi.CSI + "13~"), KeyEvent{Code: KeyF3}},
	{[]byte(ansi.CSI + "14~"), KeyEvent{Code: KeyF4}},
	{[]byte(ansi.CSI + "15~"), KeyEvent{Code: KeyF5}},
	{[]byte(ansi.CSI + "17~"), KeyEvent{Code: KeyF6}},
	{[]byte(ansi.CSI + "18~"), KeyEvent{Code: KeyF7}},
	{[]byte(ansi.CSI + "19~"), KeyEvent{Code: KeyF8}},
	{[]byte(ansi.CSI + "20~"), KeyEvent{Code: KeyF9}},
	{[]byte(ansi.CSI + "21~"), KeyEvent{Code: KeyF10}},
	{[]byte(ansi.CSI + "23~"), KeyEvent{Code: KeyF11}},
	{[]byte(ansi.CSI + "24~"), KeyEvent{Code: KeyF12}},
}

var ss3Sequences = []struct {
	sequence []byte
	event    KeyEvent
}{
	{[]byte(ansi.SS3 + "H"), KeyEvent{Code: KeyHome}},
	{[]byte(ansi.SS3 + "F"), KeyEvent{Code: KeyEnd}},
	{[]byte(ansi.SS3 + "P"), KeyEvent{Code: KeyF1}},
	{[]byte(ansi.SS3 + "Q"), KeyEvent{Code: KeyF2}},
	{[]byte(ansi.SS3 + "R"), KeyEvent{Code: KeyF3}},
	{[]byte(ansi.SS3 + "S"), KeyEvent{Code: KeyF4}},
}

var controlBytes = []struct {
	b     byte
	event KeyEvent
}{
	{'\r', KeyEvent{Code: KeyEnter}},
	{'\n', KeyEvent{Code: KeyEnter}},
	{'\t', KeyEvent{Code: KeyTab}},
	{0x7f, KeyEvent{Code: KeyBackspace}},
	{0x08, KeyEvent{Code: KeyBackspace}},
}
