package input

import (
	"github.com/NLipatov/tuigo/ansi"
	"github.com/NLipatov/tuigo/keyboard"
)

var escapeSequences = []struct {
	sequence []byte
	event    keyboard.KeyEvent
}{
	{[]byte(ansi.CSI + "A"), keyboard.KeyEvent{Code: keyboard.KeyUp}},
	{[]byte(ansi.CSI + "B"), keyboard.KeyEvent{Code: keyboard.KeyDown}},
	{[]byte(ansi.CSI + "C"), keyboard.KeyEvent{Code: keyboard.KeyRight}},
	{[]byte(ansi.CSI + "D"), keyboard.KeyEvent{Code: keyboard.KeyLeft}},
	{[]byte(ansi.CSI + "H"), keyboard.KeyEvent{Code: keyboard.KeyHome}},
	{[]byte(ansi.CSI + "F"), keyboard.KeyEvent{Code: keyboard.KeyEnd}},
	{[]byte(ansi.CSI + "3~"), keyboard.KeyEvent{Code: keyboard.KeyDelete}},
	{[]byte(ansi.CSI + "5~"), keyboard.KeyEvent{Code: keyboard.KeyPageUp}},
	{[]byte(ansi.CSI + "6~"), keyboard.KeyEvent{Code: keyboard.KeyPageDown}},
	{[]byte(ansi.CSI + "Z"), keyboard.KeyEvent{Code: keyboard.KeyTab, Mod: keyboard.ModShift}},
	{[]byte(ansi.CSI + "11~"), keyboard.KeyEvent{Code: keyboard.KeyF1}},
	{[]byte(ansi.CSI + "12~"), keyboard.KeyEvent{Code: keyboard.KeyF2}},
	{[]byte(ansi.CSI + "13~"), keyboard.KeyEvent{Code: keyboard.KeyF3}},
	{[]byte(ansi.CSI + "14~"), keyboard.KeyEvent{Code: keyboard.KeyF4}},
	{[]byte(ansi.CSI + "15~"), keyboard.KeyEvent{Code: keyboard.KeyF5}},
	{[]byte(ansi.CSI + "17~"), keyboard.KeyEvent{Code: keyboard.KeyF6}},
	{[]byte(ansi.CSI + "18~"), keyboard.KeyEvent{Code: keyboard.KeyF7}},
	{[]byte(ansi.CSI + "19~"), keyboard.KeyEvent{Code: keyboard.KeyF8}},
	{[]byte(ansi.CSI + "20~"), keyboard.KeyEvent{Code: keyboard.KeyF9}},
	{[]byte(ansi.CSI + "21~"), keyboard.KeyEvent{Code: keyboard.KeyF10}},
	{[]byte(ansi.CSI + "23~"), keyboard.KeyEvent{Code: keyboard.KeyF11}},
	{[]byte(ansi.CSI + "24~"), keyboard.KeyEvent{Code: keyboard.KeyF12}},
}

var ss3Sequences = []struct {
	sequence []byte
	event    keyboard.KeyEvent
}{
	{[]byte(ansi.SS3 + "H"), keyboard.KeyEvent{Code: keyboard.KeyHome}},
	{[]byte(ansi.SS3 + "F"), keyboard.KeyEvent{Code: keyboard.KeyEnd}},
	{[]byte(ansi.SS3 + "P"), keyboard.KeyEvent{Code: keyboard.KeyF1}},
	{[]byte(ansi.SS3 + "Q"), keyboard.KeyEvent{Code: keyboard.KeyF2}},
	{[]byte(ansi.SS3 + "R"), keyboard.KeyEvent{Code: keyboard.KeyF3}},
	{[]byte(ansi.SS3 + "S"), keyboard.KeyEvent{Code: keyboard.KeyF4}},
}

var controlBytes = []struct {
	b     byte
	event keyboard.KeyEvent
}{
	{'\r', keyboard.KeyEvent{Code: keyboard.KeyEnter}},
	{'\n', keyboard.KeyEvent{Code: keyboard.KeyEnter}},
	{'\t', keyboard.KeyEvent{Code: keyboard.KeyTab}},
	{0x7f, keyboard.KeyEvent{Code: keyboard.KeyBackspace}},
	{0x08, keyboard.KeyEvent{Code: keyboard.KeyBackspace}},
}
