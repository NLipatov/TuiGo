package input

import (
	"strconv"
	"strings"
)

func eventForCSIFinal(final byte, params string) (KeyEvent, bool) {
	switch final {
	case 'A':
		return KeyEvent{Code: KeyUp, Mod: modFromCSIParams(params)}, true
	case 'B':
		return KeyEvent{Code: KeyDown, Mod: modFromCSIParams(params)}, true
	case 'C':
		return KeyEvent{Code: KeyRight, Mod: modFromCSIParams(params)}, true
	case 'D':
		return KeyEvent{Code: KeyLeft, Mod: modFromCSIParams(params)}, true
	case 'H':
		return KeyEvent{Code: KeyHome, Mod: modFromCSIParams(params)}, true
	case 'F':
		return KeyEvent{Code: KeyEnd, Mod: modFromCSIParams(params)}, true
	case 'Z':
		return KeyEvent{Code: KeyTab, Mod: ModShift}, true
	case '~':
		return eventForCSITilde(params)
	default:
		return KeyEvent{}, false
	}
}

func eventForCSITilde(params string) (KeyEvent, bool) {
	values := splitCSIParams(params)
	if len(values) == 0 {
		return KeyEvent{}, false
	}

	event := KeyEvent{Code: keyCodeForCSITilde(values[0])}
	if event.Code == KeyUnknown {
		return event, true
	}
	if len(values) > 1 {
		event.Mod = xtermModifier(values[1])
	}
	return event, true
}

//nolint:cyclop // This switch is a flat CSI numeric key mapping table.
func keyCodeForCSITilde(value int) KeyCode {
	switch value {
	case 2:
		return KeyInsert
	case 3:
		return KeyDelete
	case 5:
		return KeyPageUp
	case 6:
		return KeyPageDown
	case 11:
		return KeyF1
	case 12:
		return KeyF2
	case 13:
		return KeyF3
	case 14:
		return KeyF4
	case 15:
		return KeyF5
	case 17:
		return KeyF6
	case 18:
		return KeyF7
	case 19:
		return KeyF8
	case 20:
		return KeyF9
	case 21:
		return KeyF10
	case 23:
		return KeyF11
	case 24:
		return KeyF12
	default:
		return KeyUnknown
	}
}

func modFromCSIParams(params string) KeyMod {
	values := splitCSIParams(params)
	if len(values) < 2 {
		return ModNone
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

func xtermModifier(value int) KeyMod {
	bits := value - 1
	if bits <= 0 {
		return ModNone
	}

	var mod KeyMod
	if bits&1 != 0 {
		mod |= ModShift
	}
	if bits&2 != 0 {
		mod |= ModAlt
	}
	if bits&4 != 0 {
		mod |= ModCtrl
	}
	return mod
}
