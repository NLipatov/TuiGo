package input

type KeyEvent struct {
	Code KeyCode
	Text string
	Mod  KeyMod
}

type KeyCode int

const (
	KeyUnknown KeyCode = iota
	KeyRune

	KeyEnter
	KeyEsc
	KeyTab
	KeyBackspace
	KeyDelete
	KeyInsert

	KeyUp
	KeyDown
	KeyLeft
	KeyRight

	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown

	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

type KeyMod uint8

const (
	ModNone  KeyMod = 0 // 0000_0000
	ModCtrl  KeyMod = 1 // 0000_0001
	ModAlt   KeyMod = 2 // 0000_0010
	ModShift KeyMod = 4 // 0000_0100
)
