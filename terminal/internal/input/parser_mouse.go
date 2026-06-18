package input

import (
	"strings"

	"github.com/NLipatov/tuigo/keyboard"
	"github.com/NLipatov/tuigo/mouse"
)

const (
	sgrMousePressFinal   = 'M'
	sgrMouseReleaseFinal = 'm'

	sgrMouseButtonMask = 0b11
	sgrMouseShiftMask  = 0b00100
	sgrMouseAltMask    = 0b01000
	sgrMouseCtrlMask   = 0b10000
	sgrMouseDragMask   = 0b100000
	sgrMouseWheelMask  = 0b1000000

	legacyMouseFinal        = 'M'
	legacyMouseSequenceLen  = 6
	legacyMouseEncodingBase = 32
)

func isSGRMouseSequence(params string, final byte) bool {
	return strings.HasPrefix(params, "<") &&
		(final == sgrMousePressFinal || final == sgrMouseReleaseFinal)
}

func mouseEventFromSGR(params string, final byte) (mouse.MouseEvent, bool) {
	values := splitCSIParams(strings.TrimPrefix(params, "<"))
	if len(values) != 3 {
		return mouse.MouseEvent{}, false
	}

	code, x, y := values[0], values[1], values[2]
	if x < 1 || y < 1 {
		return mouse.MouseEvent{}, false
	}

	return mouse.MouseEvent{
		X:      x - 1,
		Y:      y - 1,
		Button: mouseButtonFromSGR(code),
		Action: mouseActionFromSGR(code, final),
		Mod:    mouseModFromSGR(code),
	}, true
}

func mouseButtonFromSGR(code int) mouse.MouseButton {
	button := code & sgrMouseButtonMask
	if code&sgrMouseWheelMask != 0 {
		switch button {
		case 0:
			return mouse.MouseButtonWheelUp
		case 1:
			return mouse.MouseButtonWheelDown
		default:
			return mouse.MouseButtonUnknown
		}
	}

	switch button {
	case 0:
		return mouse.MouseButtonLeft
	case 1:
		return mouse.MouseButtonMiddle
	case 2:
		return mouse.MouseButtonRight
	default:
		return mouse.MouseButtonUnknown
	}
}

func mouseActionFromSGR(code int, final byte) mouse.MouseAction {
	if code&sgrMouseWheelMask != 0 {
		return mouse.MouseActionWheel
	}
	if final == sgrMouseReleaseFinal {
		return mouse.MouseActionRelease
	}
	if code&sgrMouseDragMask != 0 {
		return mouse.MouseActionDrag
	}
	return mouse.MouseActionPress
}

func mouseModFromSGR(code int) keyboard.KeyMod {
	var mod keyboard.KeyMod
	if code&sgrMouseShiftMask != 0 {
		mod |= keyboard.ModShift
	}
	if code&sgrMouseAltMask != 0 {
		mod |= keyboard.ModAlt
	}
	if code&sgrMouseCtrlMask != 0 {
		mod |= keyboard.ModCtrl
	}
	return mod
}

func legacyMouseEventFromCSI(buf []byte) (Event, int, parseStatus) {
	if len(buf) < len("\x1b[M") || buf[2] != legacyMouseFinal {
		return Event{}, 0, parseNoMatch
	}
	if len(buf) < legacyMouseSequenceLen {
		return Event{}, 0, parseNeedMore
	}

	mouse, ok := mouseEventFromLegacy(buf[3], buf[4], buf[5])
	if !ok {
		return Event{Type: EventTypeKey, Key: keyboard.KeyEvent{Code: keyboard.KeyUnknown}}, legacyMouseSequenceLen, parseDone
	}
	return Event{Type: EventTypeMouse, Mouse: mouse}, legacyMouseSequenceLen, parseDone
}

func mouseEventFromLegacy(codeByte, xByte, yByte byte) (mouse.MouseEvent, bool) {
	code := int(codeByte) - legacyMouseEncodingBase
	x := int(xByte) - legacyMouseEncodingBase - 1
	y := int(yByte) - legacyMouseEncodingBase - 1
	if code < 0 || x < 0 || y < 0 {
		return mouse.MouseEvent{}, false
	}

	return mouse.MouseEvent{
		X:      x,
		Y:      y,
		Button: mouseButtonFromLegacy(code),
		Action: mouseActionFromLegacy(code),
		Mod:    mouseModFromSGR(code),
	}, true
}

func mouseButtonFromLegacy(code int) mouse.MouseButton {
	if code&sgrMouseWheelMask != 0 {
		return mouseButtonFromSGR(code)
	}
	if code&sgrMouseButtonMask == sgrMouseButtonMask {
		return mouse.MouseButtonUnknown
	}
	return mouseButtonFromSGR(code)
}

func mouseActionFromLegacy(code int) mouse.MouseAction {
	if code&sgrMouseWheelMask != 0 {
		return mouse.MouseActionWheel
	}
	if code&sgrMouseButtonMask == sgrMouseButtonMask {
		return mouse.MouseActionRelease
	}
	if code&sgrMouseDragMask != 0 {
		return mouse.MouseActionDrag
	}
	return mouse.MouseActionPress
}
