package mouse

import "github.com/NLipatov/tuigo/keyboard"

type MouseEvent struct {
	X, Y   int
	Button MouseButton
	Action MouseAction
	Mod    keyboard.KeyMod
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
