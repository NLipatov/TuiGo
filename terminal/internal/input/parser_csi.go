package input

import (
	"strconv"
	"strings"

	"github.com/NLipatov/tuigo/keyboard"
)

func eventForCSIFinal(final byte, params string) (keyboard.KeyEvent, bool) {
	switch final {
	case 'A':
		return keyboard.KeyEvent{Code: keyboard.KeyUp, Mod: modFromCSIParams(params)}, true
	case 'B':
		return keyboard.KeyEvent{Code: keyboard.KeyDown, Mod: modFromCSIParams(params)}, true
	case 'C':
		return keyboard.KeyEvent{Code: keyboard.KeyRight, Mod: modFromCSIParams(params)}, true
	case 'D':
		return keyboard.KeyEvent{Code: keyboard.KeyLeft, Mod: modFromCSIParams(params)}, true
	case 'H':
		return keyboard.KeyEvent{Code: keyboard.KeyHome, Mod: modFromCSIParams(params)}, true
	case 'F':
		return keyboard.KeyEvent{Code: keyboard.KeyEnd, Mod: modFromCSIParams(params)}, true
	case 'Z':
		return keyboard.KeyEvent{Code: keyboard.KeyTab, Mod: keyboard.ModShift}, true
	case '~':
		return eventForCSITilde(params)
	default:
		return keyboard.KeyEvent{}, false
	}
}

func eventForCSITilde(params string) (keyboard.KeyEvent, bool) {
	values := splitCSIParams(params)
	if len(values) == 0 {
		return keyboard.KeyEvent{}, false
	}

	event := keyboard.KeyEvent{Code: keyCodeForCSITilde(values[0])}
	if event.Code == keyboard.KeyUnknown {
		return event, true
	}
	if len(values) > 1 {
		event.Mod = xtermModifier(values[1])
	}
	return event, true
}

//nolint:cyclop // This switch is a flat CSI numeric key mapping table.
func keyCodeForCSITilde(value int) keyboard.KeyCode {
	switch value {
	case 2:
		return keyboard.KeyInsert
	case 3:
		return keyboard.KeyDelete
	case 5:
		return keyboard.KeyPageUp
	case 6:
		return keyboard.KeyPageDown
	case 11:
		return keyboard.KeyF1
	case 12:
		return keyboard.KeyF2
	case 13:
		return keyboard.KeyF3
	case 14:
		return keyboard.KeyF4
	case 15:
		return keyboard.KeyF5
	case 17:
		return keyboard.KeyF6
	case 18:
		return keyboard.KeyF7
	case 19:
		return keyboard.KeyF8
	case 20:
		return keyboard.KeyF9
	case 21:
		return keyboard.KeyF10
	case 23:
		return keyboard.KeyF11
	case 24:
		return keyboard.KeyF12
	default:
		return keyboard.KeyUnknown
	}
}

func modFromCSIParams(params string) keyboard.KeyMod {
	values := splitCSIParams(params)
	if len(values) < 2 {
		return keyboard.ModNone
	}
	return xtermModifier(values[1])
}

func splitCSIParams(params string) []int {
	if params == "" {
		return nil
	}

	parts := strings.Split(params, ";")
	values := make([]int, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			values = append(values, 0)
			continue
		}
		value, err := strconv.Atoi(part)
		if err != nil {
			return nil
		}
		values = append(values, value)
	}
	return values
}

func xtermModifier(value int) keyboard.KeyMod {
	bits := value - 1
	if bits <= 0 {
		return keyboard.ModNone
	}

	var mod keyboard.KeyMod
	if bits&1 != 0 {
		mod |= keyboard.ModShift
	}
	if bits&2 != 0 {
		mod |= keyboard.ModAlt
	}
	if bits&4 != 0 {
		mod |= keyboard.ModCtrl
	}
	return mod
}
