package input

import "strings"

const (
	sgrMousePressFinal   = 'M'
	sgrMouseReleaseFinal = 'm'

	sgrMouseButtonMask = 0b11
	sgrMouseShiftMask  = 0b00100
	sgrMouseAltMask    = 0b01000
	sgrMouseCtrlMask   = 0b10000
	sgrMouseDragMask   = 0b100000
	sgrMouseWheelMask  = 0b1000000
)

func isSGRMouseSequence(params string, final byte) bool {
	return strings.HasPrefix(params, "<") &&
		(final == sgrMousePressFinal || final == sgrMouseReleaseFinal)
}

func mouseEventFromSGR(params string, final byte) (MouseEvent, bool) {
	values := splitCSIParams(strings.TrimPrefix(params, "<"))
	if len(values) != 3 {
		return MouseEvent{}, false
	}

	code, x, y := values[0], values[1], values[2]
	if x < 1 || y < 1 {
		return MouseEvent{}, false
	}

	return MouseEvent{
		X:      x - 1,
		Y:      y - 1,
		Button: mouseButtonFromSGR(code),
		Action: mouseActionFromSGR(code, final),
		Mod:    mouseModFromSGR(code),
	}, true
}

func mouseButtonFromSGR(code int) MouseButton {
	button := code & sgrMouseButtonMask
	if code&sgrMouseWheelMask != 0 {
		switch button {
		case 0:
			return MouseButtonWheelUp
		case 1:
			return MouseButtonWheelDown
		default:
			return MouseButtonUnknown
		}
	}

	switch button {
	case 0:
		return MouseButtonLeft
	case 1:
		return MouseButtonMiddle
	case 2:
		return MouseButtonRight
	default:
		return MouseButtonUnknown
	}
}

func mouseActionFromSGR(code int, final byte) MouseAction {
	if code&sgrMouseWheelMask != 0 {
		return MouseActionWheel
	}
	if final == sgrMouseReleaseFinal {
		return MouseActionRelease
	}
	if code&sgrMouseDragMask != 0 {
		return MouseActionDrag
	}
	return MouseActionPress
}

func mouseModFromSGR(code int) KeyMod {
	var mod KeyMod
	if code&sgrMouseShiftMask != 0 {
		mod |= ModShift
	}
	if code&sgrMouseAltMask != 0 {
		mod |= ModAlt
	}
	if code&sgrMouseCtrlMask != 0 {
		mod |= ModCtrl
	}
	return mod
}
