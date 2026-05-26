package terminal

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
