package input

type MouseEvent struct {
	X, Y   int
	Button MouseButton
	Action MouseAction
	Mod    KeyMod
}

type MouseButton int

const (
	MouseButtonUnknown MouseButton = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
	MouseButtonWheelUp
	MouseButtonWheelDown
)

type MouseAction int

const (
	MouseActionUnknown MouseAction = iota
	MouseActionPress
	MouseActionRelease
	MouseActionDrag
	MouseActionWheel
)
