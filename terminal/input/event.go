package input

import (
	"github.com/NLipatov/tuigo/keyboard"
	"github.com/NLipatov/tuigo/mouse"
)

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeKey
	EventTypeMouse
)

type Event struct {
	Type  EventType
	Key   keyboard.KeyEvent
	Mouse mouse.MouseEvent
}
